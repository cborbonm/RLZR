package rlzr

import (
	"sync"
)

type pStateShared struct {
	items        map[string]*packet_state
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

type pState []*pStateShared
