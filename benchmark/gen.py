import numpy as np


parameters_zipf =[
    [1.5,20_000_000,0,1,0.7],
    [2,20_000_000,0,1,0.7],
    [2,30_000_000,0,1,0.8],
    [3.5,30_000_000,0,1,0.7],
    [2,20_000_000,0,1,0.6]
]

max_int32 = np.iinfo(np.int32).max

for i, param in enumerate(parameters_zipf, start=0):
    output_file = f"distributions/test_distribution_{i}.txt"
    zipf_param = parameters_zipf[i][0]
    num_elements = parameters_zipf[i][1]
    noise_mean = parameters_zipf[i][2] 
    noise_std = parameters_zipf[i][3]
    num_proportion = parameters_zipf[i][4]
    num_zipf_elements = int(num_elements * num_proportion)
    num_noise_elements = num_elements - num_zipf_elements
    zipf_data = np.random.zipf(zipf_param, num_zipf_elements)
    filtered_data = zipf_data[zipf_data <= max_int32]
    
    white_noise = np.random.normal(noise_mean, noise_std, num_noise_elements)


    combined_data = np.concatenate((filtered_data, white_noise))

    sorted_data = np.sort(combined_data)

    with open(output_file, "w") as f:
        for value in sorted_data:
            f.write(f"{value}\n")


parameters_normal=[
    [0,1,10_000_000,5,2,10_000_000],
    [0,1,15_000_000,5,2,15_000_000],
    [0,2,10_000_000,5,2,10_000_000],
    [1,1,10_000_000,4,3,10_000_000],
    [0,1,20_000_000,5,2,10_000_000],
    [0,1,10_000_000,5,2,30_000_000]
]
for i, param in enumerate(parameters_zipf, start=0):
    output_file = f"distributions/test_distribution_{i+5}.txt"
    n1_mean=parameters_normal[i][0]
    n1_std=parameters_normal[i][1]
    n1_num=parameters_normal[i][2]
    n2_mean=parameters_normal[i][3]
    n2_std=parameters_normal[i][4]
    n2_num=parameters_normal[i][5]
    n1=np.random.normal(n1_mean,n1_std,n1_num)
    n2=np.random.normal(n2_mean,n2_std,n2_num)
    n=np.concatenate((n1,n2))
    n_sorted=np.sort(n)

    with open(output_file, "w") as f:
        for value in n_sorted:
            f.write(f"{value}\n")