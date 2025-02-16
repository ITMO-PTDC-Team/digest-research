import numpy as np
import matplotlib.pyplot as plt

import math

compression=1000.0
def integrated_q_sin(k: float) -> float:

    return (math.sin(min(k, compression) * math.pi / compression - math.pi / 2.0) + 1.0) / 2.0

def integrated_location_sin(q: float) -> float:

    return compression * (math.asin(2.0 * q - 1.0) + math.pi / 2.0) / math.pi

def logarithmic_location(q: float) -> float:
    return compression * math.log(1 + q)

l = 0
r = 1
k_values = np.linspace(l, r, 100000)

q_values = [logarithmic_location(k) for k in k_values]

plt.plot(k_values, q_values, label='q_sin(x)')

plt.xlabel('x')
plt.ylabel('f(x)')
plt.title('График функции f(x)')

plt.legend()

plt.show()