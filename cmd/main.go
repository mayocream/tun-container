package main

import (
	"flag"

	"github.com/txthinking/socks5"
	tunhijack "github.com/mayocream/tun-hijack"
)

var config = new(tunhijack.Config)

func init() {
	flag.StringVar(&config.TunName, "tun", "hijack0", "TUN device name")
	flag.Parse()
}

func main() {
	go socks5Server()
	tunhijack.Run(config)
}

func socks5Server() error {
	s, _ := socks5.NewClassicServer("127.0.0.1:1234", "127.0.0.1", "", "", 5, 5)
	return s.ListenAndServe(nil)
}
