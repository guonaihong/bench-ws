#!/bin/bash

# Source the config file
source "$(dirname "$0")/config.sh"

# Detect OS type to use correct binary extension
OSTYPE=$(uname -s)
if [[ "$OSTYPE" == "Linux" ]]; then
    BIN_EXT="linux"
elif [[ "$OSTYPE" == "Darwin" ]]; then
    BIN_EXT="mac"
else
    echo "Unsupported OS: $OSTYPE"
    exit 1
fi

# Function to get port range for a server
get_port_range() {
    local lib_name=$1
    local lib_upper=$(echo "$lib_name" | tr '[:lower:]' '[:upper:]')  # Convert to uppercase using tr
    local start_var="${lib_upper}_START_PORT"
    local end_var="${lib_upper}_END_PORT"
    start_var=${start_var//-/_}  # Replace - with _ for variable names
    end_var=${end_var//-/_}
    
    # Use indirect reference to get the values
    echo "${!start_var} ${!end_var}"
}

# Function to start a server
start_server() {
    local lib_name=$1
    local port_range
    port_range=$(get_port_range "$lib_name")
    read -r start_port end_port <<< "$port_range"
    
    echo "Starting $lib_name server on port range $start_port-$end_port"
    
    # Export port ranges as environment variables
    local lib_upper=$(echo "$lib_name" | tr '[:lower:]' '[:upper:]')  # Convert to uppercase using tr
    lib_upper="${lib_upper//-/_}"
    export "${lib_upper}_START_PORT"="$start_port"
    export "${lib_upper}_END_PORT"="$end_port"
    
    # Determine binary path
    local bin_path="$(dirname "$0")/../bin/${lib_name}.${BIN_EXT}"
    
    if [ ! -x "$bin_path" ]; then
        echo "Error: Binary not found or not executable: $bin_path"
        return 1
    fi
    
    # Start the server in the background using pre-compiled binary
    "$bin_path" &
    
    # Store the PID
    echo $! > "/tmp/bench-tcp-$lib_name.pid"
}

# Start enabled servers
for server in "${ENABLED_SERVERS[@]}"; do
    start_server "$server"
done

echo "All enabled servers started. PIDs stored in /tmp/bench-tcp-*.pid"
echo "Use stop-servers.sh to stop all servers" 