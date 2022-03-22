package upnp

import (
	"errors"
	"log"
	"sync"
)

/**
*	MappingPortStruct
**/
type MappingPortStruct struct {
	lock         *sync.Mutex
	mappingPorts map[string][][]int
}

/**
*	addMapping
**/
func (this *MappingPortStruct) addMapping(localPort, remotePort int, protocol string) {

	this.lock.Lock()
	defer this.lock.Unlock()
	if this.mappingPorts == nil {
		one := make([]int, 0)
		one = append(one, localPort)
		two := make([]int, 0)
		two = append(two, remotePort)
		portMapping := [][]int{one, two}
		this.mappingPorts = map[string][][]int{protocol: portMapping}
		return
	}
	portMapping := this.mappingPorts[protocol]
	if portMapping == nil {
		one := make([]int, 0)
		one = append(one, localPort)
		two := make([]int, 0)
		two = append(two, remotePort)
		this.mappingPorts[protocol] = [][]int{one, two}
		return
	}
	one := portMapping[0]
	two := portMapping[1]
	one = append(one, localPort)
	two = append(two, remotePort)
	this.mappingPorts[protocol] = [][]int{one, two}
}

/**
*	delMapping
**/
func (this *MappingPortStruct) delMapping(remotePort int, protocol string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.mappingPorts == nil {
		return
	}
	tmp := MappingPortStruct{lock: new(sync.Mutex)}
	mappings := this.mappingPorts[protocol]
	for i := 0; i < len(mappings[0]); i++ {
		if mappings[1][i] == remotePort {
			break
		}
		tmp.addMapping(mappings[0][i], mappings[1][i], protocol)
	}
	this.mappingPorts = tmp.mappingPorts
}
func (this *MappingPortStruct) GetAllMapping() map[string][][]int {
	return this.mappingPorts
}

/**
*	Upnp
**/
type Upnp struct {
	Active             bool              // for UPnP device active
	LocalHost          string            // Local IP Address
	GatewayInsideIP    string            // LAN Gateway IP Address
	GatewayOutsideIP   string            // Gateway Public IP Address
	OutsideMappingPort map[string]int    // Outside Mapping Port
	InsideMappingPort  map[string]int    // Inside(default) Mapping Port
	Gateway            *Gateway          // Gateway Info
	CtrlUrl            string            // Control URL
	MappingPort        MappingPortStruct // Mapping Port Info {"TCP":[1990],"UDP":[1991]}
}

/**
*	SearchGateway
**/
func (this *Upnp) SearchGateway() (err error) {
	defer func(err error) {
		if errTemp := recover(); errTemp != nil {
			log.Println("upnp error", errTemp)
			err = errTemp.(error)
		}
	}(err)

	if this.LocalHost == "" {
		this.MappingPort = MappingPortStruct{
			lock: new(sync.Mutex),
			// mappingPorts: map[string][][]int{},
		}
		this.LocalHost = GetLocalIntenetIp()
	}
	searchGateway := SearchGateway{upnp: this}
	if searchGateway.Send() {
		return nil
	}
	return errors.New("Gateway not found")
}

/**
*	deviceStatus
**/
func (this *Upnp) deviceStatus() {

}

/**
*	deviceDesc
**/
func (this *Upnp) deviceDesc() (err error) {
	if this.GatewayInsideIP == "" {
		if err := this.SearchGateway(); err != nil {
			return err
		}
	}
	device := DeviceDesc{upnp: this}
	device.Send()
	this.Active = true
	//log.Println("Control Url:", this.CtrlUrl)
	return
}

/**
*	ExternalIPAddr
**/
func (this *Upnp) ExternalIPAddr() (err error) {
	if this.CtrlUrl == "" {
		if err := this.deviceDesc(); err != nil {
			return err
		}
	}
	eia := ExternalIPAddress{upnp: this}
	eia.Send()
	return nil
	// log.Println("Public IP Address：", this.GatewayOutsideIP)
}

/**
*	AddPortMapping
* 	PIROGOM
**/
func (this *Upnp) AddPortMapping(localPort, remotePort int, protocol string, mappingDesc string, leaseDur int) (err error) {
	defer func(err error) {
		if errTemp := recover(); errTemp != nil {
			log.Println("UPnP Error", errTemp)
			err = errTemp.(error)
		}
	}(err)
	if this.GatewayOutsideIP == "" {
		if err := this.ExternalIPAddr(); err != nil {
			return err
		}
	}
	addPort := AddPortMapping{upnp: this}
	if issuccess := addPort.Send(localPort, remotePort, protocol, mappingDesc, leaseDur); issuccess {
		this.MappingPort.addMapping(localPort, remotePort, protocol)
		// log.Println("Port Mapping Added：protocol:", protocol, "local:", localPort, "remote:", remotePort)
		return nil
	} else {
		this.Active = false
		// log.Println("Add Port Mapping failed")
		return errors.New("Add Port Mapping failed")
	}
}

/**
*	DelPortMapping
**/
func (this *Upnp) DelPortMapping(remotePort int, protocol string) bool {
	delMapping := DelPortMapping{upnp: this}
	issuccess := delMapping.Send(remotePort, protocol)
	if issuccess {
		this.MappingPort.delMapping(remotePort, protocol)
		log.Println("Delete Port Mapping： remote:", remotePort)
	}
	return issuccess
}

/**
*	Reclaim
**/
func (this *Upnp) Reclaim() {
	mappings := this.MappingPort.GetAllMapping()
	tcpMapping, ok := mappings["TCP"]
	if ok {
		for i := 0; i < len(tcpMapping[0]); i++ {
			this.DelPortMapping(tcpMapping[1][i], "TCP")
		}
	}
	udpMapping, ok := mappings["UDP"]
	if ok {
		for i := 0; i < len(udpMapping[0]); i++ {
			this.DelPortMapping(udpMapping[0][i], "UDP")
		}
	}
}

/**
*	GetAllMapping
**/
func (this *Upnp) GetAllMapping() map[string][][]int {
	return this.MappingPort.GetAllMapping()
}
