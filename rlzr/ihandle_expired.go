/*
Copyright 2020 The Board of Trustees of The Leland Stanford Junior University

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package rlzr

import (
	//"fmt"
)

func handleExpired( opts *options, packet * packet_metadata, ipMeta * pState,
	timeoutQueue chan *packet_metadata, writingQueue chan packet_metadata ) {


	// first close the existing connection unless
	// its already been terminated
	if !( packet.RST && !packet.ACK ) {

		rst := constructRST( packet )
		_ = handle.WritePacketData( rst )

	}

	//grab which handshake
	handshakeNum := ipMeta.getHandshake( packet )

	//if we are all not trying anymore handshakes, so sad. 
	//if ( packet.ExpectedRToLZR == SYN_ACK ||
	if ( packet.HyperACKtive  || (handshakeNum >= (len( opts.Handshakes ) - 1)) ){

		packet.syncHandshakeNum( handshakeNum )

		//document failure if its a handshake response that hasnt succeeded before
		if !packet.HyperACKtive && !( ForceAllHandshakes() && ipMeta.getData( packet ) && len(packet.Data) == 0 ) {
			writingQueue <- *packet
		}

		//remove from state, we are done now
		packet = ipMeta.remove( packet )
		if HyperACKtiveFiltering() {
			packet.HyperACKtive = true
			packet = ipMeta.remove( packet )
		}
	} else { // lets try another handshake


		//record all succesful fingerprints if forcing all handshakes
		if ForceAllHandshakes() && len(packet.Data) > 0 {
			packet.syncHandshakeNum( handshakeNum )
			writingQueue <- *packet
		}

		packet.updatePacketFlow()
		ipMeta.incHandshake( packet )
		SendSyn( packet, ipMeta, timeoutQueue )

		//lets also filter for cananda-like things
		if ( handshakeNum == 0 &&  HyperACKtiveFiltering() ) {
			for i := 0; i < getNumFilters(); i++ {
				highPortPacket := createFilterPacket( packet )
				SendSyn( highPortPacket, ipMeta, timeoutQueue )
				ipMeta.incHandshake( highPortPacket )
				ipMeta.setHyperACKtiveStatus( highPortPacket )

				ipMeta.setParentSport( highPortPacket, packet.Sport )
				ipMeta.FinishProcessing( highPortPacket )
			}
		}

	}

}

