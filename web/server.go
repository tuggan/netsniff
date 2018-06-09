package web

import (
	"html/template"
	"log"
	"net/http"
	"netsniff/dhcp"
)

type server struct {
	devices       map[string]dhcp.Device
	indextemplate *template.Template
}

type indexpagedata struct {
	PageTitle string
	Devices   []dhcp.Device
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var devices []dhcp.Device
	for _, val := range s.devices {
		devices = append(devices, val)
	}
	data := indexpagedata{PageTitle: "Index", Devices: devices}
	s.indextemplate.Execute(w, data)
	//fmt.Fprintf(w, "hello, you've hit %s\n", r.URL.Path)
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s requested %s", r.RemoteAddr, r.URL)
		h.ServeHTTP(w, r)
	})
}

func Serve(devices map[string]dhcp.Device) {

	fs := http.FileServer(http.Dir("html/static/"))
	h := http.NewServeMux()

	tmpl := template.Must(template.ParseFiles("web/index.html"))

	h.Handle("/static/", http.StripPrefix("/static/", fs))
	h.Handle("/", server{devices: devices, indextemplate: tmpl})

	hl := logger(h)

	err := http.ListenAndServe(":8080", hl)

	log.Fatal(err)
}
