package web

import (
	"fmt"
	"net/http"
	"netsniff/dhcp"
)

type server struct {
	devices map[string]dhcp.Device
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for key, val := range s.devices {
		fmt.Fprintf(w,
			"Key: %s\tName: %s\tFQDN: %s\tMAC: %s\tIP: %s\n",
			key,
			val.Name,
			val.FQDN,
			val.MAC.String(),
			val.IP.String())
	}
	fmt.Fprintf(w, "hello, you've hit %s\n", r.URL.Path)
}

func Serve(devices map[string]dhcp.Device) {
	http.ListenAndServe(":8080", server{devices: devices})
}
