package lzr

import "time"

/*
Variables needed for main routines
*/
var (
	source_mac string
	dest_mac   string
)

/*
InitParams

Initialize `source_mac` and `dest_mac`

NOTE: Depends on `construct_responses.go`
*/
func InitParams() {
	source_mac = getSourceMacAddr()
	dest_mac = getHostMacAddr()
}

/*
ConstructWritingQueue

Takes `workers` as input to construct writing queue.
*/
func ConstructWritingQueue(workers int) chan packet_metadata {
	write_q := make(chan packet_metadata, 10000)
	return write_q
}

func ConstructIncomingRoutine(workers int) chan *packet_metadata {

}

func ConstructPcapRoutine(workers int) chan *packet_metadata {

}

// NOTE: pState is dependent on `concurrentMap.go`
func PollTimeoutRoutine(ipMeta *pState, timeoutQueue chan *packet_metadata, retransmitQueue chan *packet_metadata,
	workers int, timeoutT int, timeoutR int) chan *packet_metadata {

}

func timeoutAlg(ipMeta *pState, queue chan *packet_metadata, timeoutIncoming chan *packet_metadata,
	timeout time.Duration) {

}

func ConstructRetransmitQueue(workers int) chan *packet_metadata {

}

func ConstructTimeoutQueue(workers int) chan *packet_metadata {

}
