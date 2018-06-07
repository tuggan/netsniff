package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/krolaw/dhcp4"
	"log"
	"net"
	"netsniff/dhcp"
	"netsniff/web"
	"os"
	"strings"
)

type DHCPHandler struct {
	ip      net.IP
	devices map[string]dhcp.Device
}

func (h *DHCPHandler) ServeDHCP(p dhcp4.Packet, msgType dhcp4.MessageType, options dhcp4.Options) (d dhcp4.Packet) {

	dev, ok := h.devices[p.CHAddr().String()]
	if ok == false {
		mac := make(net.HardwareAddr, len(p.CHAddr()))
		copy(mac, p.CHAddr())
		dev = dhcp.Device{MAC: mac}
	}

	switch msgType {
	case dhcp4.Discover:
		fmt.Println("DHCP Discover")
	case dhcp4.Request:
		fmt.Println("DHCP Request")
		if options[dhcp4.OptionRequestedIPAddress] != nil {
			//fmt.Println("Requested IP: " + string(options[dhcp4.OptionRequestedIPAddress]))
			ip := make(net.IP, len(options[dhcp4.OptionRequestedIPAddress]))
			copy(ip, options[dhcp4.OptionRequestedIPAddress])
			dev.IP = ip
		}
		if bytes.Compare(p.CIAddr(), net.IP{0, 0, 0, 0}) != 0 {
			fmt.Println("Using actual ip address: ", p.CIAddr().String())
			ip := make(net.IP, len(p.CIAddr()))
			copy(ip, p.CIAddr())
			dev.IP = ip
		}
	case dhcp4.Release:
		fmt.Println("DHCP Release")
	case dhcp4.Decline:
		fmt.Println("DHCP Decline")
	}
	//fmt.Println("Name: " + string(p.SName()))
	//fmt.Println("MAC: " + p.CHAddr().String())
	fmt.Println("CIAddr: " + p.CIAddr().String())
	fmt.Println("GIAddr: " + p.GIAddr().String())
	fmt.Println("SIAddr: " + p.SIAddr().String())
	fmt.Println("YIAddr: " + p.YIAddr().String())

	if options[dhcp4.OptionHostName] != nil {
		//fmt.Println("Hostname: " + string(options[dhcp4.OptionHostName]))
		dev.Name = string(options[dhcp4.OptionHostName])
	}

	//if options[dhcp4.OptionClientIdentifier] != nil {
	//	fmt.Println("Client Identifier: " + string(options[dhcp4.OptionClientIdentifier]))
	//}

	if options[81] != nil {
		//fmt.Println("FQDN: " + string(options[81]))
		fmt.Println(hex.Dump(options[81]))
		dev.FQDN = string(options[81])
	}

	h.devices[p.CHAddr().String()] = dev

	log.Printf("%s (%s) added\n", dev.Name, dev.MAC.String())
	//fmt.Println(hex.Dump(p))
	//fmt.Println(options)
	//fmt.Println("\n")

	return nil
}

func main() {
	serverIP := net.IP{0, 0, 0, 0}
	handler := &DHCPHandler{ip: serverIP, devices: make(map[string]dhcp.Device)}
	conn, err := net.ListenPacket("udp", ":67")
	if err != nil {
		os.Exit(1)
	}
	go dhcp4.Serve(conn, handler)
	//conn2, err := net.ListenPacket("udp", ":68")
	//if err != nil {
	//	os.Exit(2)
	//}
	//go dhcp4.Serve(conn2, handler)

	go web.Serve(handler.devices)

	reader := bufio.NewReader(os.Stdin)
	run := true
	for run {
		fmt.Print("Press Enter for status\n")
		text, _ := reader.ReadString('\n')
		if strings.Compare(text, "q") == 0 {
			run = false
		} else {
			for key, val := range handler.devices {
				fmt.Println("Key: " + key)
				fmt.Println("Name: " + val.Name)
				fmt.Println("FQDN: " + val.FQDN)
				fmt.Println("MAC: " + val.MAC.String())
				fmt.Println("IP: " + val.IP.String())
			}
		}
	}
}
