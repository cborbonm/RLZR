package vnc

import (
  "rlzr"
	"strings"
)

// Handshake implements the lzr.Handshake interface
type HandshakeMod struct {
}

//server first protocol
func (h *HandshakeMod) GetData( dst string ) []byte {
    return []byte("")
}

func (h *HandshakeMod) Verify( data string ) string {
	if strings.Contains( data, "RFB" ) {
         return "vnc"
	}
    return ""
}

func RegisterHandshake() {
	var h HandshakeMod
	rlzr.AddHandshake( "vnc", &h )
}

