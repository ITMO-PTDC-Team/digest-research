import os
import numpy as np

output_folder = "hour"
os.makedirs(output_folder, exist_ok=True)

total_files = 60 * 60
sample_size = 100_000_000 // total_files

total = []

for group in range(1, 61):
    for sub in range(1, 61):
        mean = np.random.normal(0, 10)
        sigma = np.random.uniform(0.8, 1.2)
        data = np.random.normal(mean, sigma, sample_size)
        total.append(data)
        
        filename = os.path.join(output_folder, f"merge_{group}_{sub}.txt")
        sorted_data = np.sort(data)
        np.savetxt(filename, sorted_data, fmt="%.18e")

total = np.concatenate(total)
sorted_total = np.sort(total)

total_file = os.path.join(output_folder, "total.txt")
print("Создаю файл:", os.path.abspath(total_file))
np.savetxt(total_file, sorted_total, fmt="%.18e")