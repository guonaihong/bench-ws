#!/bin/bash
SEC="30s"
SLEEP_SEC="40"

run_gws_test() {
    echo "gws:"
    killall gws-std.linux &>/dev/null
    ./gws-std.linux --ws-addr ":9001" &
    PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9001/
    kill $PID
    sleep $SLEEP_SEC
    echo ""
    
    echo "gws.asyncwrite:"
    killall gws-std.linux &>/dev/null
    ./gws-std.linux -a --ws-addr ":9001" &
    PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9001/
    kill $PID
    sleep $SLEEP_SEC
    echo ""
    
    echo "gws-std:"
    killall gws-std.linux &>/dev/null
    ./gws-std.linux --ws-addr ":9002" &
    PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9002/
    kill $PID
    sleep $SLEEP_SEC
    echo ""
    
    echo "gws.tcp-delay:"
    killall gws-std.linux &>/dev/null
    ./gws.linux -o  --ws-addr ":9002" &
    PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9002/
    kill $PID
    sleep $SLEEP_SEC
    echo ""
    
    echo "gws-std.tcp-delay:"
    killall gws-std.linux &>/dev/null
    ./gws-std.linux -o  --ws-addr ":9002" &
    PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9002/
    kill $PID
    sleep $SLEEP_SEC
    echo ""
    
    echo "gws-std.asyncwrite:"
    killall gws-std.linux &>/dev/null
    ./gws-std.linux -a --ws-addr ":9002" &
    PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9002/
    kill $PID
    sleep $SLEEP_SEC
    echo ""

}

run_greatws_windows_test() {
    killall greatws.linux &>/dev/null
    ./greatws.linux -w $1 -a ":9000" &
    PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9000/
    kill $PID
    sleep $SLEEP_SEC
    killall greatws.linux &>/dev/null
    echo ""
}

run_quickws_windows_test() {
    killall quickws.linux &>/dev/null
    ./quickws.linux -w $1 -a ":9000" &
    PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9000/
    kill $PID
    sleep $SLEEP_SEC
    killall quickws.linux &>/dev/null
    echo ""
}

run_quickws_windows_delay_and_tcp_delay_test() {
    killall quickws.linux &>/dev/null
    ./quickws.linux -o -w $1 --use-delay-write -a ":9000" --delay-write-init-buffer-size 11264 &
    PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9000/
    kill $PID
    sleep $SLEEP_SEC
    killall quickws.linux &>/dev/null
    echo ""
}

run_quickws_windows_delay_test() {
    killall quickws.linux &>/dev/null
    ./quickws.linux -w $1 --use-delay-write -a ":9000" --delay-write-init-buffer-size 11264 &
    PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9000/
    kill $PID
    sleep $SLEEP_SEC
    killall quickws.linux &>/dev/null
    echo ""
}

run_quickws_windows_tcp_delay_test() {
    killall quickws.linux &>/dev/null
    ./quickws.linux -o -w $1 -a ":9000" &
    PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9000/
    kill $PID
    sleep $SLEEP_SEC
    killall quickws.linux &>/dev/null
    echo ""
}

run_quickws_bufio_delay_test() {
    killall quickws.linux &>/dev/null
    ./quickws.linux -u -b $1 --use-delay-write -a ":9000" --delay-write-init-buffer-size 11264 &
    PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9000/
    kill $PID
    sleep $SLEEP_SEC
    killall quickws.linux &>/dev/null
    echo ""
}

run_quickws_bufio_tcp_delay_test() {
    killall quickws.linux &>/dev/null
    ./quickws.linux -u -o -b $1 -a ":9000" &
    PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9000/
    kill $PID
    sleep $SLEEP_SEC
    killall quickws.linux &>/dev/null
    echo ""
}

run_gorilla_test() {
    echo "gorilla-linux-ReadMessage:"
    killall gorilla.linux  &>/dev/null
    ./gorilla.linux --ws-addr ":9003" &
    PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9003/
    kill $PID
    sleep $SLEEP_SEC
    echo ""
    
    echo "gorilla-linux-UseReader:"
    killall gorilla.linux &>/dev/null
    ./gorilla.linux --ws-addr ":9004" -u &
    PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9004/
    kill $PID
    sleep $SLEEP_SEC
    echo ""

    echo "gorilla-linux-ReadMessage.tcp.delay:"
    killall gorilla.linux  &>/dev/null
    ./gorilla.linux -o --ws-addr ":9003" &
    PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9003/
    kill $PID
    sleep $SLEEP_SEC
    echo ""

    echo "gorilla-linux-UseReader.tcp.delay:"
    killall gorilla.linux &>/dev/null
    ./gorilla.linux -o --ws-addr ":9004" -u &
    PID=$!
    sleep 1
    tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9004/
    kill $PID
    sleep $SLEEP_SEC
    echo ""
}

run_quickws_all_test() {
  echo "quickws.windows.delay.and.tcp.delay.test.32x:"
  run_quickws_windows_delay_and_tcp_delay_test 32
  echo "quickws.windows.delay.and.tcp.delay.test.24x:"
  run_quickws_windows_delay_and_tcp_delay_test 24
  echo "quickws.windows.delay.and.tcp.delay.test.16x:"
  run_quickws_windows_delay_and_tcp_delay_test 16
  echo "quickws.windows.delay.and.tcp.delay.test.8x:"
  run_quickws_windows_delay_and_tcp_delay_test 8
  echo "quickws.windows.delay.and.tcp.delay.test.4x:"
  run_quickws_windows_delay_and_tcp_delay_test 4
  echo "quickws.windows.delay.and.tcp.delay.test.1x:"
  run_quickws_windows_delay_and_tcp_delay_test 1

  echo "quickws.windows.delay.32x:"
  run_quickws_windows_delay_test 32
  echo "quickws.windows.delay.24x:"
  run_quickws_windows_delay_test 24
  echo "quickws.windows.delay.16x:"
  run_quickws_windows_delay_test 16
  echo "quickws.windows.delay.8x:"
  run_quickws_windows_delay_test 8
  echo "quickws.windows.delay.4x:"
  run_quickws_windows_delay_test 4
  echo "quickws.windows.delay.1x:"
  run_quickws_windows_delay_test 1

  echo "quickws.windows.tcp.delay.32x:"
  run_quickws_windows_tcp_delay_test 32
  echo "quickws.windows.tcp.delay.24x:"
  run_quickws_windows_tcp_delay_test 24
  echo "quickws.windows.tcp.delay.16x:"
  run_quickws_windows_tcp_delay_test 16
  echo "quickws.windows.tcp.delay.8x:"
  run_quickws_windows_tcp_delay_test 8
  echo "quickws.windows.tcp.delay.4x:"
  run_quickws_windows_tcp_delay_test 4
  echo "quickws.windows.tcp.delay.1x:"
  run_quickws_windows_tcp_delay_test 1


  echo "quickws.bufio.delay.32x:"
  run_quickws_bufio_delay_test 32
  echo "quickws.bufio.delay.24:"
  run_quickws_bufio_delay_test 24
  echo "quickws.bufio.delay.16x:"
  run_quickws_bufio_delay_test 16
  echo "quickws.bufio.delay.8x:"
  run_quickws_bufio_delay_test 8
  echo "quickws.bufio.delay.4x:"
  run_quickws_bufio_delay_test 4
  echo "quickws.bufio.delay.1x:"
  run_quickws_bufio_delay_test 1

  echo "quickws.bufio.tcpdelay.32x:"
  run_quickws_bufio_tcp_delay_test 32
  echo "quickws.bufio.tcpdelay.24:"
  run_quickws_bufio_tcp_delay_test 24
  echo "quickws.bufio.tcpdelay.16:"
  run_quickws_bufio_tcp_delay_test 16
  echo "quickws.bufio.tcpdelay.8x:"
  run_quickws_bufio_tcp_delay_test 8
  echo "quickws.bufio.tcpdelay.4x:"
  run_quickws_bufio_tcp_delay_test 4
  echo "quickws.bufio.tcpdelay.1x:"
  run_quickws_bufio_tcp_delay_test 1
}

run_quickws_mini_test() {
  echo "quickws.windows.delay.1x:"
  run_quickws_windows_delay_test 1
}

run_greatws_test() {
  run_greatws_windows_test 1
}

run_nettyws_test() {
  echo "nettyws:"
  killall nettyws.linux &>/dev/null
  ./nettyws.linux --ws-addr ":9005" &
  PID=$!
  sleep 1
  tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9005/
  kill $PID
  sleep $SLEEP_SEC
  echo ""
}

run_gobwas_test() {
  echo "gobwas:"
  killall gobwas &>/dev/null
  ./gobwas.linux --ws-addr ":9006" &
  PID=$!
  sleep 1
  tcpkali -c 10000 --connect-rate 10000 -r 10000 -T 30s -f 1K.txt --ws 127.0.0.1:9006/
  kill $PID
  sleep $SLEEP_SEC
}


#run_gorilla_test

#run_gws_test

run_quickws_mini_test

run_greatws_test


