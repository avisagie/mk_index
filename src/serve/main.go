
package main

import (
	"log"
	"net/http"
	"flag"
	"net"
	"os"
)

var (
	listenAddr = flag.String("addr", ":8080", "Address to listen on")
)

func init() {
	flag.Parse()
}

func main() {
	h := http.FileServer(http.Dir("./"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s -> %s", r.RemoteAddr, r.URL.String())
		h.ServeHTTP(w, r)
	})
	
	// Print this machine's interface addresses
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Println("Error getting interfaces: ", err)
	} else {
		for i, iface := range ifaces {
			if (iface.Flags & net.FlagUp) == 0 {
				continue
			}

			log.Printf("%02d: %s (%s)", i, iface.Name, iface.Flags) 
			addrs, err := iface.Addrs()
			if err == nil {
				for j, a := range addrs {
					log.Printf("  :%02d %s - %s", j, a.Network(), a)
				}
			}
		}
	}

	// whereami?
	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Could not get current working directory:", err)
		cwd = ""
	}
	
	
	// start serving
	log.Printf("Listening on %s, serving in %s", *listenAddr, cwd)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
