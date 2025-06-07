#!/bin/bash

# 报表生成脚本
# Usage: generate_report.sh [pid]
# 如果提供PID参数，则只处理该PID的文件；否则处理最新的文件

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
OUTPUT_FILE="$SCRIPT_DIR/benchmark_results.md"
TARGET_PID=$1  # 可选的PID参数

# 创建报表文件头部
echo "# Benchmark Results" > "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"
echo "Test completed at: $(date)" >> "$OUTPUT_FILE"
echo "Operating System: $(uname -s)" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"
echo "| 框架名 | TPS(开始) | TPS(中间) | TPS(结束) | CPU(最大)% | CPU(最小)% | CPU(平均)% | 内存(最大)MB | 内存(最小)MB | 内存(平均)MB | 线程(最大) | 线程(最小) | 线程(平均) | FD(最大) | FD(最小) | FD(平均) |" >> "$OUTPUT_FILE"
echo "|--------|-----------|-----------|-----------|------------|------------|------------|-------------|-------------|-------------|------------|------------|------------|---------|---------|---------|" >> "$OUTPUT_FILE"

# 在屏幕上显示生成的报表
echo ""
echo "=========================================="
echo "Generated Report Content:"
echo "=========================================="

# 打印美化的表格到屏幕
echo "# Benchmark Results"
echo ""
echo "Test completed at: $(date)"
echo "Operating System: $(uname -s)"
echo ""

# 打印表格头部 - 使用固定宽度对齐
printf "%-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s\n" \
    "Framework" "TPS-Start" "TPS-Mid" "TPS-End" "CPU-Max" "CPU-Min" "CPU-Avg" "Mem-Max" "Mem-Min" "Mem-Avg" "Thr-Max" "Thr-Min" "Thr-Avg" "FD-Max" "FD-Min" "FD-Avg"

printf "%-8s-+-%-8s-+-%-8s-+-%-8s-+-%-8s-+-%-8s-+-%-8s-+-%-8s-+-%-8s-+-%-8s-+-%-8s-+-%-8s-+-%-8s-+-%-8s-+-%-8s-+-%-8s\n" \
    "--------" "--------" "--------" "--------" "--------" "--------" "--------" "--------" "--------" "--------" "--------" "--------" "--------" "--------" "--------" "--------"

# 处理每个框架的数据并格式化输出
if [ -n "$TARGET_PID" ]; then
    # 如果指定了PID，处理该PID的文件
    tps_pattern="$SCRIPT_DIR/output/*.$TARGET_PID.tps"
else
    # 如果没有指定PID，处理最新的文件（按修改时间排序）
    tps_pattern="$SCRIPT_DIR/output/*.tps"
fi

# 获取所有匹配的TPS文件
for tps_file in $tps_pattern; do
    if [ -f "$tps_file" ]; then
        # 从文件名提取框架名称和PID
        filename=$(basename "$tps_file")
        if [[ "$filename" =~ ^(.+)\.([0-9]+)\.tps$ ]]; then
            framework="${BASH_REMATCH[1]}"
            file_pid="${BASH_REMATCH[2]}"
        else
            # 兼容旧格式（没有PID的文件名）
            framework=$(basename "$tps_file" .tps)
            file_pid=""
        fi
        
        # 如果没有指定TARGET_PID，需要确保我们处理的是每个框架的最新文件
        if [ -z "$TARGET_PID" ]; then
            # 查找该框架的最新文件
            latest_file=$(ls -t "$SCRIPT_DIR/output/$framework".*.tps 2>/dev/null | head -1)
            if [ "$tps_file" != "$latest_file" ]; then
                continue  # 跳过不是最新的文件
            fi
        fi
        
        # 读取TPS数据
        if [ -f "$tps_file" ]; then
            start_tps=$(grep "^start" "$tps_file" | awk '{print $2}' || echo "N/A")
            middle_tps=$(grep "^middle" "$tps_file" | awk '{print $2}' || echo "N/A")
            end_tps=$(grep "^end" "$tps_file" | awk '{print $2}' || echo "N/A")
        else
            start_tps="N/A"
            middle_tps="N/A"
            end_tps="N/A"
        fi
        
        # 读取CPU数据
        if [ -n "$file_pid" ]; then
            cpu_file="$SCRIPT_DIR/output/$framework.$file_pid.cpu"
        else
            cpu_file="$SCRIPT_DIR/output/$framework.cpu"
        fi
        if [ -f "$cpu_file" ]; then
            max_cpu=$(grep "^max" "$cpu_file" | awk '{print $2}' || echo "N/A")
            min_cpu=$(grep "^min" "$cpu_file" | awk '{print $2}' || echo "N/A")
            avg_cpu=$(grep "^avg" "$cpu_file" | awk '{print $2}' || echo "N/A")
        else
            max_cpu="N/A"
            min_cpu="N/A"
            avg_cpu="N/A"
        fi
        
        # 读取内存数据
        if [ -n "$file_pid" ]; then
            mem_file="$SCRIPT_DIR/output/$framework.$file_pid.mem"
        else
            mem_file="$SCRIPT_DIR/output/$framework.mem"
        fi
        if [ -f "$mem_file" ]; then
            max_mem=$(grep "^max" "$mem_file" | awk '{print $2}' || echo "N/A")
            min_mem=$(grep "^min" "$mem_file" | awk '{print $2}' || echo "N/A")
            avg_mem=$(grep "^avg" "$mem_file" | awk '{print $2}' || echo "N/A")
        else
            max_mem="N/A"
            min_mem="N/A"
            avg_mem="N/A"
        fi
        
        # 读取线程数据
        if [ -n "$file_pid" ]; then
            thread_file="$SCRIPT_DIR/output/$framework.$file_pid.threads"
        else
            thread_file="$SCRIPT_DIR/output/$framework.threads"
        fi
        if [ -f "$thread_file" ]; then
            max_threads=$(grep "^max" "$thread_file" | awk '{print $2}' || echo "N/A")
            min_threads=$(grep "^min" "$thread_file" | awk '{print $2}' || echo "N/A")
            avg_threads=$(grep "^avg" "$thread_file" | awk '{print $2}' || echo "N/A")
        else
            max_threads="N/A"
            min_threads="N/A"
            avg_threads="N/A"
        fi
        
        # Read FD data
        if [ -n "$file_pid" ]; then
            fd_file="$SCRIPT_DIR/output/$framework.$file_pid.fd"
        else
            fd_file="$SCRIPT_DIR/output/$framework.fd"
        fi
        if [ -f "$fd_file" ]; then
            max_fd=$(grep "^max" "$fd_file" | awk '{print $2}' || echo "N/A")
            min_fd=$(grep "^min" "$fd_file" | awk '{print $2}' || echo "N/A")
            avg_fd=$(grep "^avg" "$fd_file" | awk '{print $2}' || echo "N/A")
        else
            max_fd="N/A"
            min_fd="N/A"
            avg_fd="N/A"
        fi
        
        # 格式化打印到屏幕
        printf "%-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s | %-8s\n" \
            "$framework" "$start_tps" "$middle_tps" "$end_tps" \
            "${max_cpu}%" "${min_cpu}%" "${avg_cpu}%" \
            "${max_mem}MB" "${min_mem}MB" "${avg_mem}MB" \
            "$max_threads" "$min_threads" "$avg_threads" \
            "$max_fd" "$min_fd" "$avg_fd"
        
        # 写入到文件（保持原格式）
        echo "| $framework | $start_tps | $middle_tps | $end_tps | $max_cpu% | $min_cpu% | $avg_cpu% | ${max_mem}MB | ${min_mem}MB | ${avg_mem}MB | $max_threads | $min_threads | $avg_threads | $max_fd | $min_fd | $avg_fd |" >> "$OUTPUT_FILE"
    fi
done

echo ""
echo "=========================================="

echo "" >> "$OUTPUT_FILE"
echo "Report generated successfully at: $OUTPUT_FILE" 