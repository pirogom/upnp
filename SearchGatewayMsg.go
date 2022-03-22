package upnp

import (
	"log"
	"net"
	"strings"
	"time"
)

/**
*	Gateway
**/
type Gateway struct {
	GatewayName   string // Gateway Name
	Host          string //	Gateway IP Address and Port
	DeviceDescUrl string // Gateway description URL
	Cache         string // Cache
	ST            string //
	USN           string //
	deviceType    string // Device URN  "urn:schemas-upnp-org:service:WANIPConnection:1"
	ControlURL    string // Device port mapping request URL
	ServiceType   string // UPnP Service Type
}

/**
*	SearchGateway
**/
type SearchGateway struct {
	searchMessage string
	upnp          *Upnp
}

/**
*	Send
**/
func (this *SearchGateway) Send() bool {
	this.buildRequest()
	c := make(chan string)
	go this.send(c)
	result := <-c
	if result == "" {
		// timeout
		this.upnp.Active = false
		return false
	}
	this.resolve(result)

	this.upnp.Gateway.ServiceType = "urn:schemas-upnp-org:service:WANIPConnection:1"
	this.upnp.Active = true
	return true
}

/**
*	send
**/
func (this *SearchGateway) send(c chan string) {
	// Send multicast message (239.255.255.250:1900)
	var conn *net.UDPConn
	defer func() {
		if r := recover(); r != nil {
			// with timeout
		}
	}()
	go func(conn *net.UDPConn) {
		defer func() {
			if r := recover(); r != nil {
				// without timeout
			}
		}()

		// timeout is three second
		time.Sleep(time.Second * 3)
		c <- ""
		conn.Close()
	}(conn)
	remotAddr, err := net.ResolveUDPAddr("udp", "239.255.255.250:1900")
	if err != nil {
		log.Println("Invalid multicast address", err.Error())
	}
	locaAddr, err := net.ResolveUDPAddr("udp", this.upnp.LocalHost+":")

	if err != nil {
		log.Println("Invalud local ip address", err.Error())
	}
	conn, err = net.ListenUDP("udp", locaAddr)
	defer conn.Close()
	if err != nil {
		log.Println("Monitoring UDP error", err.Error())
	}
	_, err = conn.WriteToUDP([]byte(this.searchMessage), remotAddr)
	if err != nil {
		log.Println("Multicast message write failed(1)", err.Error())
	}
	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		log.Println("Multicast message write failed(2)", err.Error())
	}

	result := string(buf[:n])
	c <- result
}

/**
*	buildRequest
**/
func (this *SearchGateway) buildRequest() {
	this.searchMessage = "M-SEARCH * HTTP/1.1\r\n" +
		"HOST: 239.255.255.250:1900\r\n" +
		"ST: urn:schemas-upnp-org:service:WANIPConnection:1\r\n" +
		"MAN: \"ssdp:discover\"\r\n" + "MX: 3\r\n\r\n"
}

/**
*	resolve
**/
func (this *SearchGateway) resolve(result string) {
	this.upnp.Gateway = &Gateway{}

	lines := strings.Split(result, "\r\n")
	for _, line := range lines {
		// split token ':'
		nameValues := strings.SplitAfterN(line, ":", 2)
		if len(nameValues) < 2 {
			continue
		}
		switch strings.ToUpper(strings.Trim(strings.Split(nameValues[0], ":")[0], " ")) {
		case "ST":
			this.upnp.Gateway.ST = nameValues[1]
		case "CACHE-CONTROL":
			this.upnp.Gateway.Cache = nameValues[1]
		case "LOCATION":
			urls := strings.Split(strings.Split(nameValues[1], "//")[1], "/")
			this.upnp.Gateway.Host = urls[0]
			this.upnp.Gateway.DeviceDescUrl = "/" + urls[1]
		case "SERVER":
			this.upnp.Gateway.GatewayName = nameValues[1]
		default:
		}
	}
}
