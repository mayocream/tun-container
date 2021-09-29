package tunhijack

import (
	"log"
	"testing"

	"tailscale.com/net/tstun"
)

func TestTun(t *testing.T) {
	dev, devName, err := tstun.New(log.Printf, "tun0")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(devName)
	defer dev.Close()
}
