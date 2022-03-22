package upnp

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
)

/**
*	DeviceDesc
**/
type DeviceDesc struct {
	upnp *Upnp
}

/**
*	Send
**/
func (this *DeviceDesc) Send() bool {
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
func (this *DeviceDesc) BuildRequest() *http.Request {
	// header
	header := http.Header{}
	header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	header.Set("User-Agent", "preston")
	header.Set("Host", this.upnp.Gateway.Host)
	header.Set("Connection", "keep-alive")

	// request
	request, requestErr := http.NewRequest("GET", "http://"+this.upnp.Gateway.Host+this.upnp.Gateway.DeviceDescUrl, nil)
	if requestErr != nil {
		return nil
	}
	request.Header = header
	// request := http.Request{Method: "GET", Proto: "HTTP/1.1",
	// 	Host: this.upnp.Gateway.Host, Url: this.upnp.Gateway.DeviceDescUrl, Header: header}
	return request
}

/**
*	resolve
**/
func (this *DeviceDesc) resolve(resultStr string) {
	inputReader := strings.NewReader(resultStr)

	// read from file：
	// content, err := ioutil.ReadFile("studygolang.xml")
	// decoder := xml.NewDecoder(bytes.NewBuffer(content))

	lastLabel := ""

	ISUpnpServer := false

	IScontrolURL := false
	var controlURL string //`controlURL`
	// var eventSubURL string //`eventSubURL`
	// var SCPDURL string     //`SCPDURL`

	decoder := xml.NewDecoder(inputReader)
	for t, err := decoder.Token(); err == nil && !IScontrolURL; t, err = decoder.Token() {
		switch token := t.(type) {
		// element processing start
		case xml.StartElement:
			if ISUpnpServer {
				name := token.Name.Local
				lastLabel = name
			}

		// element processing end
		case xml.EndElement:
			// PIROGOM
			//log.Println("End Tag :", token.Name.Local)
		// char(string) data processing
		case xml.CharData:
			// URL을 가져온 후 다른 태그는 처리되지 않습니다
			content := string([]byte(token))

			// 포트 매핑을 제공하는 서비스를 찾습니다
			if content == this.upnp.Gateway.ServiceType {
				ISUpnpServer = true
				continue
			}
			//urn:upnp-org:serviceId:WANIPConnection
			if ISUpnpServer {
				switch lastLabel {
				case "controlURL":
					controlURL = content
					IScontrolURL = true
				case "eventSubURL":
					// eventSubURL = content
				case "SCPDURL":
					// SCPDURL = content
				}
			}
		default:
			// ...
		}
	}
	this.upnp.CtrlUrl = controlURL
}
