import sys
import numpy as np
import matplotlib.pyplot as plt
def cdf_to_pdf(cdf, x_values):
    cdf_values = np.array([cdf(x) for x in x_values])
    pdf_values = np.diff(cdf_values) / np.diff(x_values)
    return pdf_values

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

    plt.plot(x_values, y_values, label="dist", color='green')
    # plt.xlabel('x')
    # plt.ylabel('CDF(x)')
    # plt.grid(True)
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
    # plt.plot(x_values[:-1], pdf_values, label="PDF", color='blue')
    plt.xlabel('x')
    plt.ylabel('CDF(x)')
    plt.grid(True)
    plt.legend()
    # plt.show()


def read_data(filename):
    with open(filename, 'r') as file:
        try:
            N, M = map(int, file.readline().split())
        except ValueError as e:
            exit(1)
        try:
            X = list(map(float, file.readline().split()))
        except ValueError as e:
            exit(1)
        if len(X) != N:
            exit(1)
        
        Y = []
        for i in range(M - 1):
            try:
                y_values = list(map(float, file.readline().split()))
            except ValueError as e:
                exit(1)
            if len(y_values) != N:
                exit(1)
            
            Y.append(y_values)
        
        if not Y:
            exit(1)
        
        return X, Y
    
for j in range (0,33):
    print (j)
    data = read_floats_from_file("distributions/test_distribution_"+str(j)+".txt")  
    cdf = CDF(data)
    filename="sin_pow_2_5"
    Y1, X1 = read_data('quantiles/td_sin_pow_2_5.txt')
    Y2, X2 = read_data('quantiles/cdf_sin_pow_2_5.txt')
    y_values = [cdf(x) for x in X1[j]]
    # data_min = np.min(data)
    # data_max = np.max(data)
    # x_values = np.linspace(data_min, data_max, 1000000)
    pdf_values = cdf_to_pdf(cdf, X1[j])
    max_v=max(pdf_values)
    pdf_norm=[x/max_v for x in pdf_values]
    plt.plot(X1[j][:-1] , pdf_norm,label='Distribution',color='green')
    plt.plot(X1[j], Y1, label='Tdigest', color='blue', linestyle='-')
    plt.plot(X2[j], Y2, label='CDF', color='red', linestyle='-')
    plt.xlabel('X')
    plt.ylabel('Y')
    plt.title('Сравнение распределения, Tdigest и CDF')
    plt.legend()
    plt.grid(True)
    plt.tight_layout()
    plt.savefig(f'Graphic_{j}.png')
    plt.show()  
