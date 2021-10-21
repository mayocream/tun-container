package main

import (
	"log"

	tun "github.com/mayocream/tun-container"
)

func main() {
	t, err := tun.NewTun("tunc0")
	if err != nil {
		log.Fatal(err)
	}

	defer t.Close()

	log.Println("ip addr add 192.168.10.0/24 dev tunc0")
	if err := t.Up("192.168.10.0/24"); err != nil {
		log.Fatal(err)
	}

	log.Println("accepting packets")

	t.HandlePackets()
}
