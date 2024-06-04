package handshakes

import "rlzr/handshakes/amqp"

func init() {
	amqp.RegisterHandshake()
}

