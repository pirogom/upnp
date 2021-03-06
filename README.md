## pirogom/upnp
====

clone from [cokeboL/upnp](https://github.com/cokeboL/upnp)

grep PIROGOM

## upnp protocol
====

A simple implements UPnP protocol for Golang library.  Add port mapping for NAT devices.

====

## example: 

### 1. add a port mapping
~~~ go
mapping := new(upnp.Upnp)
if err := mapping.AddPortMapping(55789, 55789, "TCP", "mop", 1800); err == nil {
	fmt.Println("success !")
	// remove port mapping in gatway
	mapping.Reclaim()
} else {
	fmt.Println("fail !")
}
~~~

### 2. search gateway device.
~~~ go
upnpMan := new(upnp.Upnp)
err := upnpMan.SearchGateway()
if err != nil {
	fmt.Println(err.Error())
} else {
	fmt.Println("local ip address: ", upnpMan.LocalHost)
	fmt.Println("gateway ip address: ", upnpMan.Gateway.Host)
}
~~~
### 3. get an internet ip address in gatway.
~~~ go
upnpMan := new(upnp.Upnp)
err := upnpMan.ExternalIPAddr()
if err != nil {
	fmt.Println(err.Error())
} else {
	fmt.Println("internet ip address: ", upnpMan.GatewayOutsideIP)
}
~~~