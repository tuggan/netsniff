package ping

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/tatsushid/go-fastping"
	"log"
	"net"
	"netsniff/dhcp"
	"os"
	"time"
)

type Sweeper struct {
	Devices map[string]dhcp.Device
	Nets    []net.IPNet
}

func inc(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}

func (s Sweeper) Sweep() {
	dc := dns.Client{}
	dm := dns.Msg{}
	for _, netw := range s.Nets {
		log.Printf("Starting sweep on: %s", netw.String())
		for ip, _, _ := net.ParseCIDR(netw.String()); netw.Contains(ip); inc(ip) {
			//log.Printf("IP: %s, Contains: %t, n.IP: %s, n.Mask: %s", ip.String(), netw.Contains(ip), netw.IP.String(), netw.Mask.String())
			p := fastping.NewPinger()
			ra, err := net.ResolveIPAddr("ip4:icmp", ip.String())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			p.AddIPAddr(ra)
			p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
				log.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)

				dm.SetQuestion(ra.String()+".", dns.TypePTR)
				lookup, t, err := dc.Exchange(&dm, "10.13.37.1:53")
				if err != nil {
					log.Println(err)
				} else {
					log.Printf("DNS request took %v", t)
					for _, s := range lookup.Answer {
						log.Println(s)
					}
				}

			}
			//p.OnIdle = func() {
			//	fmt.Println("finish")
			//}
			err = p.Run()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
