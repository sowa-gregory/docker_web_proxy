package main

import (
	"log"
	"net"
	"strconv"
	"fmt"
	"os"
	"github.com/miekg/dns"
)

var ip string
const domain = "home." // domain name to response for
const port = 53
const version = "HomeDNSServer 0.9"

func getDockerHostIP(dns string) string {
	ips, err := net.LookupIP(dns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		os.Exit(1)
	}
	return ips[0].String()
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		m.Authoritative = true
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s, response:%s\n", q.Name, ip)

				m.Answer = append(m.Answer, &dns.A{
				Hdr: dns.RR_Header{ Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60 },
				A: net.ParseIP(ip),
			})

		}
	}
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false
	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}
	w.WriteMsg(m)
}

func main() {
	log.Printf(version+"\n")
	proxy_host_dns := os.Getenv("PROXY_HOST")
	if (proxy_host_dns == "" ) {
		log.Printf("PROXY_HOST env not set\n")
		os.Exit(1)
}
	log.Printf("Proxy DNS:%s\n", proxy_host_dns)
	ip = getDockerHostIP(proxy_host_dns)
	log.Printf("DNS host ip: %s\n", ip)
	// attach request handler func
	dns.HandleFunc(domain, handleDNSRequest)

	// start server
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at UDP port:%d\n", port)
	log.Printf("Resolves DNS for domain:%s\n", domain)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
	log.Printf("running\n")
}
