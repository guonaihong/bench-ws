#!/bin/bash

# Source the config file
source "$(dirname "$0")/config.sh"

# Function to stop a server
stop_server() {
    local lib_name=$1
    local pid_file="/tmp/bench-tcp-$lib_name.pid"
    
    if [ -f "$pid_file" ]; then
        local pid
        pid=$(cat "$pid_file")
        if ps -p "$pid" > /dev/null; then
            echo "Stopping $lib_name server (PID: $pid)"
            kill "$pid"
            rm "$pid_file"
        else
            echo "$lib_name server not running (PID: $pid)"
            rm "$pid_file"
        fi
    else
        echo "$lib_name server not running (no PID file)"
    fi
}

# Stop enabled servers
for server in "${ENABLED_SERVERS[@]}"; do
    stop_server "$server"
done

echo "All enabled servers stopped" 