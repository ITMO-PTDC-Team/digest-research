import matplotlib.pyplot as plt


X = ['0.50',  '0.75',  '0.80',  '0.85',  '0.90',  '0.95',  '0.99',  '0.995',  '0.999'  ]
Y = [0.0001,  0.0004,  0.0006,  0.0009,  0.0017,  0.0077,  0.0524,  0.0822,  0.1602]

plt.bar(X, Y)
plt.ylim(0,0.17)

plt.xlabel('Quantiles')
plt.ylabel('Relative Error %')
plt.title('Original')


plt.savefig('histogram_original.png')

plt.show()