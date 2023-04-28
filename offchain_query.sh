#!/bin/bash

file="query_performance.txt" # Replace "data.txt" with the name of your data file

# Initialize variables
min=1000000
max=0
sum=0
count=0

# Read data from the file line by line
while IFS= read -r line; do
  # Update min and max
  if (( $(echo "$line < $min" | bc -l) )); then
    min="$line"
  fi
  if (( $(echo "$line > $max" | bc -l) )); then
    max="$line"
  fi

  # Calculate sum and count
  sum=$(echo "$sum + $line" | bc -l)
  count=$((count + 1))
done < "$file"

# Calculate average
avg=$(echo "$sum / $count" | bc -l)

# Print results
printf "Min: %.3f\nMax: %.3f\nAvg: %.3f\n" "$min" "$max" "$avg"
