import os
import numpy as np

output_folder = "noise"
os.makedirs(output_folder, exist_ok=True)

total_files = 60 * 60
sample_size = 100_000_000 // total_files 

arrays = []  

for group in range(1, 61):
    for sub in range(1, 61):
        mean = np.random.normal(0, 10)
        sigma = np.random.uniform(0.8, 1.2)
        data = np.random.normal(mean, sigma, sample_size)
        arrays.append(data)
        
        filename = os.path.join(output_folder, f"merge_{group}_{sub}.txt")
        sorted_data = np.sort(data)
        np.savetxt(filename, sorted_data, fmt="%.18e")

base_total = np.concatenate(arrays)
sorted_base_total = np.sort(base_total)

noise_size = int(sorted_base_total.size * 0.07)

q90 = np.percentile(sorted_base_total, 90)

mean_noise = q90 + 5
sigma_noise = 1  

noise = np.random.normal(mean_noise, sigma_noise, noise_size)
sorted_noise = np.sort(noise)

noise_file = os.path.join(output_folder, "merge_noise.txt")
np.savetxt(noise_file, sorted_noise, fmt="%.18e")

total_all = np.concatenate([base_total, noise])
sorted_total_all = np.sort(total_all)

total_file = os.path.join(output_folder, "total.txt")
print("Создаю файл:", os.path.abspath(total_file))
np.savetxt(total_file, sorted_total_all, fmt="%.18e")