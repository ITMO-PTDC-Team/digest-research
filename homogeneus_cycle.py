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

def perform_homogeneity_analysis(csv_filepath, num_samples_per_distribution=500, alpha=0.05):
    filename = os.path.basename(csv_filepath)
    filename_no_ext = os.path.splitext(filename)[0]
    output_dir = os.path.join("homogeneous", filename_no_ext)
    os.makedirs(output_dir, exist_ok=True)
    results_filepath = os.path.join(output_dir, "kruskal_wallis_test_results.txt")

    all_quantile_arrays = parse_quantile_data_from_csv(csv_filepath)

    if not all_quantile_arrays:
        message = "Нет данных для обработки."
        print(message)
        with open(results_filepath, 'w', encoding='utf-8') as f_res:
            f_res.write(message + "\n")
        return

    if len(all_quantile_arrays) < 2:
        message = "Недостаточно данных для теста."
        print(message)
        with open(results_filepath, 'w', encoding='utf-8') as f_res:
            f_res.write(message + "\n")
        return

    print("Тест...")
    all_generated_samples = []
    for i, quantiles_arr in enumerate(all_quantile_arrays):
        try:
            q_func = create_quantile_function(quantiles_arr)
            samples = generate_samples(q_func, num_samples_per_distribution)
            all_generated_samples.append(samples)
        except Exception:
            continue

    if len(all_generated_samples) < 2:
        message = "Недостаточно выборок для теста."
        print(message)
        with open(results_filepath, 'w', encoding='utf-8') as f_res:
            f_res.write(message + "\n")
        return

    try:
        h_statistic, p_value = kruskal(*all_generated_samples)

        if p_value < alpha:
            conclusion = f"Распределения НЕ однородны (p={p_value:.4g})"
        else:
            conclusion = f"Распределения однородны (p={p_value:.4g})"

        with open(results_filepath, 'w', encoding='utf-8') as f_res:
            f_res.write(f"Распределений: {len(all_generated_samples)}\n")
            f_res.write(f"Выборок на распределение: {num_samples_per_distribution}\n")
            f_res.write(f"Уровень значимости: {alpha}\n")
            f_res.write("-" * 50 + "\n")
            f_res.write(f"H-статистика: {h_statistic}\n")
            f_res.write(f"P-value: {p_value}\n")
            f_res.write("-" * 50 + "\n")
            f_res.write(f"Вывод: {conclusion}\n\n")

    except Exception as e:
        with open(results_filepath, 'w', encoding='utf-8') as f_res:
            f_res.write(f"ОШИБКА: {e}\n")

if __name__ == '__main__':
    INPUT_CSV_FILE = "/home/lilclown/study/statshouse/datas/nginx_request_length/nginx_request_length-1s.csv"
    perform_homogeneity_analysis(INPUT_CSV_FILE, num_samples_per_distribution=2000, alpha=0.05)