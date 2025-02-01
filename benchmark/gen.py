import numpy as np

zipf_param = 2  
num_elements = 20_000_000 
noise_mean = 0  
noise_std = 1  
output_file = "distribution.txt" 
num_proportion = 0.7
num_zipf_elements = int(num_elements * num_proportion)
num_noise_elements = num_elements - num_zipf_elements

zipf_data = np.random.zipf(zipf_param, num_zipf_elements)

white_noise = np.random.normal(noise_mean, noise_std, num_noise_elements)


combined_data = np.concatenate((zipf_data, white_noise))

sorted_data = np.sort(combined_data)

with open(output_file, "w") as f:
    for value in sorted_data:
        f.write(f"{value}\n")
