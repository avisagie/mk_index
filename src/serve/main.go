
package main

import (
	"log"
	"net/http"
	"flag"
	"net"
	"os"
	"fmt"
)

var (
	listenAddr = flag.String("addr", ":8080", "Address to listen on")
	cacheMaxAge = flag.Int("max-age", 60, "Seconds to allow caching of resources on the client side")
)

func init() {
	flag.Parse()
}

func printAddresses(port string) {
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
				for _, a := range addrs {
					var ip net.IP
					switch v := a.(type) {
					case *net.IPNet:
						ip = v.IP
					case *net.IPAddr:
						ip = v.IP
					}
					log.Printf("      http://%s/", net.JoinHostPort(ip.String(), port))
				}
			}
		}
	}
}

func main() {
	maxAge := fmt.Sprintf("max-age=%d", *cacheMaxAge)

	h := http.FileServer(http.Dir("./"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s -> %s", r.RemoteAddr, r.URL.String())
		w.Header().Set("Cache-Control", maxAge)
		h.ServeHTTP(w, r)
	})

	listenHost, listenPort, err := net.SplitHostPort(*listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	
	if len(listenHost) == 0 {
		printAddresses(listenPort)
	} else {
		log.Printf("Listening on http://%s/", *listenAddr)
	}

	// whereami?
	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Could not get current working directory:", err)
		cwd = ""
	}
	
	// start serving
	log.Printf("Serving in %s", cwd)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
