# Change Log

## v1.1.0
2022-0425

### Changed
1. Use fasthttp, faster and less resource.
2. Add Options -F  Faster request ([0]/1). Select 1 faster request, but time detail cannot be counted
3. If use -F 1, Faster, increase by about 50% (be based on: ysab -n 500 -r 1500 ...)

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