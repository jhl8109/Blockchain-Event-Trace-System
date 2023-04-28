#!/bin/bash

awk '
BEGIN {
  min1 = min2 = max1 = max2 = avg1 = avg2 = count = 0
}

{
  if (NR == 1) {
    min1 = max1 = $1
    min2 = max2 = $2
  }

  if ($1 < min1) min1 = $1
  if ($2 < min2) min2 = $2

  if ($1 > max1) max1 = $1
  if ($2 > max2) max2 = $2

  avg1 += $1
  avg2 += $2

  count++
}

END {
  printf "Column 1: Min = %.3f, Max = %.3f, Avg = %.3f\n", min1, max1, avg1 / count
  printf "Column 2: Min = %.3f, Max = %.3f, Avg = %.3f\n", min2, max2, avg2 / count
}
' txquery.txt
