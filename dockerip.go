// Package example is a CoreDNS plugin that prints "example" to stdout on every packet received.
//
// It serves as an example CoreDNS plugin with numerous code comments.
package dockerip

import (
	"context"
	"net"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

// Define log to be a logger with the plugin name in it. This way we can just use log.Info and
// friends to log.
var log = clog.NewWithPlugin("dockerip")

// Example is an example plugin to show how to write a plugin.
type Dockerip struct {
	Next   plugin.Handler
	Target string
}

// ServeDNS implements the plugin.Handler interface. This method gets called when example is used
// in a Server.
func (d Dockerip) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	// regex := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	state := request.Request{W: w, Req: r}
	// fmt.Printf("QType of %v", state.QType())
	// if(state.QType() != dns.TypeA) {
	// 	var err dns.
	// 	return plugin.Error(d.Name(), dns.ErrRcode)
	// 	//return plugin.NextOrFailure(d.Name(), d.Next, ctx, w, r)
	// }

	lookup, err := net.LookupIP(d.Target)
	if err != nil {
		log.Debug("Error")
		log.Debug(err.Error())
		return plugin.NextOrFailure(d.Name(), d.Next, ctx, w, r)
	}
	ip := lookup[0]

	// log.Debug(ip.String())

	a := new(dns.Msg)
	a.SetReply(r)
	if ip != nil && state.QType() == dns.TypeA {
		// log.Debug(ip)
		a.Authoritative = true

		//var rr dns.RR
		rr := new(dns.A)
		rr.Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass()}
		rr.A = ip

		// srv := new(dns.SRV)
		// srv.Hdr = dns.RR_Header{Name: "_" + state.Proto() + "." + state.QName(), Rrtype: dns.TypeSRV, Class: state.QClass()}
		// if state.QName() == "." {
		// 	srv.Hdr.Name = "_" + state.Proto() + state.QName()
		// }
		// port, _ := strconv.ParseUint(state.Port(), 10, 16)
		// srv.Port = uint16(port)
		// srv.Target = "."

		a.Answer = []dns.RR{rr}
		// a.Extra = []dns.RR{srv}

		return dns.RcodeSuccess, w.WriteMsg(a)
	}

	a.Rcode = dns.RcodeSuccess
	a.Opcode = dns.OpcodeQuery

	return 0, w.WriteMsg(a)
}

// Name implements the Handler interface.
func (e Dockerip) Name() string { return "dockerip" }
