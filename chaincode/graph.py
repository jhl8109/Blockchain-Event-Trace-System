import matplotlib.pyplot as plt

# Your data
data = {
    "Transaction Average Execution Time": 2323.385,
    "Perform transaction + store off-chain": 2326.293,
    "Block (world state) query": 29.818,
    "Off-chain queries": 0.328,
    "Block Queries - Off-Chain Queries": 29.49,
    "Off-chain storage": 2.908,
    "Time Gain": 26.582
}

# Extracting data labels and values
labels = list(data.keys())
values = list(data.values())

# Creating the histogram
plt.bar(labels, values)
plt.xlabel("Categories")
plt.ylabel("Time (ms)")
plt.title("Performance Comparison Histogram")
plt.xticks(rotation=45, ha="right")  # Rotate the x-axis labels for better readability
plt.tight_layout()  # Adjust the layout for better readability

# Show the histogram
plt.show()