with open('queryOrigin.txt', 'r') as file:
    lines = file.readlines()

result = []

for line in lines:
    time_diff_us = int(line)
    time_diff_ms = time_diff_us / 1000.0  # Convert microseconds to milliseconds
    result.append('{:.4f}'.format(time_diff_ms))  # Include decimal point with 3 decimal places

with open('queryResult.txt', 'w') as file:
    file.write('\n'.join(result))
