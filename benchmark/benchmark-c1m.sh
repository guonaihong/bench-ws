#!/bin/bash

# Run benchmark with 1,000,000 concurrent connections
# First argument controls whether to rebuild (default: true)
REBUILD=${1:-true}

"$(dirname "$0")/benchmark-core.sh" 1000000 10s "$REBUILD"

# Exit with the same status as the core script
exit $? 