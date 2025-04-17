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

output_file = "distributions/test_distribution_caushy.txt"
heavy_tail_data = np.random.standard_cauchy(20_000_000)
sorted_heavy_tail = np.sort(heavy_tail_data)

with open(output_file, "w") as f:
    for value in heavy_tail_data:
        f.write(f"{value}\n")


parameters_bimodal_exp = [
    [5.0, 1.0, 0.5, 20_000_000], 
    [10.0, 2.0, 0.4, 15_000_000],
    [3.0, 0.5, 0.6, 25_000_000],
    [7.0, 1.5, 0.3, 30_000_000],
    [2.0, 0.8, 0.5, 20_000_000]
]

for i, param in enumerate(parameters_bimodal_exp, start=0):
    output_file = f"distributions/test_distribution_{i+11}.txt"
    shift, scale, proportion, num_elements = param
    num_left = int(num_elements * proportion)
    num_right = num_elements - num_left
    left_data = -np.random.exponential(scale=scale, size=num_left) - shift
    right_data = np.random.exponential(scale=scale, size=num_right) + shift
    combined = np.concatenate((left_data, right_data))
    sorted_data = np.sort(combined)
    with open(output_file, "w") as f:
        for value in sorted_data:
            f.write(f"{value}\n")

    
mu1 = 10    
sigma1 = 0.5 
mu2 = 0     
sigma2 = 0.5

data1 = np.random.normal(mu1, sigma1, 10000000)
data2 = np.random.normal(mu2, sigma2, 10000000)
output_file1 = f"distributions/test_distribution_merge_low.txt"
output_file2 = f"distributions/test_distribution_merge_hight.txt"

sorted_data1 = np.sort(data1)
sorted_data2 = np.sort(data2)


with open(output_file1, "w") as f:
        for value in sorted_data1:
            f.write(f"{value}\n")

with open(output_file2, "w") as f:
        for value in sorted_data2:
            f.write(f"{value}\n")


merged_data = np.concatenate((data1, data2))
sorted_merged_data = np.sort(merged_data)
output_file3 = f"distributions/test_distribution_merged.txt"
with open(output_file3, "w") as f:
        for value in sorted_merged_data:
            f.write(f"{value}\n")