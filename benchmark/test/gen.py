import numpy as np


parameters =[
    [1.5,20_000_000,0,1,0.7],
    [2,20_000_000,0,1,0.7],
    [2,30_000_000,0,1,0.8],
    [3.5,30_000_000,0,1,0.7],
    [2,20_000_000,0,1,0.6],
    [2,30_000_000,0,1,0],
    [2,30_000_000,0,1,0.1],
    [2,30_000_000,0,1,0.05]

]
for i, param in enumerate(parameters, start=0):
    output_file = f"test_distribution_{i}.txt"
    zipf_param = parameters[i][0]
    num_elements = parameters[i][1]
    noise_mean = parameters[i][2] 
    noise_std = parameters[i][3]
    num_proportion = parameters[i][4]
    num_zipf_elements = int(num_elements * num_proportion)
    num_noise_elements = num_elements - num_zipf_elements
    zipf_data = np.random.zipf(zipf_param, num_zipf_elements)
    
    white_noise = np.random.normal(noise_mean, noise_std, num_noise_elements)


    combined_data = np.concatenate((zipf_data, white_noise))

    sorted_data = np.sort(combined_data)

    with open(output_file, "w") as f:
        for value in sorted_data:
            f.write(f"{value}\n")
