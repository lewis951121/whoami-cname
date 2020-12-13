// Package whoami implements a plugin that returns details about the resolving
// querying it.
package whoami_cname

import (
	"context"
	"net"
	"strconv"

	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

const name = "whoami_cname"

// Whoami_Cname is a plugin that returns your IP address, port and the protocol used for connecting
// to CoreDNS. Different from the built-in whoami, this plugin includes the above info in a CNAME rdata.
type Whoami_Cname struct{}

// ServeDNS implements the plugin.Handler interface.
func (wh Whoami_Cname) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	a := new(dns.Msg)
	a.SetReply(r)
	a.Authoritative = true

	ip := state.IP()
	var rr dns.RR

	// the address of the egress resolver
	switch state.Family() {
	case 1:
		rr = new(dns.A)
		rr.(*dns.A).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass()}
		rr.(*dns.A).A = net.ParseIP(ip).To4()
	case 2:
		rr = new(dns.AAAA)
		rr.(*dns.AAAA).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeAAAA, Class: state.QClass()}
		rr.(*dns.AAAA).AAAA = net.ParseIP(ip)
	}

	// todo: returning the src port, 0x20 status, ...
	/*
	srv := new(dns.SRV)
	srv.Hdr = dns.RR_Header{Name: "_" + state.Proto() + "." + state.QName(), Rrtype: dns.TypeSRV, Class: state.QClass()}
	if state.QName() == "." {
		srv.Hdr.Name = "_" + state.Proto() + state.QName()
	}
	port, _ := strconv.Atoi(state.Port())
	srv.Port = uint16(port)
	srv.Target = "."
	*/

	a.Answer = []dns.RR{rr}

	w.WriteMsg(a)

	return 0, nil
}

// Name implements the Handler interface.
func (wh Whoami_Cname) Name() string { return name }
