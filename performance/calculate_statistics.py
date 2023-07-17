import numpy as np

data = []
with open('invokeResult.txt', 'r') as file:
    for line in file:
        value = float(line.strip())
        data.append(value)

# Convert the data to a NumPy array
data_arr = np.array(data)

# Calculate the values
minimum = np.min(data_arr)
maximum = np.max(data_arr)
percentile_25 = np.percentile(data_arr, 25)
average = np.mean(data_arr)
percentile_75 = np.percentile(data_arr, 75)

# Write the results to a file
with open('invoke_statistic_result.txt', 'w') as file:
    file.write(f"Minimum: {minimum}\n")
    file.write(f"25th Percentile: {percentile_25}\n")
    file.write(f"Average: {average}\n")
    file.write(f"75th Percentile: {percentile_75}\n")
    file.write(f"Maximum: {maximum}\n")
