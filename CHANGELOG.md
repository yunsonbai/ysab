# Change Log

## v1.1.0
2023-1118

### Changed
1. 当使用urlsfile参数时，r参数也生效支持多轮测试 (r takes effect when using urlsfile, supporting multiple rounds of testing)
2. 新增T参数，支持连续测试T秒。T大于0时，r失效 (Added T to support continuous testing for T seconds. r is invalid when T is greater than 0)
3. 优化一些代码 (Optimize some code)
4. 添加-b参数, 读取body体的buf大小, 默认256K (-b reader buf size, default 256K)
5. 基于go1.22.12编写 (Based on go1.22.12)

## v1.0.2
2023-0424

### Changed
1. Correction result description
2. Optimization Statistics.
3. Fix Put/Delete bug.
4. Add transfer metrics: TransferRate/sec (Byte).

Compared with 1.0.1, the performance is improved by about 10%.


## v1.0.1
2021-0804

### Changed
1. Allow setting Host in header

## v1.0.0
2021-0528

### Changed
1. More stable performance
2. Make the number of tcp connections more stable
3. Faster, increase by about 10%