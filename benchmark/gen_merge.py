import os
import numpy as np

output_folder = "merges"
os.makedirs(output_folder, exist_ok=True)

total = []

for i in range(1, 61):
    mean = np.random.normal(0, 10)
    sigma = np.abs(np.random.normal(1, 0.2))
    data = np.random.normal(mean, sigma, 1_000_000)
    total.append(data)
    
    merge_file = os.path.join(output_folder, f"merge{i}.txt")
    sorted_data = np.sort(data)
    np.savetxt(merge_file, sorted_data, fmt="%.18e")

total = np.concatenate(total)
sorted_total = np.sort(total)

total_file = os.path.join(output_folder, "total.txt")
print("Создаю файл:", os.path.abspath(total_file))
np.savetxt(total_file, sorted_total, fmt="%.18e")