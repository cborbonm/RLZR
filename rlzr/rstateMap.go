package lzr

type packet_metadata struct {
    Saddr          string
    Sport          int
    Daddr          string
    Dport          int
    Seqnum         int
    Acknum         int
    SYN            bool
    ACK            bool
    Data           []byte
    LZRResponseL   int
    HyperACKtive   bool
    ParentSport    int
    EphemeralFilters []packet_metadata
}

type packet_state struct {
    Packet          *packet_metadata
    Ack             bool
    HandshakeNum    int
    EphemeralRespNum int
    HyperACKtive    bool
    Data            bool
}

type pState struct {
    stateMap map[string]*packet_state
}

func NewpState() pState {
    return pState{stateMap: make(map[string]*packet_state)}
}

// construct key functions

func constructKey(packet *packet_metadata) string {
    return packet.Saddr + ":" + strconv.Itoa(packet.Sport)
}

func constructParentKey(packet *packet_metadata, parentSport int) string {
    return packet.Saddr + ":" + strconv.Itoa(parentSport)
}

// state management functions

func (ipMeta *pState) metaContains(p *packet_metadata) bool {
    pKey := constructKey(p)
    _, ok := ipMeta.stateMap[pKey]
    return ok
}

// func (ipMeta *pState) find(p *packet_metadata) (*packet_metadata, bool) {
// 	return nil, true
// }

func (ipMeta *pState) find(p *packet_metadata) (*packet_metadata, bool) {
    pKey := constructKey(p)
    ps, ok := ipMeta.stateMap[pKey]
    if ok {
        return ps.Packet, ok
    }
    return nil, ok
}

func (ipMeta *pState) update(p *packet_metadata) {
    pKey := constructKey(p)
    ps, ok := ipMeta.stateMap[pKey]
    if !ok {
        ps = &packet_state{
            Packet:       p,
            Ack:          false,
            HandshakeNum: 0,
        }
    } else {
        ps.Packet = p
    }
    ipMeta.stateMap[pKey] = ps
}

func (ipMeta *pState) incHandshake(p *packet_metadata) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.stateMap[pKey]
    if ok {
        ps.HandshakeNum++
        ipMeta.stateMap[pKey] = ps
    }
    return ok
}

func (ipMeta *pState) updateAck(p *packet_metadata) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.stateMap[pKey]
    if ok {
        ps.Ack = true
        ipMeta.stateMap[pKey] = ps
    }
    return ok
}

func (ipMeta *pState) getAck(p *packet_metadata) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.stateMap[pKey]
    if ok {
        return ps.Ack
    }
    return false
}

func (ipMeta *pState) incEphemeralResp(p *packet_metadata, sport int) bool {
    pKey := constructParentKey(p, sport)
    ps, ok := ipMeta.stateMap[pKey]
    if ok {
        ps.EphemeralRespNum++
        ipMeta.stateMap[pKey] = ps
    }
    return ok
}

func (ipMeta *pState) getEphemeralRespNum(p *packet_metadata) int {
    pKey := constructKey(p)
    ps, ok := ipMeta.stateMap[pKey]
    if ok {
        return ps.EphemeralRespNum
    }
    return 0
}

func (ipMeta *pState) getHyperACKtiveStatus(p *packet_metadata) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.stateMap[pKey]
    if ok {
        return ps.HyperACKtive
    }
    return false
}

func (ipMeta *pState) setHyperACKtiveStatus(p *packet_metadata) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.stateMap[pKey]
    if ok {
        ps.HyperACKtive = true
        ipMeta.stateMap[pKey] = ps
    }
    return ok
}

func (ipMeta *pState) setParentSport(p *packet_metadata, sport int) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.stateMap[pKey]
    if ok {
        ps.ParentSport = sport
        ipMeta.stateMap[pKey] = ps
    }
    return ok
}

func (ipMeta *pState) getParentSport(p *packet_metadata) int {
    pKey := constructKey(p)
    ps, ok := ipMeta.stateMap[pKey]
    if ok {
        return ps.ParentSport
    }
    return 0
}

func (ipMeta *pState) recordEphemeral(p *packet_metadata, ephemerals []packet_metadata) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.stateMap[pKey]
    if ok {
        ps.EphemeralFilters = append(ps.EphemeralFilters, ephemerals...)
        ipMeta.stateMap[pKey] = ps
    }
    return ok
}

func (ipMeta *pState) getEphemeralFilters(p *packet_metadata) ([]packet_metadata, bool) {
    pKey := constructKey(p)
    ps, ok := ipMeta.stateMap[pKey]
    if ok {
        return ps.EphemeralFilters, ok
    }
    return nil, ok
}

func (ipMeta *pState) updateData(p *packet_metadata) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.stateMap[pKey]
    if ok {
        ps.Data = true
        ipMeta.stateMap[pKey] = ps
    }
    return ok
}

func (ipMeta *pState) getData(p *packet_metadata) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.stateMap[pKey]
    if ok {
        return ps.Data
    }
    return false
}

func (ipMeta *pState) getHandshake(p *packet_metadata) int {
    pKey := constructKey(p)
    ps, ok := ipMeta.stateMap[pKey]
    if ok {
        return ps.HandshakeNum
    }
    return 0
}

func (ipMeta *pState) incrementCounter(p *packet_metadata) bool {
    pKey := constructKey(p)
    ps, ok := ipMeta.stateMap[pKey]
    if !ok {
        return false
    }
    ps.Packet.incrementCounter()
    ipMeta.stateMap[pKey] = ps
    return true
}

func (ipMeta *pState) remove(packet *packet_metadata) *packet_metadata {
    packet.ACKed = ipMeta.getAck(packet)
    packetKey := constructKey(packet)
    delete(ipMeta.stateMap, packetKey)
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

func (ipMeta *pState) verifyScanningIP(pRecv *packet_metadata) bool {
    pRecvKey := constructKey(pRecv)
    ps, ok := ipMeta.stateMap[pRecvKey]
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

func main() {
    opts := &options{} // Populate with necessary options
    ipMeta := ConstructPacketStateMap(opts)
    
    // Add test packets
    pkt := &packet_metadata{
        Saddr: "192.168.0.1",
        Sport: 12345,
        Daddr: "192.168.0.2",
        Dport: 80,
        Seqnum: 100,
        Acknum: 200,
        SYN: true,
        ACK: false,
        Data: []byte{},
        LZRResponseL: 1,
    }
    
    ipMeta.update(pkt)
    
    // Verify packet
    result := ipMeta.verifyScanningIP(pkt)
    fmt.Println("Verification Result:", result)
}
