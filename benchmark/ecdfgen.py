import numpy as np
import os
import argparse

def ecdf(data):
    x = np.sort(data)
    y = np.arange(1, len(x) + 1) / len(x)
    return x, y

def generate_from_ecdf(x, y, size=1):
    u = np.random.uniform(0, 1, size)
    return np.interp(u, y, x)

def parse_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument("input_file", type=str)
    return parser.parse_args()

def main():
    args = parse_arguments()
    num_samples = 10000
    num_files = 3600
    if not os.path.exists(args.input_file):
        raise FileNotFoundError(f"File {args.input_file} not found")
    
    data = np.loadtxt(args.input_file)
    x, y = ecdf(data)
    
    os.makedirs("generated", exist_ok=True)
    base = os.path.splitext(os.path.basename(args.input_file))[0]
    
    for i in range(1, num_files + 1):
        samples = generate_from_ecdf(x, y, num_samples)
        output = f"generated/{base}_{i}.txt"
        np.savetxt(output, samples, fmt="%.6f")

if __name__ == "__main__":
    main()