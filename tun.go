package tunhijack

import (
	"fmt"
	"io"
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
	defer dev.Close()

	log.Printf("Create TUN device %s", devName)

	for {
		buf := make([]byte, 1024)
		n, err := dev.Read(buf, 4)
		if err == io.EOF {
			if n > 0 {
				data := buf[:n]
				fmt.Println(string(data))
			}
			continue
		}
	}
}
