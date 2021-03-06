package upnp

import (
	"net/http"
	"strconv"
	"strings"
)

/**
*	SearchGatewayReq
**/
type SearchGatewayReq struct {
	host       string
	resultBody string
	ctrlUrl    string
	upnp       *Upnp
}

/**
*	Send
**/
func (this SearchGatewayReq) Send() {
	// request := this.BuildRequest()
}

/**
*	BuildRequest
**/
func (this SearchGatewayReq) BuildRequest() *http.Request {
	// header
	header := http.Header{}
	header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	header.Set("SOAPAction", `"urn:schemas-upnp-org:service:WANIPConnection:1#GetStatusInfo"`)
	header.Set("Content-Type", "text/xml")
	header.Set("Connection", "Close")
	header.Set("Content-Length", "")
	// body
	body := Node{Name: "SOAP-ENV:Envelope",
		Attr: map[string]string{"xmlns:SOAP-ENV": `"http://schemas.xmlsoap.org/soap/envelope/"`,
			"SOAP-ENV:encodingStyle": `"http://schemas.xmlsoap.org/soap/encoding/"`}}
	childOne := Node{Name: `SOAP-ENV:Body`}
	childTwo := Node{Name: `m:GetStatusInfo`,
		Attr: map[string]string{"xmlns:m": `"urn:schemas-upnp-org:service:WANIPConnection:1"`}}
	childOne.AddChild(childTwo)
	body.AddChild(childOne)
	bodyStr := body.BuildXML()
	// request
	request, requestErr := http.NewRequest("POST", "http://"+this.upnp.Gateway.Host+this.upnp.CtrlUrl, strings.NewReader(bodyStr))
	if requestErr != nil {
		return nil
	}
	request.Header = header
	request.Header.Set("Content-Length", strconv.Itoa(len([]byte(bodyStr))))
	return request
}
