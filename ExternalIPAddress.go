package upnp

import (
	"encoding/xml"
	"io/ioutil"

	"net/http"
	"strconv"
	"strings"
)

/**
*	ExternalIPAddress
**/
type ExternalIPAddress struct {
	upnp *Upnp
}

/**
*	Send
**/
func (this *ExternalIPAddress) Send() bool {
	request := this.BuildRequest()
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
*	BuildRequest
**/
func (this *ExternalIPAddress) BuildRequest() *http.Request {
	// header
	header := http.Header{}
	header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	header.Set("SOAPAction", `"urn:schemas-upnp-org:service:WANIPConnection:1#GetExternalIPAddress"`)
	header.Set("Content-Type", "text/xml")
	header.Set("Connection", "Close")
	header.Set("Content-Length", "")
	// body
	body := Node{Name: "SOAP-ENV:Envelope",
		Attr: map[string]string{"xmlns:SOAP-ENV": `"http://schemas.xmlsoap.org/soap/envelope/"`,
			"SOAP-ENV:encodingStyle": `"http://schemas.xmlsoap.org/soap/encoding/"`}}
	childOne := Node{Name: `SOAP-ENV:Body`}
	childTwo := Node{Name: `m:GetExternalIPAddress`,
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
	request.Header.Set("Content-Length", strconv.Itoa(len([]byte(body.BuildXML()))))
	return request
}

/**
*	resolve
**/
//NewExternalIPAddress
func (this *ExternalIPAddress) resolve(resultStr string) {
	inputReader := strings.NewReader(resultStr)
	decoder := xml.NewDecoder(inputReader)
	ISexternalIP := false
	for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		// element processing start
		case xml.StartElement:
			name := token.Name.Local
			if name == "NewExternalIPAddress" {
				ISexternalIP = true
			}

		// element processing end
		case xml.EndElement:
		// char(string) data processing
		case xml.CharData:
			if ISexternalIP == true {
				this.upnp.GatewayOutsideIP = string([]byte(token))
				return
			}
		default:
			// ...
		}
	}
}
