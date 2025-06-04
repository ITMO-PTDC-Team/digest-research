import csv
import ast
from datetime import datetime, timedelta
import numpy as np
from scipy.interpolate import interp1d
import os

P_LEVELS = np.linspace(0.001, 0.999, 999)

def create_quantile_function(quantile_values_at_p_levels):
    _p_levels_extended = np.concatenate(([0.0], P_LEVELS, [1.0]))
    _quantile_values_extended = np.concatenate((
        [quantile_values_at_p_levels[0]],
        quantile_values_at_p_levels,
        [quantile_values_at_p_levels[-1]]
    ))
    sorted_indices = np.argsort(_p_levels_extended)
    _p_levels_sorted = _p_levels_extended[sorted_indices]
    _quantile_values_sorted = _quantile_values_extended[sorted_indices]
    unique_p_levels, unique_indices = np.unique(_p_levels_sorted, return_index=True)
    _p_levels_final = unique_p_levels
    _quantile_values_final = _quantile_values_sorted[unique_indices]
    quantile_func = interp1d(
        _p_levels_final,
        _quantile_values_final,
        kind='linear',
        bounds_error=False,
        fill_value=(_quantile_values_final[0], _quantile_values_final[-1])
    )
    return quantile_func

def generate_samples(quantile_function, num_samples):
    uniform_variates = np.random.rand(num_samples)
    samples = quantile_function(uniform_variates)
    return samples

def parse_csv_data(filepath, date_format="%Y-%m-%d %H:%M:%S"):
    data = []
    expected_quantile_len = len(P_LEVELS)
    with open(filepath, 'r', encoding='utf-8') as f:
        reader = csv.reader(f)
        for i, row in enumerate(reader):
            if not row or len(row) < 2:
                continue
            try:
                timestamp_str = row[0].strip().strip('"')
                quantiles_str = row[1].strip()
                dt_obj = datetime.strptime(timestamp_str, date_format)
                
                try:
                    quantiles_list_from_str = ast.literal_eval(quantiles_str)
                except (SyntaxError, ValueError):
                    if not quantiles_str.startswith('['): quantiles_str = '[' + quantiles_str
                    if not quantiles_str.endswith(']'): quantiles_str = quantiles_str + ']'
                    try:
                        quantiles_list_from_str = ast.literal_eval(quantiles_str)
                    except Exception:
                        continue

                quantiles_array = np.array(quantiles_list_from_str, dtype=float)
                if len(quantiles_array) != expected_quantile_len:
                    continue
                data.append((dt_obj, quantiles_array))
            except Exception:
                continue
    data.sort(key=lambda x: x[0])
    return data

def align_data(seconds_file_path, minutes_file_path):
    print("Парсинг данных...")
    raw_seconds_data = parse_csv_data(seconds_file_path, "%Y-%m-%d %H:%M:%S")
    raw_minutes_data = parse_csv_data(minutes_file_path, "%Y-%m-%d %H:%M:00")

    if not raw_seconds_data or not raw_minutes_data:
        return {}, {}, []

    seconds_grouped_by_minute_raw = {}
    for dt_sec, quantiles_sec in raw_seconds_data:
        minute_start_dt = dt_sec.replace(second=0, microsecond=0)
        if minute_start_dt not in seconds_grouped_by_minute_raw:
            seconds_grouped_by_minute_raw[minute_start_dt] = []
        seconds_grouped_by_minute_raw[minute_start_dt].append((dt_sec, quantiles_sec))

    aligned_seconds_by_minute = {}
    valid_minute_starts_from_seconds = set()
    for minute_start_dt, sec_data_list in seconds_grouped_by_minute_raw.items():
        if len(sec_data_list) == 60:
            sec_data_list.sort(key=lambda x: x[0])
            is_complete_minute = True
            for i in range(60):
                if sec_data_list[i][0].second != i:
                    is_complete_minute = False
                    break
            if is_complete_minute:
                aligned_seconds_by_minute[minute_start_dt] = [item[1] for item in sec_data_list]
                valid_minute_starts_from_seconds.add(minute_start_dt)

    aligned_minutes_quantiles = {}
    for dt_min, quantiles_min in raw_minutes_data:
        if dt_min.second == 0 and dt_min.microsecond == 0:
            if dt_min in valid_minute_starts_from_seconds:
                aligned_minutes_quantiles[dt_min] = quantiles_min
    
    common_minute_keys = sorted(list(set(aligned_seconds_by_minute.keys()) & set(aligned_minutes_quantiles.keys())))
    final_seconds_map = {key: aligned_seconds_by_minute[key] for key in common_minute_keys}
    final_minutes_map = {key: aligned_minutes_quantiles[key] for key in common_minute_keys}

    return final_seconds_map, final_minutes_map, common_minute_keys

def run_simulation_and_comparison(seconds_data_map, minutes_data_map, minute_keys_sorted,
                                  samples_per_second_tick=100, base_output_dir_for_minute_details=""):
    results = {}
    target_quantile_p_levels_for_output = P_LEVELS

    print("Тест...")
    for minute_dt in minute_keys_sorted:
        second_quantile_arrays_for_minute = seconds_data_map[minute_dt]
        merged_minute_actual_quantiles = minutes_data_map[minute_dt]

        all_simulated_samples_for_minute = []
        for i, sec_quantiles_arr in enumerate(second_quantile_arrays_for_minute):
            try:
                q_func = create_quantile_function(sec_quantiles_arr)
                samples = generate_samples(q_func, samples_per_second_tick)
                all_simulated_samples_for_minute.extend(samples)
            except Exception:
                continue

        if not all_simulated_samples_for_minute:
            results[minute_dt] = {"error": "Не сгенерировано выборок"}
            continue
        
        all_simulated_samples_for_minute = np.array(all_simulated_samples_for_minute)
        simulated_minute_quantiles_values = np.percentile(
            all_simulated_samples_for_minute,
            target_quantile_p_levels_for_output * 100
        )

        if len(merged_minute_actual_quantiles) != len(simulated_minute_quantiles_values):
            results[minute_dt] = {"error": "Несовпадение длин массивов квантилей"}
            continue

        abs_diff = np.abs(merged_minute_actual_quantiles - simulated_minute_quantiles_values)
        rel_diff = np.zeros_like(abs_diff, dtype=float)
        mask_nonzero_actual = merged_minute_actual_quantiles != 0
        rel_diff[mask_nonzero_actual] = abs_diff[mask_nonzero_actual] / np.abs(merged_minute_actual_quantiles[mask_nonzero_actual])
        mask_zero_actual_nonzero_sim = (merged_minute_actual_quantiles == 0) & (simulated_minute_quantiles_values != 0)
        rel_diff[mask_zero_actual_nonzero_sim] = np.inf
        finite_rel_diff = rel_diff[np.isfinite(rel_diff)]
        
        mean_rel_err = np.mean(finite_rel_diff) if len(finite_rel_diff) > 0 else np.nan
        median_rel_err = np.median(finite_rel_diff) if len(finite_rel_diff) > 0 else np.nan

        comparison_metrics = {
            "mean_absolute_error": np.mean(abs_diff),
            "max_absolute_error": np.max(abs_diff),
            "median_absolute_error": np.median(abs_diff),
            "mean_relative_error": mean_rel_err,
            "median_relative_error": median_rel_err,
        }
        results[minute_dt] = comparison_metrics

        # Сохранение детальных квантилей
        if base_output_dir_for_minute_details:
            minute_dt_str = minute_dt.strftime("%Y-%m-%d_%H-%M-%S")
            current_minute_output_dir = os.path.join(base_output_dir_for_minute_details, minute_dt_str)
            os.makedirs(current_minute_output_dir, exist_ok=True)

            simulated_path = os.path.join(current_minute_output_dir, "simulated_quantiles.txt")
            actual_path = os.path.join(current_minute_output_dir, "actual_quantiles.txt")

            with open(simulated_path, 'w', encoding='utf-8') as f_sim:
                for q_val in simulated_minute_quantiles_values:
                    f_sim.write(f"{q_val}\n")
            
            with open(actual_path, 'w', encoding='utf-8') as f_act:
                for q_val in merged_minute_actual_quantiles:
                    f_act.write(f"{q_val}\n")

    return results

if __name__ == '__main__':
    seconds_csv_path_full = "/home/lilclown/study/statshouse/datas/rpc_proxy_rpc_response_ok/rpc_proxy_rpc_response_ok-1s.csv"
    minutes_csv_path_full = "/home/lilclown/study/statshouse/datas/rpc_proxy_rpc_response_ok/rpc_proxy_rpc_response_ok-1m.csv"

    try:
        dataset_name = os.path.basename(os.path.dirname(seconds_csv_path_full))
        if not dataset_name or dataset_name == "datas":
             dataset_name = os.path.basename(seconds_csv_path_full).split('-')[0]
    except Exception:
        dataset_name = "unknown_dataset"

    base_output_dir = os.path.join("testResults", dataset_name)
    minute_details_output_dir = os.path.join(base_output_dir, "minutes")
    summary_file_path = os.path.join(base_output_dir, "minutesTotal.txt")
    os.makedirs(minute_details_output_dir, exist_ok=True)

    aligned_secs_map, aligned_mins_map, common_minute_keys_sorted = align_data(seconds_csv_path_full, minutes_csv_path_full)

    if not common_minute_keys_sorted:
        print("Общих минут не найдено.")
    else:
        SAMPLES_PER_SECOND = 200
        comparison_results = run_simulation_and_comparison(
            aligned_secs_map,
            aligned_mins_map,
            common_minute_keys_sorted,
            samples_per_second_tick=SAMPLES_PER_SECOND,
            base_output_dir_for_minute_details=minute_details_output_dir
        )

        with open(summary_file_path, 'w', encoding='utf-8') as f_summary:
            f_summary.write("MinuteDateTime;MeanAbsoluteError;MaxAbsoluteError;MedianAbsoluteError;MeanRelativeError;MedianRelativeError\n")
            for minute_dt_key, metrics in comparison_results.items():
                f_summary.write(f"{minute_dt_key.strftime('%Y-%m-%d %H:%M:%S')};")
                if "error" in metrics:
                    f_summary.write(f"ERROR: {metrics['error']}\n")
                else:
                    f_summary.write(f"{metrics['mean_absolute_error']:.4f};")
                    f_summary.write(f"{metrics['max_absolute_error']:.4f};")
                    f_summary.write(f"{metrics['median_absolute_error']:.4f};")
                    
                    mre_str = f"{metrics['mean_relative_error']:.4%}" if not np.isnan(metrics['mean_relative_error']) else "NaN"
                    medre_str = f"{metrics['median_relative_error']:.4%}" if not np.isnan(metrics['median_relative_error']) else "NaN"
                    
                    f_summary.write(f"{mre_str};")
                    f_summary.write(f"{medre_str}\n")
        print("Конец")