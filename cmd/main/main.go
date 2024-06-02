package main

import (
	"flag"
	"log"

	"github.com/balazsgrill/hass/paho"
	"github.com/myhomemqtt/daikin2mqtt"
)

func main() {
	host := flag.String("h", "", "Host or IP address of device")
	broker := flag.String("b", "tcp://192.168.0.1", "URL of MQTT broker")
	flag.Parse()

	if *host == "" {
		log.Panic("Host is required")
	}

	log.Printf("connecting to %s\n", *host)
	runtime := paho.NewByURL(*broker)

	daikin2mqtt.CreateBridge(runtime, *host).Attach(daikin2mqtt.New(*host), "f2rpi40.lan")

	select {}
}
