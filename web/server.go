package web

import (
	"html/template"
	"log"
	"net/http"
	"netsniff/dhcp"
)

type server struct {
	devices map[string]dhcp.Device
}

type indexpagedata struct {
	PageTitle string
	Devices   []dhcp.Device
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//for key, val := range s.devices {
	//	fmt.Fprintf(w,
	//		"Key: %s\tName: %s\tFQDN: %s\tMAC: %s\tIP: %s\n",
	//		key,
	//		val.Name,
	//		val.FQDN,
	//		val.MAC.String(),
	//		val.IP.String())
	//}

	tmpl := template.Must(template.ParseFiles("web/index.html"))
	var devices []dhcp.Device
	for _, val := range s.devices {
		devices = append(devices, val)
	}
	data := indexpagedata{PageTitle: "Index", Devices: devices}
	tmpl.Execute(w, data)
	//fmt.Fprintf(w, "hello, you've hit %s\n", r.URL.Path)
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s requested %s", r.RemoteAddr, r.URL)
		h.ServeHTTP(w, r)
	})
}

func Serve(devices map[string]dhcp.Device) {
	h := http.NewServeMux()
	h.Handle("/", server{devices: devices})
	hl := logger(h)
	err := http.ListenAndServe(":8080", hl)
	log.Fatal(err)
}
