# Benchmark Results

Test completed at: 2025年 06月 07日 星期六 17:57:48 CST
Operating System: Linux

| 框架名 | TPS(开始) | TPS(中间) | TPS(结束) | CPU(最大)% | CPU(最小)% | CPU(平均)% | 内存(最大)MB | 内存(最小)MB | 内存(平均)MB | 线程(最大) | 线程(最小) | 线程(平均) | FD(最大) | FD(最小) | FD(平均) |
|--------|-----------|-----------|-----------|------------|------------|------------|-------------|-------------|-------------|------------|------------|------------|---------|---------|---------|
| fasthttp-ws-std | 680091 | 780138 | 787481 | 1660% | 170.0% | 1256.2% | 63MB | 26MB | 49MB | 48 | 41 | 44 | 1006 | 703 | 930 |
| gobwas | N/A | N/A | N/A | N/A% | N/A% | N/A% | N/AMB | N/AMB | N/AMB | N/A | N/A | N/A | N/A | N/A | N/A |
| gorilla | 624808 | 741727 | 749849 | 1660% | 60.0% | 1237.5% | 69MB | 17MB | 47MB | 48 | 44 | 46 | 1006 | 439 | 864 |
| greatws | 736475 | 835859 | 870725 | 1380% | 130.0% | 1045.0% | 40MB | 26MB | 36MB | 107 | 34 | 88 | 1007 | 257 | 819 |
| gws | 0 | 0 | 0 | 140.0% | 0.0% | 35.0% | 25MB | 25MB | 25MB | 34 | 34 | 34 | 188 | 188 | 188 |
| gws-std | 0 | 0 | 0 | 20.0% | 0.0% | 5.0% | 19MB | 16MB | 18MB | 34 | 32 | 33 | 354 | 36 | 115 |
| hertz | 526934 | 620649 | 627336 | 1050% | 260.0% | 840.0% | 100MB | 52MB | 87MB | 38 | 31 | 36 | 1011 | 627 | 915 |
| hertz-std | N/A | N/A | N/A | N/A% | N/A% | N/A% | N/AMB | N/AMB | N/AMB | N/A | N/A | N/A | N/A | N/A | N/A |
| nbio-blocking | N/A | N/A | N/A | N/A% | N/A% | N/A% | N/AMB | N/AMB | N/AMB | N/A | N/A | N/A | N/A | N/A | N/A |
| nbio-mixed | N/A | N/A | N/A | N/A% | N/A% | N/A% | N/AMB | N/AMB | N/AMB | N/A | N/A | N/A | N/A | N/A | N/A |
| nbio-nonblocking | N/A | N/A | N/A | N/A% | N/A% | N/A% | N/AMB | N/AMB | N/AMB | N/A | N/A | N/A | N/A | N/A | N/A |
| nbio-std | N/A | N/A | N/A | N/A% | N/A% | N/A% | N/AMB | N/AMB | N/AMB | N/A | N/A | N/A | N/A | N/A | N/A |
| nettyws | 1087455 | 1287342 | 1317332 | 1650% | 40.0% | 1186.7% | 30MB | 16MB | 26MB | 53 | 33 | 48 | 1006 | 256 | 818 |
| nhooyr | N/A | N/A | N/A | N/A% | N/A% | N/A% | N/AMB | N/AMB | N/AMB | N/A | N/A | N/A | N/A | N/A | N/A |
| quickws | 1097063 | 1252112 | 1306932 | 1460% | 470.0% | 1173.2% | 26MB | 21MB | 24MB | 55 | 42 | 51 | 1006 | 1006 | 1006 |

