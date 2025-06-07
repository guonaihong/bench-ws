#!/bin/bash

# Add timestamp function
log_with_timestamp() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*"
}

# Add signal handler
cleanup() {
    log_with_timestamp "Received interrupt signal, cleaning up..."
    # Stop all servers
    "$(dirname "$0")/stop-servers.sh"
    # Kill all child processes
    pkill -P $$ 2>/dev/null || true
    exit 1
}

# Set signal handler
trap cleanup SIGINT SIGTERM

# Add time conversion function
duration_to_seconds() {
    local duration=$1
    local value=$(echo "$duration" | sed 's/[^0-9]*//g')  # Extract numeric part
    local unit=$(echo "$duration" | sed 's/[0-9]*//g')    # Extract unit part
    
    case "$unit" in
        "s"|"")  # Seconds or no unit, default to seconds
            echo "$value"
            ;;
        "m")     # Minutes
            echo $((value * 60))
            ;;
        "h")     # Hours
            echo $((value * 3600))
            ;;
        "d")     # Days
            echo $((value * 86400))
            ;;
        *)
            log_with_timestamp "Warning: Unknown time unit '$unit', assuming seconds"
            echo "$value"
            ;;
    esac
}

# This is a core benchmark script that can be called with different parameters
# Usage: benchmark-core.sh [concurrency] [duration] [rebuild] [extra_args]
# Example: benchmark-core.sh 10000 100s true

# Default values
CONCURRENCY=${1:-10000}  # Default to 10k connections if not specified
DURATION=${2:-100s}      # Default to 100 seconds if not specified
REBUILD=${3:-false}      # Whether to rebuild the project
EXTRA_ARGS=${4:-""}      # Additional arguments to pass to bench-tcp

# Rebuild if requested
if [ "$REBUILD" = "true" ]; then
    log_with_timestamp "Rebuilding project..."
    make clean
    make
    if [ $? -ne 0 ]; then
        log_with_timestamp "Error: Build failed"
        exit 1
    fi
    log_with_timestamp "Build completed successfully"
fi

# Stop any existing servers first
log_with_timestamp "Stopping any existing servers..."
"$(dirname "$0")/stop-servers.sh"

# Source the config file to get ENABLED_SERVERS
source "$(dirname "$0")/config.sh"

# Export port ranges as environment variables for lib servers
for server in "${ENABLED_SERVERS[@]}"; do
    # Convert server name to uppercase and replace - with _
    server_upper=$(echo "$server" | tr '[:lower:]' '[:upper:]' | tr '-' '_')
    start_var="${server_upper}_START_PORT"
    end_var="${server_upper}_END_PORT"
    
    # Export port ranges to environment variables
    export "$server_upper"_START_PORT="${!start_var}"
    export "$server_upper"_END_PORT="${!end_var}"
    log_with_timestamp "Exported $server_upper port range: ${!start_var}-${!end_var}"
done

# Start all enabled servers
log_with_timestamp "Starting servers..."
"$(dirname "$0")/start-servers.sh"

# Wait a bit for servers to be ready
sleep 2

# Detect OS type to use correct binary
OSTYPE=$(uname -s)
if [[ "$OSTYPE" == "Linux" ]]; then
    BENCH_BIN="bin/bench-tcp.linux"
elif [[ "$OSTYPE" == "Darwin" ]]; then
    BENCH_BIN="bin/bench-tcp.mac"
else
    log_with_timestamp "Unsupported OS: $OSTYPE"
    exit 1
fi

# 运行基准测试
run_benchmark() {
    local server=$1
    local duration=$2
    local concurrency=$3
    local start_port=$4
    local end_port=$5
    
    log_with_timestamp "Running benchmark for $server..."
    
    if [ ! -x "$BENCH_BIN" ]; then
        log_with_timestamp "Error: Benchmark executable not found or not executable: $BENCH_BIN"
        return 1
    fi

    log_with_timestamp "Running bench-tcp with parameters: -d $duration -c $concurrency --addr 127.0.0.1:$start_port-$end_port --open-tmp-result $EXTRA_ARGS"
    # 运行bench-tcp客户端并将输出重定向到对应的.tps文件（使用PID后缀）
    "$BENCH_BIN" -d "$duration" -c "$concurrency" --addr "127.0.0.1:$start_port-$end_port" --open-tmp-result $EXTRA_ARGS | tee "$SCRIPT_DIR/output/$server.$$.tps"
}

log_with_timestamp "Running benchmarks with concurrency: $CONCURRENCY, duration: $DURATION"

# 创建报表输出目录和文件
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
OUTPUT_FILE="$SCRIPT_DIR/benchmark_results.md"
mkdir -p "$SCRIPT_DIR/output"

# 创建报表文件头部
echo "# Benchmark Results" > "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"
echo "Test started at: $(date)" >> "$OUTPUT_FILE"
echo "Operating System: $(uname -s)" >> "$OUTPUT_FILE"
echo "Concurrency: $CONCURRENCY" >> "$OUTPUT_FILE"
echo "Duration: $DURATION" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"
echo "| Framework Name | TPS(Start) | TPS(Middle) | TPS(End) | CPU(Start) | CPU(Middle) | CPU(End) | Memory(Start) | Memory(Middle) | Memory(End) |" >> "$OUTPUT_FILE"
echo "|----------------|------------|-------------|----------|------------|-------------|----------|---------------|---------------|-------------|" >> "$OUTPUT_FILE"

# Run benchmarks for each enabled server
log_with_timestamp "Running benchmarks..."
for server in "${ENABLED_SERVERS[@]}"; do
    log_with_timestamp "Benchmarking $server..."
    
    # Convert server name to uppercase and replace - with _
    server_upper=$(echo "$server" | tr '[:lower:]' '[:upper:]' | tr '-' '_')
    start_var="${server_upper}_START_PORT"
    end_var="${server_upper}_END_PORT"
    
    # Run bench-tcp with port range
    log_with_timestamp "Testing $server on port range ${!start_var}-${!end_var}..."
    
    if [ ! -x "$BENCH_BIN" ]; then
        log_with_timestamp "Error: Benchmark executable not found or not executable: $BENCH_BIN"
        exit 1
    fi
    
    # 启动性能数据收集（在后台）
    # 获取服务器PID
    SERVER_PID_FILE="/tmp/bench-tcp-$server.pid"
    if [ -f "$SERVER_PID_FILE" ]; then
        SERVER_PID=$(cat "$SERVER_PID_FILE")
        log_with_timestamp "Found server PID for $server: $SERVER_PID"
        
        # 将 DURATION 转换为秒数
        DURATION_SECONDS=$(duration_to_seconds "$DURATION")
        log_with_timestamp "Duration converted: $DURATION -> $DURATION_SECONDS seconds"
        
        "$SCRIPT_DIR/collect_metrics.sh" "$server" $SERVER_PID $DURATION_SECONDS "$OUTPUT_FILE" $$ &
        COLLECTOR_PID=$!
    else
        log_with_timestamp "Warning: PID file not found for $server: $SERVER_PID_FILE"
        log_with_timestamp "Skipping metrics collection for $server"
    fi
    
    # Run the benchmark with specified parameters and capture output
    log_with_timestamp "Starting benchmark for $server..."
    BENCHMARK_OUTPUT=$(run_benchmark "$server" "$DURATION" "$CONCURRENCY" "${!start_var}" "${!end_var}")
    BENCHMARK_RESULT=$?
    
    if [ $BENCHMARK_RESULT -ne 0 ]; then
        log_with_timestamp "Error: Benchmark failed for $server with exit code $BENCHMARK_RESULT"
    fi
    
    # 解析TPS数据
    log_with_timestamp "Parsing TPS data for $server..."
    echo "$BENCHMARK_OUTPUT" | grep -E "[0-9]+s:[0-9]+/s" | tail -1 > "$SCRIPT_DIR/output/$server.$$.tps.raw"
    
    # 检查是否有TPS数据
    if [ ! -s "$SCRIPT_DIR/output/$server.$$.tps.raw" ]; then
        log_with_timestamp "Warning: No TPS data found for $server. Benchmark output was:"
        log_with_timestamp "$BENCHMARK_OUTPUT"
        # 创建空的TPS文件以避免后续错误
        echo "start N/A" > "$SCRIPT_DIR/output/$server.$$.tps"
        echo "middle N/A" >> "$SCRIPT_DIR/output/$server.$$.tps"
        echo "end N/A" >> "$SCRIPT_DIR/output/$server.$$.tps"
    else
        # 从最后一行TPS统计中提取开始、中间、结束的TPS值
        TPS_LINE=$(cat "$SCRIPT_DIR/output/$server.$$.tps.raw")
        log_with_timestamp "TPS raw data for $server: $TPS_LINE"
        
        # 提取所有TPS值：1s:485714/s 2s:502793/s ... 格式
        ALL_TPS=($(echo "$TPS_LINE" | grep -o '[0-9]\+s:[0-9]\+/s' | sed 's/[0-9]\+s:\([0-9]\+\)\/s/\1/'))
        TPS_COUNT=${#ALL_TPS[@]}
        
        if [ $TPS_COUNT -gt 0 ]; then
            START_TPS=${ALL_TPS[0]}
            MIDDLE_TPS=${ALL_TPS[$((TPS_COUNT/2))]}
            END_TPS=${ALL_TPS[$((TPS_COUNT-1))]}
            log_with_timestamp "Extracted TPS values for $server: Start=$START_TPS, Middle=$MIDDLE_TPS, End=$END_TPS"
        else
            log_with_timestamp "Warning: Failed to parse TPS values from: $TPS_LINE"
            START_TPS="N/A"
            MIDDLE_TPS="N/A" 
            END_TPS="N/A"
        fi
        
        # 保存TPS数据到标准格式文件
        echo "start $START_TPS" > "$SCRIPT_DIR/output/$server.$$.tps"
        echo "middle $MIDDLE_TPS" >> "$SCRIPT_DIR/output/$server.$$.tps"
        echo "end $END_TPS" >> "$SCRIPT_DIR/output/$server.$$.tps"
        log_with_timestamp "TPS data saved to $SCRIPT_DIR/output/$server.$$.tps"
    fi
    
    # 等待性能数据收集完成
    if [ -n "$COLLECTOR_PID" ]; then
        wait $COLLECTOR_PID 2>/dev/null || true
    fi
    
    log_with_timestamp "Completed benchmarking $server"
    
    # 在每个框架运行完后等待4秒
    log_with_timestamp "Waiting 4 seconds before next framework..."
    sleep 4
done

# Stop all servers
log_with_timestamp "Stopping servers..."
"$(dirname "$0")/stop-servers.sh"

log_with_timestamp "Benchmark complete!"

# 生成最终报表
log_with_timestamp "Generating final report..."
"$SCRIPT_DIR/generate_report.sh" $$

log_with_timestamp "Results saved to: $OUTPUT_FILE" 