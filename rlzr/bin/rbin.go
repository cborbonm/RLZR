package bin

import "fmt"

func RLZRMain() {
	fmt.Println("Hello world!")
	// set up

	// read from Zmap

	// window size 0? exit

	// from ephemeral port?
	// if yes, num_received >= eph_limit? exit

	// if no (to either of last two questions),
	// SEND ACK HANDSHAKE[i]

	// 1) receive ack and data? try all fingerprinting modules

	// 2) no ack? if possible, retransmit (depending on # of retransmits left):
	// SEND ACK W/ PSH W/ HANDSHAKE[i]

	// 3) receive ack BUT no data?
	// no RST, no FIN, and retransmit possible?
	// SEND ACK W/ PSH W/ HANDSHAKE[i]

	// no RST but receive FIN? or no RST, no FIN and no retransmit possible?
	// send RST

	// if RST sent or RST received:
	// if more handshakes, i++, restart cycle

}
