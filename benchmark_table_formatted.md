# WebSocket Framework Performance Benchmark

## Option 1: Improved ASCII Table
```
Framework     │ TPS Start │ TPS Mid   │ TPS End   │ CPU Max │ CPU Min │ CPU Avg  │ Mem Max │ Mem Min │ Mem Avg │ Thr Max │ Thr Min │ Thr Avg │ FD Max  │ FD Min │ FD Avg
──────────────┼───────────┼───────────┼───────────┼─────────┼─────────┼──────────┼─────────┼─────────┼─────────┼─────────┼─────────┼─────────┼─────────┼────────┼───────
greatws       │   658,890 │   964,027 │   965,463 │  1730%  │   0.0%  │  1172.2% │  114 MB │  38 MB  │ 104 MB  │     200 │      79 │     197 │  10,039 │     38 │  7,404
greatws-event │   725,910 │ 1,085,014 │ 1,087,399 │  1560%  │   0.0%  │  1052.7% │   86 MB │  36 MB  │  79 MB  │     230 │      62 │     227 │  10,038 │     38 │  7,404
quickws       │   673,534 │   957,436 │   957,599 │  1730%  │   0.0%  │  1101.7% │  130 MB │  28 MB  │ 125 MB  │      68 │      44 │      64 │  10,006 │      6 │  7,205
```

## Option 2: Markdown Table
| Framework     | TPS Start | TPS Mid   | TPS End   | CPU Max | CPU Min | CPU Avg  | Mem Max | Mem Min | Mem Avg | Thr Max | Thr Min | Thr Avg | FD Max  | FD Min | FD Avg |
|---------------|----------:|----------:|----------:|--------:|--------:|---------:|--------:|--------:|--------:|--------:|--------:|--------:|--------:|-------:|-------:|
| greatws       |   658,890 |   964,027 |   965,463 |  1730%  |   0.0%  |  1172.2% |  114 MB |  38 MB  | 104 MB  |     200 |      79 |     197 |  10,039 |     38 |  7,404 |
| greatws-event |   725,910 | 1,085,014 | 1,087,399 |  1560%  |   0.0%  |  1052.7% |   86 MB |  36 MB  |  79 MB  |     230 |      62 |     227 |  10,038 |     38 |  7,404 |
| quickws       |   673,534 |   957,436 |   957,599 |  1730%  |   0.0%  |  1101.7% |  130 MB |  28 MB  | 125 MB  |      68 |      44 |      64 |  10,006 |      6 |  7,205 |

## Option 3: Grouped by Metric Categories
### Throughput (TPS)
| Framework     |     Start |       Mid |       End |
|---------------|----------:|----------:|----------:|
| greatws       |   658,890 |   964,027 |   965,463 |
| greatws-event |   725,910 | 1,085,014 | 1,087,399 |
| quickws       |   673,534 |   957,436 |   957,599 |

### CPU Usage
| Framework     |   Max |  Min |     Avg |
|---------------|------:|-----:|--------:|
| greatws       | 1730% | 0.0% | 1172.2% |
| greatws-event | 1560% | 0.0% | 1052.7% |
| quickws       | 1730% | 0.0% | 1101.7% |

### Memory Usage
| Framework     |   Max |  Min |   Avg |
|---------------|------:|-----:|------:|
| greatws       | 114MB | 38MB | 104MB |
| greatws-event |  86MB | 36MB |  79MB |
| quickws       | 130MB | 28MB | 125MB |

### Thread Count
| Framework     | Max | Min | Avg |
|---------------|----:|----:|----:|
| greatws       | 200 |  79 | 197 |
| greatws-event | 230 |  62 | 227 |
| quickws       |  68 |  44 |  64 |

### File Descriptors
| Framework     |   Max | Min |   Avg |
|---------------|------:|----:|------:|
| greatws       | 10039 |  38 | 7,404 |
| greatws-event | 10038 |  38 | 7,404 |
| quickws       | 10006 |   6 | 7,205 |

## Option 4: Summary Performance Ranking
| Metric          | 🥇 Best        | 🥈 Second      | 🥉 Third    |
|-----------------|----------------|----------------|-------------|
| Peak TPS        | greatws-event  | quickws        | greatws     |
| CPU Efficiency  | greatws-event  | quickws        | greatws     |
| Memory Usage    | greatws-event  | greatws        | quickws     |
| Thread Count    | quickws        | greatws        | greatws-event | 