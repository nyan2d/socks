package main

import (
	"flag"
)

func main() {
	// parsing cmd args
	host := flag.String("host", ":1080", "hostname")
	allowGuest := flag.Bool("guest", true, "Allows join to the server with no authentication.")
	authFile := flag.String("auth", "", "A JSON file with a user accounts.")
	dnsFile := flag.String("dns", "", "A JSON file with a DNS records.")
	flag.Parse()

	// reading configs
	dns := NewDNS()
	if len(*dnsFile) > 0 {
		if err := dns.ReadFromJsonFile(*dnsFile); err != nil {
			panic(err)
		}
	}
	auth := NewAuth(*allowGuest)
	if len(*authFile) > 0 {
		if err := auth.ReadFromJsonFile(*authFile); err != nil {
			panic(err)
		}
	}

	// starting socks server
	server := NewSocksServer(dns, auth)
	err := server.Listen(*host)
	if err != nil {
		panic(err)
	}
}
