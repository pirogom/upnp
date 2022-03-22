package upnp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// PIROGOM
const (
	defaultMappingDesc   = "mop"
	defailtLeaseDuration = 1800
)

/**
*	AddPortMapping
**/
type AddPortMapping struct {
	upnp *Upnp
}

/**
*	Send
* 	PIROGOM
**/
func (this *AddPortMapping) Send(localPort, remotePort int, protocol string, mappingDesc string, leaseDur int) bool {
	request := this.buildRequest(localPort, remotePort, protocol, mappingDesc, leaseDur)
	if request == nil {
		return false
	}
	response, responseErr := http.DefaultClient.Do(request)
	if responseErr != nil {
		return false
	}
	resultBody, resultBodyErr := ioutil.ReadAll(response.Body)
	if resultBodyErr != nil {
		return false
	}
	if response.StatusCode == 200 {
		this.resolve(string(resultBody))
		return true
	}
	return false
}

/**
*	buildRequest
* 	PIROGOM
**/
func (this *AddPortMapping) buildRequest(localPort, remotePort int, protocol string, mappingDesc string, leaseDur int) *http.Request {

	// PIROGOM
	if mappingDesc == "" {
		mappingDesc = defaultMappingDesc
	}
	if leaseDur <= 0 {
		leaseDur = defailtLeaseDuration
	}

	header := http.Header{}
	header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	header.Set("SOAPAction", `"urn:schemas-upnp-org:service:WANIPConnection:1#AddPortMapping"`)
	header.Set("Content-Type", "text/xml")
	header.Set("Connection", "Close")
	header.Set("Content-Length", "")

	body := Node{Name: "SOAP-ENV:Envelope",
		Attr: map[string]string{"xmlns:SOAP-ENV": `"http://schemas.xmlsoap.org/soap/envelope/"`,
			"SOAP-ENV:encodingStyle": `"http://schemas.xmlsoap.org/soap/encoding/"`}}
	childOne := Node{Name: `SOAP-ENV:Body`}
	childTwo := Node{Name: `m:AddPortMapping`,
		Attr: map[string]string{"xmlns:m": `"urn:schemas-upnp-org:service:WANIPConnection:1"`}}

	childList1 := Node{Name: "NewExternalPort", Content: strconv.Itoa(remotePort)}
	childList2 := Node{Name: "NewInternalPort", Content: strconv.Itoa(localPort)}
	childList3 := Node{Name: "NewProtocol", Content: protocol}
	childList4 := Node{Name: "NewEnabled", Content: "1"}
	childList5 := Node{Name: "NewInternalClient", Content: this.upnp.LocalHost}
	childList6 := Node{Name: "NewLeaseDuration", Content: fmt.Sprintf("%d", leaseDur)} // PIROGOM
	childList7 := Node{Name: "NewPortMappingDescription", Content: mappingDesc}        // PIROGOM
	childList8 := Node{Name: "NewRemoteHost"}
	childTwo.AddChild(childList1)
	childTwo.AddChild(childList2)
	childTwo.AddChild(childList3)
	childTwo.AddChild(childList4)
	childTwo.AddChild(childList5)
	childTwo.AddChild(childList6)
	childTwo.AddChild(childList7)
	childTwo.AddChild(childList8)

	childOne.AddChild(childTwo)
	body.AddChild(childOne)
	bodyStr := body.BuildXML()

	request, requestErr := http.NewRequest("POST", "http://"+this.upnp.Gateway.Host+this.upnp.CtrlUrl, strings.NewReader(bodyStr))
	if requestErr != nil {
		return nil
	}
	request.Header = header
	request.Header.Set("Content-Length", strconv.Itoa(len([]byte(bodyStr))))
	return request
}

/**
*	resolve
**/
func (this *AddPortMapping) resolve(resultStr string) {
}
