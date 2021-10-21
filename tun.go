package tun

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
	wtun "golang.zx2c4.com/wireguard/tun"
	"tailscale.com/net/packet"
	"tailscale.com/net/tstun"
	"tailscale.com/wgengine/filter"
)

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

func (t *Tun) Read() ([]byte, error) {
	buf := make([]byte, 2048)
	n, err := t.dev.Read(buf, tstun.PacketStartOffset)
	if err != nil {
		return nil, err
	}

	return buf[tstun.PacketStartOffset : tstun.PacketStartOffset+n], nil
}

func (t *Tun) Close() {
	t.dev.Close()
}

// HandlePackets ...
func (t *Tun) HandlePackets() error {
	for {
		buf, err := t.Read()
		if err != nil {
			log.Printf("[E] read iface fail: %v\n", err)
			break
		}

		p := new(packet.Parsed)
		p.Decode(buf)
		fmt.Println(p.String())
		t.respondToICMP(p)
	}

	return nil
}

func (t *Tun) Write(buf []byte) (int, error){
	return t.dev.Write(buf, tstun.PacketStartOffset)
}

func (t *Tun) respondToICMP(p *packet.Parsed)  {
	if p.IsEchoRequest() {
		fmt.Println("response to ping")
		header := p.ICMP4Header()
		header.ToResponse()
		outp := packet.Generate(&header, p.Payload())
		fmt.Println("output: ", outp)
		fmt.Println("output len: ", len(outp))


		buf := make([]byte, 2048)
		copy(buf[tstun.PacketStartOffset:], outp)

		n, err := t.dev.Write(buf[:tstun.PacketStartOffset+len(outp)], tstun.PacketStartOffset)
		if err != nil {
			fmt.Println("[E]: write, ", err)
		}
		fmt.Println("write len: ", n)
		t.dev.Flush()
	}
}

// handleLocalPackets inspects packets coming from the local network
// stack, and intercepts any packets that should be handled by
// tailscaled directly. Other packets are allowed to proceed into the
// main ACL filter.
func handleLocalPackets(p *packet.Parsed, t *tstun.Wrapper) filter.Response {

	return filter.Accept
}
