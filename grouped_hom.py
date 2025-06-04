import csv
import ast
import numpy as np
from scipy.interpolate import interp1d
from scipy.stats import kruskal
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

    if len(_p_levels_final) < 2:
         if len(_p_levels_final) == 1:
            return lambda x: np.full_like(x, _quantile_values_final[0], dtype=float)

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

def parse_quantile_data_from_csv(filepath):
    quantile_arrays = []
    expected_quantile_len = len(P_LEVELS)
    print("Парсинг файла...")
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            reader = csv.reader(f)
            for i, row in enumerate(reader):
                if not row or len(row) < 2:
                    continue
                
                quantiles_str = row[1].strip()
                try:
                    quantiles_list_from_str = ast.literal_eval(quantiles_str)
                except (SyntaxError, ValueError):
                    if not quantiles_str.startswith('['): quantiles_str = '[' + quantiles_str
                    if not quantiles_str.endswith(']'): quantiles_str = quantiles_str + ']'
                    try:
                        quantiles_list_from_str = ast.literal_eval(quantiles_str)
                    except Exception:
                        continue
                
                try:
                    quantiles_array = np.array(quantiles_list_from_str, dtype=float)
                    if len(quantiles_array) != expected_quantile_len:
                        continue
                    quantile_arrays.append(quantiles_array)
                except Exception:
                    continue
    except FileNotFoundError:
        return []
    except Exception:
        return []

    return quantile_arrays

def perform_homogeneity_analysis_for_chunk(quantile_chunk, chunk_identifier, output_file_handle, num_samples_per_distribution=500, alpha=0.05):
    if not quantile_chunk:
        message = f"Группа {chunk_identifier}: Нет данных."
        output_file_handle.write(message + "\n\n")
        return

    if len(quantile_chunk) < 2:
        message = f"Группа {chunk_identifier}: Недостаточно данных для теста."
        output_file_handle.write(message + "\n\n")
        return

    all_generated_samples = []
    for i, quantiles_arr in enumerate(quantile_chunk):
        try:
            q_func = create_quantile_function(quantiles_arr)
            samples = generate_samples(q_func, num_samples_per_distribution)
            all_generated_samples.append(samples)
        except Exception:
            continue
    
    if len(all_generated_samples) < 2:
        message = f"Группа {chunk_identifier}: Недостаточно выборок для теста."
        output_file_handle.write(message + "\n\n")
        return

    try:
        h_statistic, p_value = kruskal(*all_generated_samples)
        
        if p_value < alpha:
            conclusion = f"Распределения НЕ однородны (p={p_value:.4g})"
        else:
            conclusion = f"Распределения однородны (p={p_value:.4g})"

        output_file_handle.write(f"--- {chunk_identifier} ---\n")
        output_file_handle.write(f"Распределений: {len(all_generated_samples)}\n")
        output_file_handle.write(f"H-статистика: {h_statistic}\n")
        output_file_handle.write(f"P-value: {p_value}\n")
        output_file_handle.write(f"Вывод: {conclusion}\n\n")

    except Exception as e:
        output_file_handle.write(f"--- {chunk_identifier} ---\n")
        output_file_handle.write(f"ОШИБКА: {e}\n\n")

if __name__ == '__main__':
    INPUT_CSV_FILE = "/home/lilclown/study/statshouse/datas/rpc_proxy_rpc_response_ok/rpc_proxy_rpc_response_ok-1s.csv"
    GROUP_SIZE = 3600
    NUM_SAMPLES_PER_DISTRIBUTION = 1000
    ALPHA = 0.05

    if not os.path.exists(INPUT_CSV_FILE):
        exit()

    all_quantile_data_from_file = parse_quantile_data_from_csv(INPUT_CSV_FILE)

    if not all_quantile_data_from_file:
        print("Нет данных для обработки.")
    else:
        output_dir_name = "grouped_homogeneus"
        os.makedirs(output_dir_name, exist_ok=True)
        
        base_name = os.path.basename(INPUT_CSV_FILE)
        filename_no_ext = os.path.splitext(base_name)[0]
        results_filename = f"{filename_no_ext}_grouped_homogeneity_results.txt"
        results_filepath = os.path.join(output_dir_name, results_filename)

        print("Тест...")
        with open(results_filepath, 'w', encoding='utf-8') as f_out:
            num_total_seconds = len(all_quantile_data_from_file)
            num_processed_groups = 0
            for i in range(0, num_total_seconds, GROUP_SIZE):
                chunk_data = all_quantile_data_from_file[i : i + GROUP_SIZE]
                start_second_in_file = i
                end_second_in_file = min(i + GROUP_SIZE - 1, num_total_seconds - 1)
                chunk_identifier = f"Группа {num_processed_groups + 1} (записи {start_second_in_file+1}-{end_second_in_file+1})"
                
                perform_homogeneity_analysis_for_chunk(
                    quantile_chunk=chunk_data,
                    chunk_identifier=chunk_identifier,
                    output_file_handle=f_out,
                    num_samples_per_distribution=NUM_SAMPLES_PER_DISTRIBUTION,
                    alpha=ALPHA
                )
                num_processed_groups += 1