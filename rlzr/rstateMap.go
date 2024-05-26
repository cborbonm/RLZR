package rlzr

import (
	"strconv"
	"fmt"
)

func ConstructPacketStateMap( opts *options ) rCMap {
	ipMeta := New()
	return ipMeta
}

//type rCMap struct {
//    stateMap map[string]*packet_state
//}

//func NewpState() rCMap {
//    return rCMap{stateMap: make(map[string]*packet_state)}
//}

// construct key functions

func constructKey(packet *packet_metadata) string {
    return packet.Saddr + ":" + strconv.Itoa(packet.Sport)
}

func constructParentKey(packet *packet_metadata, parentSport int) string {
    return packet.Saddr + ":" + strconv.Itoa(parentSport)
}

// state management functions

func (ipMeta *rCMap) metaContains(p *packet_metadata) bool {
    pKey := constructKey(p)
    _, ok := ipMeta.Get(pKey)
    return ok
}

// func (ipMeta *rCMap) find(p *packet_metadata) (*packet_metadata, bool) {
// 	return nil, true
// }

func (ipMeta *rCMap) find(p *packet_metadata) (*packet_metadata, bool) {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        return ps.Packet, ok
    }
    return nil, ok
}

func (ipMeta *rCMap) update(p *packet_metadata) {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if !ok {
        ps = &packet_state{
            Packet:       p,
            Ack:          false,
            HandshakeNum: 0,
        }
    } else {
        ps.Packet = p
    }
    ipMeta.Set(pKey, ps)
}

func (ipMeta *rCMap) incHandshake(p *packet_metadata) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        ps.HandshakeNum++
        ipMeta.Set(pKey, ps)
    }
    return ok
}

func (ipMeta *rCMap) updateAck(p *packet_metadata) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        ps.Ack = true
        ipMeta.Set(pKey, ps)
    }
    return ok
}

func (ipMeta *rCMap) getAck(p *packet_metadata) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        return ps.Ack
    }
    return false
}

func (ipMeta *rCMap) incEphemeralResp(p *packet_metadata, sport int) bool {
    pKey := constructParentKey(p, sport)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        ps.EphemeralRespNum++
        ipMeta.Set(pKey, ps)
    }
    return ok
}

func (ipMeta *rCMap) getEphemeralRespNum(p *packet_metadata) int {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        return ps.EphemeralRespNum
    }
    return 0
}

func (ipMeta *rCMap) getHyperACKtiveStatus(p *packet_metadata) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        return ps.HyperACKtive
    }
    return false
}

func (ipMeta *rCMap) setHyperACKtiveStatus(p *packet_metadata) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        ps.HyperACKtive = true
        ipMeta.Set(pKey, ps)
    }
    return ok
}

func (ipMeta *rCMap) setParentSport(p *packet_metadata, sport int) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        ps.ParentSport = sport
        ipMeta.Set(pKey, ps)
    }
    return ok
}

func (ipMeta *rCMap) getParentSport(p *packet_metadata) int {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        return ps.ParentSport
    }
    return 0
}

func (ipMeta *rCMap) recordEphemeral(p *packet_metadata, ephemerals []packet_metadata) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        ps.EphemeralFilters = append(ps.EphemeralFilters, ephemerals...)
        ipMeta.Set(pKey, ps)
    }
    return ok
}

func (ipMeta *rCMap) getEphemeralFilters(p *packet_metadata) ([]packet_metadata, bool) {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        return ps.EphemeralFilters, ok
    }
    return nil, ok
}

func (ipMeta *rCMap) updateData(p *packet_metadata) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        ps.Data = true
	ipMeta.Set(pKey, ps)
    }
    return ok
}

func (ipMeta *rCMap) getData(p *packet_metadata) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        return ps.Data
    }
    return false
}

func (ipMeta *rCMap) getHandshake(p *packet_metadata) int {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if ok {
        return ps.HandshakeNum
    }
    return 0
}

func (ipMeta *rCMap) incrementCounter(p *packet_metadata) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.Get(pKey)
    if !ok {
        return false
    }
    ps.Packet.incrementCounter()
    ipMeta.Set(pKey, ps)
    return true
}

func (ipMeta *rCMap) remove(packet *packet_metadata) *packet_metadata {
    packet.ACKed = ipMeta.getAck(packet)
    packetKey := constructKey(packet)
    ipMeta.Remove(packetKey)
    return packet
}

// Verify if the packet is a response to a SYN-ACK

func verifySA(pMap *packet_metadata, pRecv *packet_metadata) bool {
    if pRecv.SYN && pRecv.ACK {
        if pRecv.Acknum == pMap.Seqnum+1 {
            return true
        }
    } else {
        if pRecv.Seqnum == pMap.Seqnum || pRecv.Seqnum == pMap.Seqnum+1 {
            if pRecv.Acknum == pMap.Acknum+pMap.LZRResponseL || pRecv.Acknum == 0 {
                return true
            }
        }
    }
    return false
}

func (ipMeta *rCMap) verifyScanningIP(pRecv *packet_metadata) bool {
    pRecvKey := constructKey(pRecv)
    ps, ok := ipMeta.Get(pRecvKey)
    if !ok {
        return false
    }
    pMap := ps.Packet
    if pMap.Saddr == pRecv.Saddr && pMap.Dport == pRecv.Dport && pMap.Sport == pRecv.Sport {
        if verifySA(pMap, pRecv) {
            return true
        }
    }

    if DebugOn() {
        fmt.Println(pMap.Saddr, "====")
        fmt.Println("recv seq num:", pRecv.Seqnum)
        fmt.Println("stored seqnum: ", pMap.Seqnum)
        fmt.Println("recv ack num:", pRecv.Acknum)
        fmt.Println("stored acknum: ", pMap.Acknum)
        fmt.Println("received response length: ", len(pRecv.Data))
        fmt.Println("stored response length: ", pMap.LZRResponseL)
        fmt.Println(pMap.Saddr, "====")
    }
    return false
}

// TEST

