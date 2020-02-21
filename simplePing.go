package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	dest := "1.1.1.1"

	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	transTime := strconv.FormatInt(time.Now().UnixNano(), 10)
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte(transTime),
		},
	}
	wb, err := wm.Marshal(nil)

	if err != nil {
		log.Fatal(err)
	}
	if _, err := c.WriteTo(wb, &net.IPAddr{IP: net.ParseIP(dest)}); err != nil {
		log.Fatal(err)
	}

	rb := make([]byte, 1500)
	n, _, err := c.ReadFrom(rb)
	if err != nil {
		log.Fatal(err)
	}
	rm, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), rb[:n])
	rcvTime := time.Now().UnixNano()
	if err != nil {
		log.Fatal(err)
	}
	if rm.Type == ipv4.ICMPTypeEchoReply {
		fromTime, _ := strconv.ParseInt(string(rm.Body.(*icmp.Echo).Data), 10, 64)
		rtt := rcvTime - fromTime
		fmt.Println("OK RTT = ", rtt/1e6, "ms")
	} else {
		fmt.Printf("NG")
	}

}
