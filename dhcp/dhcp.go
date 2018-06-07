package dhcp

import "net"

type Device struct {
	Name string
	FQDN string
	MAC  net.HardwareAddr
	IP   net.IP
}
