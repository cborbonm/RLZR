package rlzr

import (
	"log"
	//"time"
)

func SendAck( opts *options, synack  *packet_metadata, ipMeta * rCMap,
timeoutQueue  chan *packet_metadata, retransmitQueue chan *packet_metadata,
writingQueue  chan packet_metadata, toACK bool, toPUSH bool, expectedResponse string ) {


	if synack.windowZero() {
		//not a real s/a
		writingQueue <- *synack
		return
	}

	//grab which handshake
	handshakeNum := ipMeta.getHandshake(synack)
	handshake, _ := GetHandshake( opts.Handshakes[ handshakeNum ] )

	//Send Ack with Data
	ack, payload := constructData( handshake, synack, toACK, toPUSH )//true, false )
	//add to map
	synack.updateResponse( expectedResponse )//ACK )
	synack.updateResponseL( payload )
	synack.updateTimestamp()
	ipMeta.update( synack )
	err := handle.WritePacketData(ack)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	synack.updateTimestamp()
	retransmitQueue <-synack
	return

}


