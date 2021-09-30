package tunhijack

import (
	"fmt"
	"log"

	"tailscale.com/net/tstun"
)

// Config config
type Config struct {
	TunName string
}

// Run runs hijack mode
func Run(config *Config) {
	dev, devName, err := tstun.New(log.Printf, config.TunName)
	if err != nil {
		tstun.Diagnose(log.Printf, devName)
		log.Fatalf("TUN device create err: %s", err)
	}

	log.Printf("Create TUN device %s", devName)

	wrap := tstun.Wrap(log.Printf, dev)
	defer wrap.Close()

	for {
		buf := make([]byte, 1024)
		n, err := wrap.Read(buf, 0)
		if err != nil {
			log.Fatal(err)
		}
		if n > 0 {
			data := buf[:n]
			fmt.Println(string(data))
		}
	}
}
