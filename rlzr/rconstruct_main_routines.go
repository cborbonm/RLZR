package rlzr

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

/*
Variables needed for main routines
*/
var (
	source_mac   string
	dest_mac     string
	wq_buf       int
	in_buf       int
	snapshot_len int32
	promiscuous  bool
	handle       *pcap.Handle
	err          error
)

/*
InitParams

Initialize `source_mac` and `dest_mac`

NOTE: Depends on `construct_responses.go`
I consulted the lzr documentation for the parameters.
*/
func InitParams() {
	source_mac = getSourceMacAddr()
	dest_mac = getHostMacAddr()
	wq_buf = 10000
	in_buf = 1000000
	snapshot_len = 1024
	promiscuous = false
}

/*
ConstructWritingQueue

Constructs writing queue channel
Takes `workers` as input to construct writing queue.
*/
func ConstructWritingQueue(workers int) chan packet_metadata {
	write_q := make(chan packet_metadata, wq_buf)
	return write_q
}

/*
ZMapRead

Reads STDIN data into ZMap. Edge cases handled.
*/
func ZMapRead(incoming chan *packet_metadata) {
	// Create STDIN Reader
	std_read := bufio.NewReader(os.Stdin)
	for {
		// Read from ZMap
		input, err := std_read.ReadString(byte('\n'))

		// Edge case: Finished reading
		if err != nil && err == io.EOF {
			fmt.Fprintln(os.Stderr, "Finished Reading Input")
			close(incoming)
			return
		}

		// Retrieve packet, either via ZMap or Input List
		var packet *packet_metadata
		if ReadZMap() {
			packet = convertFromZMapToPacket(input)
		} else {
			packet = convertFromInputListToPacket(input)
		}

		// Edge case: Packet is not loaded
		if packet == nil {
			continue
		}

		// Load packet into incoming channel routine
		incoming <- packet
	}
}

/*
ConstructIncomingRoutine

Creates channel routine to read in ZMap Data.
*/
func ConstructIncomingRoutine(workers int) chan *packet_metadata {
	// Create channel routine to read from ZMap
	zmap_routine := make(chan *packet_metadata, in_buf)
	go ZMapRead(zmap_routine)
	return zmap_routine
}

/*
ProcessWorkerChannelData

Listen in for worker channel data - i.e. queue packets.
Save host mac address and load packet into routine channel.
*/
func ProcessWorkerChannelData(routine chan *packet_metadata, queue chan *gopacket.Packet, worker_id int) {
	for {
		select {
		case data := <-queue:
			packet := convertToPacketM(data)
			if packet == nil {
				continue
			}
			if dest_mac == "" {
				saveHostMacAddr(packet)
			}
			routine <- packet
		}
	}
}

/*
LoadHandlePacket

Generates new packet sources to monitor packets.
Detected packets passed onto queue channel for processing.
Post-cleanup handled.
*/
func LoadHandlePacket(queue chan *gopacket.Packet) {
	// Run task at function exit
	defer handle.Close()

	// Generate new packet source based on link type
	packet_source := gopacket.NewPacketSource(handle, handle.LinkType())

	// Continously scan for packets and capture
	for {
		pcap_packet, _ := packet_source.NextPacket()
		queue <- &pcap_packet
	}
}

/*
ConstructPcapRoutine

Creates routine channel to listen for ZMap SYN Packets.
Handles queue and network cleanup logistics.
*/
func ConstructPcapRoutine(workers int) chan *packet_metadata {
	// Create channel and queue routines
	pcap_routine := make(chan *packet_metadata, in_buf)
	pcap_queue := make(chan *gopacket.Packet, in_buf)

	// Open device
	handle, err = pcap.OpenLive(getDevice(), snapshot_len, promiscuous, pcap.BlockForever)
	// Shut down if error
	if err != nil {
		panic(err)
	}
	// Filter out ZMap SYN packets
	filter_err := handle.SetBPFFilter("tcp[tcpflags] != tcp-syn") // NOTE: Referenced source code for filter
	// Shut down if error while filtering
	if filter_err != nil {
		panic(filter_err)
	}

	// Spawn channel workers to listen in for incoming data for packet processing
	for i := 0; i < workers; i++ {
		go ProcessWorkerChannelData(pcap_routine, pcap_queue, i)
	}

	// Handle packet monitoring
	go LoadHandlePacket(pcap_queue)

	return pcap_routine
}

/*
timeoutAlg

Sleep until timeout, and process packet data.

NOTE: pState is dependent on `concurrentMap.go`
*/
func timeoutAlg(ipMeta *rCMap, queue chan *packet_metadata, timeoutIncoming chan *packet_metadata,
	timeout time.Duration) {
	// Run anonymous coroutine
	go func() {
		time_diff := time.Duration(timeout)

		// Monitor channel communications
		for {
			select {
			case packet := <-queue:
				// Calculate time difference
				time_diff = time.Now().Sub(packet.Timestamp)

				// Edge case: Timeout
				if time_diff < timeout {
					time.Sleep(timeout - time_diff)
				}

				fpacket, ok := ipMeta.find(packet)
				// Edge case: Error finding packet
				if !ok {
					continue
				}

				// Packet finding process okay. Check state
				if fpacket.ExpectedRToLZR != packet.ExpectedRToLZR {
					// No state change
					continue
				} else {
					timeoutIncoming <- packet
				}
			}
		}
	}()
}

/*
PollTimeoutRoutine

Create and spawn timeout and retransmit routines for packet channel processing.

NOTE: pState is dependent on `concurrentMap.go`
*/
func PollTimeoutRoutine(ipMeta *rCMap, timeoutQueue chan *packet_metadata, retransmitQueue chan *packet_metadata,
	workers int, timeoutT int, timeoutR int) chan *packet_metadata {
	// Convert duration time to second for interval routine processing
	TIMEOUT_T := time.Duration(timeoutT) * time.Second
	TIMEOUT_R := time.Duration(timeoutR) * time.Second

	// Create channel routine to receive timeout and retransmit packets
	timeout_routine := make(chan *packet_metadata, in_buf)

	// Spawn routines to poll in timeout & retransmit queues given intervals
	timeoutAlg(ipMeta, timeoutQueue, timeout_routine, TIMEOUT_T)
	timeoutAlg(ipMeta, retransmitQueue, timeout_routine, TIMEOUT_R)

	return timeout_routine
}

/*
ConstructRetransmitQueue
*/
func ConstructRetransmitQueue(workers int) chan *packet_metadata {
	retransmit_queue := make(chan *packet_metadata, in_buf)
	return retransmit_queue
}

/*
ConstructTimeoutQueue
*/
func ConstructTimeoutQueue(workers int) chan *packet_metadata {
	timeout_queue := make(chan *packet_metadata, in_buf)
	return timeout_queue
}
