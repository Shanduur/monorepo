package main

import (
	"net"
	"os"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"my.mod/cproc"
)

func main() {
	cp := cproc.NewProcessor(os.Stdout, os.Stderr, os.Stdin, "> ")
	sp, err := cproc.NewSafeProcessor(cp, "ping.bin")
	if err != nil {
		cp.Log.Fatalf("failed to create safe processor: %s", err.Error())
	}
	defer sp.Close()

	cp.AddCommand("ping", func(s ...string) {
		for _, sub := range s {
			ping(*cp, sub)
		}
	}, "", "sends ping request to the given host (can be IP or FQDN), IPv4 only", "long")

	cp.ExitHandler = func() error {
		return nil
	}

	sp.Run()
}

func ping(cp cproc.CommandProcessor, target string) {
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		cp.Logf("listen err, %s", err)
		return
	}
	defer c.Close()

	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("HELLO-R-U-THERE"),
		},
	}
	wb, err := wm.Marshal(nil)
	if err != nil {
		cp.Logf("failed to marshal: %s", err.Error())
		return
	}

	addr, err := net.LookupIP(target)
	if err != nil {
		cp.Logf("unable to lookup ip: %s", err.Error())
		return
	}

	for _, n := range addr {

		if _, err := c.WriteTo(wb, &net.IPAddr{IP: n}); err != nil {
			cp.Logf("write to error: %s", err.Error())
			return
		}

		rb := make([]byte, 1500)
		n, peer, err := c.ReadFrom(rb)
		if err != nil {
			cp.Logf("read from error: %s", err.Error())
			return
		}
		rm, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), rb[:n])
		if err != nil {
			cp.Logf("read from error: %s", err.Error())
			return
		}
		switch rm.Type {
		case ipv4.ICMPTypeEchoReply:
			cp.Printf("got reflection from %v (target: %s)\n", peer, target)
		default:
			cp.Logf("got %+v; want echo reply (target: %s)", rm, target)
		}
	}
}
