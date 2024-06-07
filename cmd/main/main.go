package main

import (
	"flag"
	"log"

	"github.com/home2mqtt/daikin2mqtt"
	"github.com/home2mqtt/hass/paho"
)

func main() {
	host := flag.String("h", "", "Host or IP address of device")
	broker := flag.String("b", "tcp://192.168.0.1:1883", "URL of MQTT broker")
	flag.Parse()

	if *host == "" {
		log.Panic("Host is required")
	}

	log.Printf("connecting to %s\n", *host)
	runtime := paho.NewByURL(*broker)
	err := runtime.Connect()
	if err != nil {
		log.Panic(err)
	}

	daikin2mqtt.CreateBridge(runtime, *host).Attach(daikin2mqtt.New(*host), "f2rpi40.lan")

	select {}
}
