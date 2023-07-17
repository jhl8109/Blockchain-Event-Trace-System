with open('invoke.txt', 'r') as file:
    lines = file.readlines()

result = []

for line in lines:
    data = line.split()
    if len(data) == 2:
        start, end = data
        time_diff_ns = int(end) - int(start)
        time_diff_ms = time_diff_ns / 1_000_000.0  # Convert nanoseconds to milliseconds
        result.append('{:.4f}'.format(time_diff_ms))  # Include decimal point with 6 decimal places

with open('invokeResult.txt', 'w') as file:
    file.write('\n'.join(result))
