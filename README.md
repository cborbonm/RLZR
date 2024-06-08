# CS244 RLZR
## Overview
This project focuses on the re-implementation of the [LZR](https://github.com/stanford-esrg/lzr) module - see website for information on LZR and usage.

The repository is structured as follows:
- `lzr`: LZR package version as of the time of its publication (August 2021)
- `rlzr`: RLZR package version.

## Contribution
We aim to build `RLZR`, centered around the design of `LZR`. We are specifically focusing on implementing harnessing and end-to-end testing frameworks that can provide clearer metrics for the performance of `RLZR`.

## Requirements
- Golang
- ZMap
- A public IP Address and associated gateway MAC address ([reference](#identifying-machines-ip-address-and-mac-gateway-address))

## Quick Start
> NOTE: This repository can only be run in Linux OS.

Start by cloning the repository to your directory of preference:
```
git clone https://github.com/cborbonm/RLZR.git cs244-rlzr
```
Then navigate to `cs244-rlzr`. Set up the repository dependecies via:
```
go mod tidy
```

Next, navigate to the `rlzr` in `cs244-rlzr` repository, and set up via Makefile.

Run make to retrieve `rlzr` executable via `make`

```
make all source-ip=YOURIPADDR
```

Follow [LZR](https://github.com/stanford-esrg/lzr)'s directions on running LZR. Example command to scan HTTP handshakes on port 80:

```
sudo zmap --target-port=80 --bandwidth=1G --output-filter="success = 1 && repeat = 0" -f "saddr,daddr,sport,dport,seqnum,acknum,window" -O json --source-ip=YOURIPADDR --gateway-mac=YOURGATEWAYMACADDR -i enp0s9 --blacklist-file=/etc/zmap/blocklist.conf --max-targets 1% | sudo ./rlzr --handshakes http --sendInterface enp0s9
```
>NOTE: Make sure to replace each instance of `YOURIPADDR` and `YOURGATEMACADDR` values with your actual associated values.

Once finished, run `make clean` to clean executables.

## Identifying Machine's IP Address and Mac Gateway Address
Provided you have a Linux machine with a valid IP address, you can retrieve its relevant information via the following:
- `ip addr show | grep inet`: Identify your machine's valid public IP address.
- `ip route show`: Identify the line associated with your machine's valid public IP address. Its MAC gateway address is provided as well.


## Contributors

- Wilmer Zuna
- Ola Adekola
- Carolina Borbon
