#!/bin/bash

if [ ! -d ./node_modules ];then
    npm install
fi

npm run serve &
./bin/web.linux &
./script/tps-all-benchmark.sh
