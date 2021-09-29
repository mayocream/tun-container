package main

import (
	"flag"

	tunhijack "github.com/mayocream/tun-hijack"
)

var config = new(tunhijack.Config)

func init() {
	flag.StringVar(&config.TunName, "tun", "hijack0", "TUN device name")
	flag.Parse()
}

func main() {
	tunhijack.Run(config)
}
