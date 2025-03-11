import matplotlib.pyplot as plt

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
        
        return X, Y[9]

X1, Y1 = read_data('quantiles/100sin_pow1_5.txt')
X2, Y2 = read_data('quantiles/cdf100sin_pow1_5.txt')

plt.plot(X1, Y1, label='Tdigest',color='blue', linestyle='-')
plt.plot(X2, Y2, label='CDF',color='red', linestyle='-')
plt.xlabel('X')
plt.ylabel('Y')
plt.title('Сравнение TDgest и CDF')
plt.legend()
plt.grid(True)
plt.show()