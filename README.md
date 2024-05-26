# CS244 RLZR
## Overview
This project focuses on the re-implementation of the [LZR](https://github.com/stanford-esrg/lzr) module - see website for information on LZR and usage.

## Contribution
We aim to build `RLZR`, centered around the design of `LZR`. We are specifically focusing on implementing harnessing and end-to-end testing frameworks that can provide clearer metrics for the performance of `RLZR`.

## Quick Start
> NOTE: This repository can only be run in Linux OS. Alternative ways of testing this repository involve Virtual Machines and Docker.

Start by cloning the repository to your directory of preference:
```
git clone https://github.com/cborbonm/RLZR.git cs244-rlzr
```
Navigate to the `rlzr` in `cs244-rlzr` repository, and set up via Makefile.

1. Run make to retrieve `rlzr` executable via `make`
> NOTE: You can only run the `rlzr` executable if in Linux OS. Otherwise, you can only run go tests.

```
make all source-ip=256.256.256.256/32
```

2. Follow [LZR](https://github.com/stanford-esrg/lzr)'s directions on running LZR. Example using random port (9002):

```
sudo zmap --target-port=9002 --output-filter="success = 1 && repeat = 0" \
-f "saddr,daddr,sport,dport,seqnum,acknum,window" -O json --source-ip=$source-ip | \
sudo ./lzr --handshakes http,tls
```

3. Once finished, run `make clean` to clean executables.

## Actually Running on Ola server
sudo zmap --target-port=80 --output-filter="success = 1 && repeat = 0" \
-f "saddr,daddr,sport,dport,seqnum,acknum,window" -O json -n 40000 --source-ip=10.129.44.6 --gateway-mac=74:56:3c:fb:86:a5 | \
sudo ./lzr --handshakes http,tls --sendInterface enp11s0

## Contributors

- Wilmer Zuna
- Ola Adekola
- Carolina Borbon
