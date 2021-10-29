package tun

import (
	"bytes"
	"fmt"
	"log"
	"sync"

	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
	wtun "golang.zx2c4.com/wireguard/tun"
	"tailscale.com/net/packet"
	"tailscale.com/net/tstun"
	"tailscale.com/types/ipproto"
)

// parsedPacketPool holds a pool of Parsed structs for use in filtering.
// This is needed because escape analysis cannot see that parsed packets
// do not escape through {Pre,Post}Filter{In,Out}.
var parsedPacketPool = sync.Pool{New: func() interface{} { return new(packet.Parsed) }}

var bufferPool = sync.Pool{New: func() interface{} {
	buf := make([]byte, tstun.MaxPacketSize)
	return bytes.NewBuffer(buf)
}}

var offset = tstun.PacketStartOffset

// Tun device
type Tun struct {
	dev     wtun.Device
	devName string
}

func NewTun(tunName string) (*Tun, error) {
	dev, devName, err := tstun.New(log.Printf, tunName)
	if err != nil {
		return nil, err
	}

	log.Printf("create tun device: %s\n", devName)

	t := &Tun{
		dev:     dev,
		devName: devName,
	}

	return t, nil
}

func (t *Tun) Up(addr string) error {
	tunNet, err := netlink.LinkByName(t.devName)
	if err != nil {
		return errors.Wrap(err, "ip link")
	}
	if err := netlink.LinkSetUp(tunNet); err != nil {
		return errors.Wrap(err, "ip link up")
	}
	naddr, err := netlink.ParseAddr(addr)
	if err != nil {
		return errors.Wrap(err, "parse addr")
	}
	if err := netlink.AddrAdd(tunNet, naddr); err != nil {
		return errors.Wrap(err, "ip addr add")
	}

	return nil
}

// HandlePackets ...
func (t *Tun) HandlePackets() error {
	for {
		buf := bufferPool.Get().(*bytes.Buffer)
		defer bufferPool.Put(buf)

		n, err := t.dev.Read(buf.Bytes(), offset)
		if err != nil {
			log.Printf("[E] read iface fail: %v\n", err)
			break
		}

		p := parsedPacketPool.Get().(*packet.Parsed)
		defer parsedPacketPool.Put(p)
		p.Decode(buf.Bytes()[offset : offset+n])

		log.Println(p.String())
		t.respondToICMP(p)
		t.respondToDNS(p)
	}

	return nil
}

func (t *Tun) Close() {
	t.dev.Close()
}

func (t *Tun) respondToICMP(p *packet.Parsed) {
	if p.IsEchoRequest() {
		log.Println("response to ping")
		header := p.ICMP4Header()
		header.ToResponse()
		outp := packet.Generate(&header, p.Payload())

		buf := bufferPool.Get().(*bytes.Buffer)
		defer bufferPool.Put(buf)

		copy(buf.Bytes()[offset:], outp)

		defer t.dev.Flush()
		if _, err := t.dev.Write(buf.Bytes()[:offset+len(outp)], offset); err != nil {
			fmt.Println("[E]: write, ", err)
		}
	}
}

func (t *Tun) respondToICMP2(p *packet.Parsed) {
	if p.IsEchoRequest() {
		log.Println("response to ping")

		buf := bufferPool.Get().(*bytes.Buffer)
		defer bufferPool.Put(buf)

		copy(buf.Bytes()[offset:], p.Buffer())

		defer t.dev.Flush()
		if _, err := t.dev.Write(buf.Bytes()[:offset+len(p.Buffer())], offset); err != nil {
			fmt.Println("[E]: write, ", err)
		}
	}
}

func (t *Tun) respondToDNS(p *packet.Parsed) {
	// DNS query
	if p.Dst.Port() == 53 && p.IPProto == ipproto.UDP {
		log.Println("forward dns")

		buf := bufferPool.Get().(*bytes.Buffer)
		defer bufferPool.Put(buf)

		copy(buf.Bytes()[offset:], p.Buffer())

		defer t.dev.Flush()
		if _, err := t.dev.Write(buf.Bytes()[:offset+len(p.Buffer())], offset); err != nil {
			fmt.Println("[E]: write, ", err)
		}
	}
}
