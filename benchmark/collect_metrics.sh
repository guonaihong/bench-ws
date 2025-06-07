#!/bin/bash

# Add timestamp function
log_with_timestamp() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*"
}

# Add signal handler
cleanup() {
    log_with_timestamp "Received interrupt signal, cleaning up..."
    exit 0
}

# Set signal handler
trap cleanup SIGINT SIGTERM

# Performance metrics collection script
# Usage: ./collect_metrics.sh <framework> <pid> <duration> <output_file> <benchmark_pid>

FRAMEWORK=$1
PID=$2
DURATION=$3
OUTPUT_FILE=$4
BENCHMARK_PID=$5  # Added: benchmark script PID for file naming

# Create table header if output file doesn't exist
if [ ! -f "$OUTPUT_FILE" ]; then
    echo "| Framework | TPS(Start) | TPS(Middle) | TPS(End) | CPU(Start) | CPU(Middle) | CPU(End) | Memory(Start) | Memory(Middle) | Memory(End) | Threads(Start) | Threads(Middle) | Threads(End) | FD(Start) | FD(Middle) | FD(End) |" > "$OUTPUT_FILE"
    echo "|-----------|------------|-------------|----------|------------|-------------|----------|---------------|----------------|-------------|---------------|----------------|--------------|-----------|------------|----------|" >> "$OUTPUT_FILE"
fi

log_with_timestamp "=================== Starting metrics collection for $FRAMEWORK (PID: $PID) ==================="

# Initialize array to store CPU and memory data
cpu_values=()
mem_values=()
thread_values=()  # Added: store thread count data
fd_values=()      # Added: store file descriptor count data

# Sampling interval (seconds)
SAMPLE_INTERVAL=1
total_samples=$((DURATION / SAMPLE_INTERVAL))

log_with_timestamp "Will collect $total_samples samples every $SAMPLE_INTERVAL seconds for $DURATION seconds"

# Collect data continuously
for ((i=1; i<=total_samples; i++)); do
    # Use top command to collect CPU and memory usage
    # Detect operating system type to use the correct top command format
    OS_TYPE=$(uname -s)
    
    if [ "$OS_TYPE" = "Darwin" ]; then
        # macOS system
        top_output=$(top -pid $PID -l 1 -stats pid,cpu,mem 2>/dev/null | tail -n +13 | head -1)
        if [ -n "$top_output" ]; then
            current_cpu=$(echo "$top_output" | awk '{print $2}' | sed 's/%//')
            current_mem_str=$(echo "$top_output" | awk '{print $3}')
            
            # On macOS, parse memory format (e.g., 3904K, 256M, 1.2G)
            if echo "$current_mem_str" | grep -q "K"; then
                current_mem_kb=$(echo "$current_mem_str" | sed 's/K.*//')
                current_mem_mb=$((current_mem_kb / 1024))
            elif echo "$current_mem_str" | grep -q "M"; then
                current_mem_mb=$(echo "$current_mem_str" | sed 's/M.*//' | cut -d. -f1)
            elif echo "$current_mem_str" | grep -q "G"; then
                current_mem_gb=$(echo "$current_mem_str" | sed 's/G.*//')
                current_mem_mb=$(echo "$current_mem_gb * 1024" | bc)
            else
                # If it's just a number, assume it's KB
                current_mem_mb=$(echo "$current_mem_str / 1024" | bc 2>/dev/null || echo "0")
            fi
        else
            # If top command fails, fall back to ps command
            current_cpu=$(ps -p $PID -o %cpu --no-headers 2>/dev/null | tr -d ' ')
            current_mem_kb=$(ps -p $PID -o rss --no-headers 2>/dev/null | tr -d ' ')
            current_mem_mb=$((current_mem_kb / 1024))
        fi
    else
        # Linux system
        # Use a more precise top command to ensure getting the correct process line
        top_output=$(top -p $PID -n 1 -b 2>/dev/null | grep "^ *$PID " | head -1)
        
        if [ -n "$top_output" ]; then
            # Linux top output format is usually: PID USER PR NI VIRT RES SHR S %CPU %MEM TIME+ COMMAND
            current_cpu=$(echo "$top_output" | awk '{print $9}')
            
            # Directly use RES field (actual memory usage, in KB), more accurate
            current_mem_res=$(echo "$top_output" | awk '{print $6}')
            if [ -n "$current_mem_res" ] && [ "$current_mem_res" -gt 0 ]; then
                current_mem_mb=$((current_mem_res / 1024))
            else
                current_mem_mb="0"
            fi
        else
            # If top command fails, fall back to ps command
            current_cpu=$(ps -p $PID -o %cpu --no-headers 2>/dev/null | tr -d ' ')
            current_mem_kb=$(ps -p $PID -o rss --no-headers 2>/dev/null | tr -d ' ')
            current_mem_mb=$((current_mem_kb / 1024))
        fi
    fi
    
    # Check if process still exists
    if [ -z "$current_cpu" ] || [ -z "$current_mem_mb" ]; then
        log_with_timestamp "Warning: Process $PID not found at sample $i, stopping collection"
        break
    fi
    
    # Get thread count
    if [ "$OS_TYPE" = "Darwin" ]; then
        # macOS: Use ps command to get thread count
        thread_count=$(ps -M -p $PID 2>/dev/null | wc -l | tr -d ' ')
        # ps -M includes header, so need to subtract 1
        thread_count=$((thread_count - 1))
    else
        # Linux: Use /proc/PID/stat to get thread count, or ps -eLf
        if [ -f "/proc/$PID/status" ]; then
            thread_count=$(grep "^Threads:" /proc/$PID/status | awk '{print $2}')
        else
            # Backup method: Use ps command
            thread_count=$(ps -eLf | grep "^ *[^ ]* *$PID " | wc -l | tr -d ' ')
        fi
    fi
    
    # Ensure thread count is valid
    if ! [[ "$thread_count" =~ ^[0-9]+$ ]] || [ "$thread_count" -le 0 ]; then
        thread_count="1"  # Default at least 1 thread
    fi
    
    # Ensure values are valid
    if ! [[ "$current_cpu" =~ ^[0-9]+\.?[0-9]*$ ]]; then
        current_cpu="0.0"
    fi
    if ! [[ "$current_mem_mb" =~ ^[0-9]+$ ]]; then
        current_mem_mb="0"
    fi
    
    # Get file descriptor count
    if [ "$OS_TYPE" = "Darwin" ]; then
        # macOS: Use lsof to count file descriptors
        fd_count=$(lsof -p $PID 2>/dev/null | wc -l)
        # lsof includes header, so subtract 1
        fd_count=$((fd_count - 1))
    else
        # Linux: Use /proc/PID/fd to count file descriptors
        if [ -d "/proc/$PID/fd" ]; then
            fd_count=$(ls -l /proc/$PID/fd 2>/dev/null | wc -l)
            # ls -l includes total line, so subtract 1
            fd_count=$((fd_count - 1))
        else
            fd_count=0
        fi
    fi
    
    # Ensure fd count is valid
    if ! [[ "$fd_count" =~ ^[0-9]+$ ]]; then
        fd_count=0
    fi
    
    # Store data
    cpu_values+=("$current_cpu")
    mem_values+=("$current_mem_mb")
    thread_values+=("$thread_count")
    fd_values+=("$fd_count")
    
    log_with_timestamp "Sample $i/$total_samples: CPU=${current_cpu}%, Memory=${current_mem_mb}MB, Threads=${thread_count}, FD=${fd_count} (via top)"
    
    # Wait for next sample (except for the last one)
    if [ $i -lt $total_samples ]; then
        sleep $SAMPLE_INTERVAL
    fi
done

# Calculate CPU statistics
if [ ${#cpu_values[@]} -gt 0 ]; then
    # Calculate max
    max_cpu=${cpu_values[0]}
    for cpu in "${cpu_values[@]}"; do
        if (( $(echo "$cpu > $max_cpu" | bc -l) )); then
            max_cpu=$cpu
        fi
    done
    
    # Calculate min
    min_cpu=${cpu_values[0]}
    for cpu in "${cpu_values[@]}"; do
        if (( $(echo "$cpu < $min_cpu" | bc -l) )); then
            min_cpu=$cpu
        fi
    done
    
    # Calculate avg
    sum_cpu=0
    for cpu in "${cpu_values[@]}"; do
        sum_cpu=$(echo "$sum_cpu + $cpu" | bc -l)
    done
    avg_cpu=$(echo "scale=1; $sum_cpu / ${#cpu_values[@]}" | bc -l)
else
    max_cpu="N/A"
    min_cpu="N/A"
    avg_cpu="N/A"
fi

# Calculate memory statistics
if [ ${#mem_values[@]} -gt 0 ]; then
    # Calculate max
    max_mem=${mem_values[0]}
    for mem in "${mem_values[@]}"; do
        if [ $mem -gt $max_mem ]; then
            max_mem=$mem
        fi
    done
    
    # Calculate min
    min_mem=${mem_values[0]}
    for mem in "${mem_values[@]}"; do
        if [ $mem -lt $min_mem ]; then
            min_mem=$mem
        fi
    done
    
    # Calculate avg
    sum_mem=0
    for mem in "${mem_values[@]}"; do
        sum_mem=$((sum_mem + mem))
    done
    avg_mem=$((sum_mem / ${#mem_values[@]}))
else
    max_mem="N/A"
    min_mem="N/A"
    avg_mem="N/A"
fi

log_with_timestamp "CPU Statistics: Max=${max_cpu}%, Min=${min_cpu}%, Avg=${avg_cpu}%"
log_with_timestamp "Memory Statistics: Max=${max_mem}MB, Min=${min_mem}MB, Avg=${avg_mem}MB"

# Calculate thread count statistics
if [ ${#thread_values[@]} -gt 0 ]; then
    # Calculate max
    max_threads=${thread_values[0]}
    for threads in "${thread_values[@]}"; do
        if [ $threads -gt $max_threads ]; then
            max_threads=$threads
        fi
    done
    
    # Calculate min
    min_threads=${thread_values[0]}
    for threads in "${thread_values[@]}"; do
        if [ $threads -lt $min_threads ]; then
            min_threads=$threads
        fi
    done
    
    # Calculate avg
    sum_threads=0
    for threads in "${thread_values[@]}"; do
        sum_threads=$((sum_threads + threads))
    done
    avg_threads=$((sum_threads / ${#thread_values[@]}))
else
    max_threads="N/A"
    min_threads="N/A"
    avg_threads="N/A"
fi

log_with_timestamp "Thread Statistics: Max=${max_threads}, Min=${min_threads}, Avg=${avg_threads}"

# Calculate FD statistics
if [ ${#fd_values[@]} -gt 0 ]; then
    # Calculate max
    max_fd=${fd_values[0]}
    for fd in "${fd_values[@]}"; do
        if [ $fd -gt $max_fd ]; then
            max_fd=$fd
        fi
    done
    
    # Calculate min
    min_fd=${fd_values[0]}
    for fd in "${fd_values[@]}"; do
        if [ $fd -lt $min_fd ]; then
            min_fd=$fd
        fi
    done
    
    # Calculate avg
    sum_fd=0
    for fd in "${fd_values[@]}"; do
        sum_fd=$((sum_fd + fd))
    done
    avg_fd=$((sum_fd / ${#fd_values[@]}))
else
    max_fd="N/A"
    min_fd="N/A"
    avg_fd="N/A"
fi

log_with_timestamp "FD Statistics: Max=${max_fd}, Min=${min_fd}, Avg=${avg_fd}"

# Extract data from corresponding framework's TPS file
TPS_FILE="$(dirname "$OUTPUT_FILE")/output/$FRAMEWORK.$BENCHMARK_PID.tps"
start_tps="N/A"
mid_tps="N/A"
end_tps="N/A"

# Wait for TPS file to be generated, up to 5 retries
log_with_timestamp "Waiting for TPS data file: $TPS_FILE"
retry_count=0
max_retries=5

while [ $retry_count -lt $max_retries ]; do
    if [ -f "$TPS_FILE" ] && [ -s "$TPS_FILE" ]; then
        log_with_timestamp "TPS file found and not empty"
        break
    fi
    
    retry_count=$((retry_count + 1))
    log_with_timestamp "TPS file not ready, retry $retry_count/$max_retries..."
    sleep 2
done

if [ -f "$TPS_FILE" ] && [ -s "$TPS_FILE" ]; then
    start_tps=$(grep "^start" "$TPS_FILE" | awk '{print $2}' | head -1)
    mid_tps=$(grep "^middle" "$TPS_FILE" | awk '{print $2}' | head -1)
    end_tps=$(grep "^end" "$TPS_FILE" | awk '{print $2}' | head -1)
    log_with_timestamp "TPS data extracted: Start=$start_tps, Middle=$mid_tps, End=$end_tps"
else
    log_with_timestamp "Warning: TPS file not found or empty after $max_retries retries: $TPS_FILE"
    if [ -f "$TPS_FILE" ]; then
        log_with_timestamp "TPS file content: $(cat "$TPS_FILE")"
    fi
fi

# Append data to markdown table - using new statistics format
echo "| $FRAMEWORK | $start_tps | $mid_tps | $end_tps | $max_cpu% | $min_cpu% | $avg_cpu% | ${max_mem}MB | ${min_mem}MB | ${avg_mem}MB | $max_threads | $min_threads | $avg_threads | $max_fd | $min_fd | $avg_fd |" >> "$OUTPUT_FILE"

# Save CPU, memory, thread count and FD data to separate file (also use PID suffix)
CPU_FILE="$(dirname "$OUTPUT_FILE")/output/$FRAMEWORK.$BENCHMARK_PID.cpu"
MEM_FILE="$(dirname "$OUTPUT_FILE")/output/$FRAMEWORK.$BENCHMARK_PID.mem"
THREAD_FILE="$(dirname "$OUTPUT_FILE")/output/$FRAMEWORK.$BENCHMARK_PID.threads"
FD_FILE="$(dirname "$OUTPUT_FILE")/output/$FRAMEWORK.$BENCHMARK_PID.fd"

echo "max $max_cpu" > "$CPU_FILE"
echo "min $min_cpu" >> "$CPU_FILE"
echo "avg $avg_cpu" >> "$CPU_FILE"

echo "max $max_mem" > "$MEM_FILE"
echo "min $min_mem" >> "$MEM_FILE"
echo "avg $avg_mem" >> "$MEM_FILE"

echo "max $max_threads" > "$THREAD_FILE"
echo "min $min_threads" >> "$THREAD_FILE"
echo "avg $avg_threads" >> "$THREAD_FILE"

echo "max $max_fd" > "$FD_FILE"
echo "min $min_fd" >> "$FD_FILE"
echo "avg $avg_fd" >> "$FD_FILE"

log_with_timestamp "Metrics collection completed for $FRAMEWORK" 