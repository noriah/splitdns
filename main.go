package main

import (
	"encoding/json"
	"log"
	"net"
	"os"

	"github.com/miekg/dns"
)

type Config struct {
	ListenAddress string `json:"listen"`
	Zones         []Zone `json:"zones"`
}

type Zone struct {
	Name    string   `json:"name"`
	Servers []string `json:"servers"`
}

// returns a handler that forwards requests to the servers given, trying each
// next server only after failure of the previous. on failure of all servers,
// replies to request with dns.RcodeServerFailure
func nameHandler(servers []string) dns.HandlerFunc {
	return func(w dns.ResponseWriter, r *dns.Msg) {
		reply := dns.Msg{}
		reply.SetReply(r)

		for _, s := range servers {
			res, err := dns.Exchange(r, s)

			if err == nil {
				res.CopyTo(&reply)
				w.WriteMsg(&reply)
				return
			}

			log.Printf("got err attempting query %s: %s\n", s, err)
		}

		reply.SetRcode(r, dns.RcodeServerFailure)
		w.WriteMsg(&reply)
	}
}

// registers handlers for all Zones in config, then starts a dns server on the
// ListenAddress in config
func dnsServer(config *Config) error {
	for _, e := range config.Zones {
		dns.HandleFunc(e.Name, nameHandler(e.Servers))
	}

	addr, err := net.ResolveUDPAddr("udp", config.ListenAddress)
	if err != nil {
		return err
	}

	socket, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	return dns.ActivateAndServe(nil, socket, dns.DefaultServeMux)
}

// reads the configuration from `path`
func readConfig(path string) (config Config, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &config)

	return
}

// runs splitdns
func main() {
	if len(os.Args) < 2 {
		log.Fatalf("path to config required\n")
	}

	config, err := readConfig(os.Args[1])

	if err != nil {
		log.Fatalf("config read failed: %s\n", err)
	}

	if err := dnsServer(&config); err != nil {
		log.Fatalf("dns server quit with non-zero: %s\n", err)
	}
}
