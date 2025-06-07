#!/bin/bash

# Run benchmark with 100,000 concurrent connections
# First argument controls whether to rebuild (default: true)
REBUILD=${1:-true}

"$(dirname "$0")/benchmark-core.sh" 100000 10s "$REBUILD"

# Exit with the same status as the core script
exit $? 