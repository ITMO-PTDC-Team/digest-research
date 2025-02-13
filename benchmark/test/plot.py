import sys
import numpy as np
import matplotlib.pyplot as plt

def read_floats_from_file(filename):
    with open(filename, 'r') as file:
        return np.array([float(line.strip()) for line in file if line.strip()])
    
class CDF:
    def __init__(self, data):
        self.sorted_data = np.sort(data)
        self.n = len(self.sorted_data)

    def __call__(self, x):
        return np.searchsorted(self.sorted_data, x, side='right') / self.n

def cdf_to_pdf(cdf, x_values):
    cdf_values = np.array([cdf(x) for x in x_values])
    pdf_values = np.diff(cdf_values) / np.diff(x_values)
    return pdf_values

def plot_cdf(cdf,x_values):
    y_values = [cdf(x) for x in x_values]

    plt.plot(x_values, y_values, label="CDF", color='red')
    plt.xlabel('x')
    plt.ylabel('CDF(x)')
    plt.grid(True)
    plt.legend()
    plt.show()

def plot_pdf(x_values, pdf_values):
    plt.plot(x_values[:-1], pdf_values, label="PDF", color='blue')
    plt.xlabel('x')
    plt.ylabel('PDF(x)')
    plt.grid(True)
    plt.legend()
    plt.show()
def plot(cdf,x_values,pdf_values):
    y_values = [cdf(x) for x in x_values]

    plt.plot(x_values, y_values, label="CDF", color='red')
    plt.plot(x_values[:-1], pdf_values, label="PDF", color='blue')
    plt.xlabel('x')
    plt.ylabel('CDF(x)')
    plt.grid(True)
    plt.legend()
    plt.show()
    
def main():
    if len(sys.argv) > 1:
        filename = sys.argv[1]
    else:
        filename = 'input.txt'

    try:
        data = read_floats_from_file(filename)
    except FileNotFoundError:
        return
    cdf = CDF(data)

    data_min = np.min(data)
    data_max = np.max(data)
    x_values = np.linspace(data_min-1, data_max, 1000000)

    pdf_values = cdf_to_pdf(cdf, x_values)
    plot(cdf,x_values,pdf_values)

if __name__ == "__main__":
    main()