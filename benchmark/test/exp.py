import numpy as np
import matplotlib.pyplot as plt


def exp(x):
    c = 4.0
    return (np.exp(c * x) - 1.0) / (np.exp(c) - 1.0)



k_values = np.linspace(0, 1, 1000)  


y_values = exp(k_values)


plt.plot(k_values, y_values, label=f"exp(x)")
plt.xlabel("x")
plt.ylabel("exp(x)")
plt.title("exp(x)")
plt.legend()
plt.grid(True)
plt.show()