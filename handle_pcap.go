package lzr

import (
	//"fmt"
	"log"
)


func closeConnection( packet *packet_metadata, ipMeta * pState, writingQueue chan packet_metadata, write bool, ackingFirewall bool ) {

	//close connection
	rst := constructRST(packet)
	err := handle.WritePacketData(rst)
	if err != nil {
		log.Fatal(err)
	}
	//remove from state, we are done now
	packet = ipMeta.remove(packet)
	if write {
		packet.setHyperACKtive(ackingFirewall)
		writingQueue <- *packet
	}
	return
}


func HandlePcap( opts *options, packet *packet_metadata, ipMeta * pState, timeoutQueue	chan *packet_metadata,
	retransmitQueue chan *packet_metadata, writingQueue chan packet_metadata ) {


	//verify
	verified := ipMeta.verifyScanningIP( packet )
	if !verified {
		packet.incrementCounter()
		packet.updateTimestamp()
		packet.validationFail()
		timeoutQueue <-packet
		return
	}

	isHyperACKtive := ipMeta.getHyperACKtiveStatus( packet )

	//for every ack received, mark as accepting data
	if (!packet.SYN) && packet.ACK {
		ipMeta.updateAck( packet )
	}
	 //exit condition
	 if len(packet.Data) > 0 {
		packet.updateResponse(DATA)
		ipMeta.updateData( packet )

		// if not stopping here, send off to handle_expire
		if ForceAllHandshakes() {
			handleExpired( opts,packet, ipMeta, timeoutQueue, writingQueue )
			return
		}

		handshakeNum := ipMeta.getHandshake( packet )
		packet.syncHandshakeNum( handshakeNum )

		/*if HyperACKtiveFiltering() && handshakeNum == 1 {
			cleanEphState( packet, ipMeta, writingQueue, false)
		}*/

		closeConnection( packet, ipMeta, writingQueue, true,  isHyperACKtive)
		return

	}
	//deal with closed connection 
	if packet.RST || packet.FIN {

		handleExpired( opts,packet, ipMeta, timeoutQueue, writingQueue )
		return

	 }


	//checking if max filter syn acks reached
	//( filterACKs + original ACK + this ack)
     if HyperACKtiveFiltering() && !isHyperACKtive {
			//fmt.Println( ipMeta.getEphemeralRespNum( packet ) )
			//fmt.Println(getNumFilters())
            if ipMeta.getEphemeralRespNum( packet )   > getNumFilters() {
                closeConnection( packet, ipMeta, writingQueue, true, true)
				return
            }
     }

	 //for every ack received, mark as accepting data
	 if (!packet.SYN) && packet.ACK {
		 //add to map
		 packet.updateResponse(DATA)
		 packet.updateTimestamp()
		 ipMeta.update(packet)

		 //add to map
		 timeoutQueue <-packet
		 return
	}

	//for every s/a send the appropriate ack
	if packet.SYN && packet.ACK {

		if  HyperACKtiveFiltering() {

			//just close and record
			if isHyperACKtive {

                parentSport := ipMeta.getParentSport( packet )

				ipMeta.incEphemeralResp( packet, parentSport )
				closeConnection( packet, ipMeta, writingQueue, false, isHyperACKtive)
				return
			} else {
				ipMeta.incEphemeralResp( packet, packet.Sport )
			}
		}

		SendAck( opts, packet, ipMeta, timeoutQueue, retransmitQueue, writingQueue )
		return
	}

}

