#!/bin/bash

# Create a file of 5MB made up of 5 chunks of 1MB each
chunk_size=$((1 * 1024 * 1024))  # 1MB
num_chunks=5
filename="test_5chunks.dat"

# Clear the file if it exists
> "$filename"

# Generate each chunk
for i in $(seq 0 $((num_chunks - 1))); do
  char=$(printf "\\$(printf '%03o' $((65 + i)))")  # 'A' + i

  tmp_chunk=$(mktemp)
  
  # Fill chunk with pattern (char), with slight variation every 100 bytes
  for j in $(seq 0 $((chunk_size - 1))); do
    if (( j % 100 == 0 )); then
      variation=$(printf "\\$(printf '%03o' $((65 + i + (j % 26)) ))")
      printf "$variation" >> "$tmp_chunk"
    else
      printf "$char" >> "$tmp_chunk"
    fi
  done

  cat "$tmp_chunk" >> "$filename"
  rm "$tmp_chunk"
done

echo "Created file $filename with size $((chunk_size * num_chunks)) bytes"
echo "Each chunk is $chunk_size bytes with varying ASCII patterns (A-E)"
