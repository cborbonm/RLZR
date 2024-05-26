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

func closeConnection(packet *packet_metadata, ipMeta *rCMap, writingQueue chan packet_metadata, write bool, ackingFirewall bool) {
	rst := constructRST(packet)
	err := handle.WritePacketData(rst)
	if err != nil {
		log.Fatal(err)
	}
	packet = ipMeta.remove(packet)
	if write {
		packet.setHyperACKtive(ackingFirewall)
		writingQueue <- *packet
	}
}

func HandlePcap(opts *options, packet *packet_metadata, ipMeta *rCMap, timeoutQueue chan *packet_metadata, retransmitQueue chan *packet_metadata, writingQueue chan packet_metadata) {
	verified := ipMeta.verifyScanningIP(packet)
	if !verified {
		packet.incrementCounter()
		packet.updateTimestamp()
		packet.validationFail()
		timeoutQueue <- packet
		return
	}

	isHyperACKtive := ipMeta.getHyperACKtiveStatus(packet)
	handshakeNum := ipMeta.getHandshake(packet)

	if !packet.SYN && packet.ACK {
		ipMeta.updateAck(packet)
	}

	if len(packet.Data) > 0 {
		packet.updateResponse("DATA")
		ipMeta.updateData(packet)

		if ForceAllHandshakes() {
			handleExpired(opts, packet, ipMeta, timeoutQueue, writingQueue)
			return
		}

		packet.syncHandshakeNum(handshakeNum)
		closeConnection(packet, ipMeta, writingQueue, true, isHyperACKtive)
		return
	}

	if packet.RST || packet.FIN {
		handleExpired(opts, packet, ipMeta, timeoutQueue, writingQueue)
		return
	}

	if handshakeNum == 1 && HyperACKtiveFiltering() && !isHyperACKtive {
		if ipMeta.getEphemeralRespNum(packet) > getNumFilters() {
			closeConnection(packet, ipMeta, writingQueue, true, true)
			return
		}
	}

	if !packet.SYN && packet.ACK {
		packet.updateResponse("DATA")
		packet.updateTimestamp()
		ipMeta.update(packet)
		timeoutQueue <- packet
		return
	}

	if packet.SYN && packet.ACK {
		if handshakeNum == 1 && HyperACKtiveFiltering() {
			if isHyperACKtive {
				parentSport := ipMeta.getParentSport(packet)
				ipMeta.incEphemeralResp(packet, parentSport)
				closeConnection(packet, ipMeta, writingQueue, false, isHyperACKtive)
				return
			} else {
				ipMeta.incEphemeralResp(packet, packet.Sport)
			}
		}
		toACK := true
		toPUSH := false
		SendAck(opts, packet, ipMeta, timeoutQueue, retransmitQueue, writingQueue, toACK, toPUSH, "ACK")
	}
}
