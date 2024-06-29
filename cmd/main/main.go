package main

import (
	"flag"
	"log"
	"os"

	"github.com/home2mqtt/daikin2mqtt"
	"github.com/home2mqtt/hass/paho"
)

func main() {
	host := flag.String("h", "", "Host or IP address of device")
	broker := flag.String("b", "tcp://192.168.0.1:1883", "URL of MQTT broker")
	bridgehost := flag.String("bh", "", "Host or IP address of bridge device")
	flag.Parse()

	if *host == "" {
		log.Panic("Host is required")
	}
	if *bridgehost == "" {
		var err error
		*bridgehost, err = os.Hostname()
		if err != nil {
			log.Panic(err)
		}
	}

	log.Printf("connecting to %s\n", *host)
	runtime := paho.NewByURL(*broker)
	err := runtime.Connect()
	if err != nil {
		log.Panic(err)
	}

	daikindevice := daikin2mqtt.New(*host)
	daikin2mqtt.CreateBridge(runtime, *host).Attach(daikindevice, *bridgehost)
	daikin2mqtt.ProvideEnergy(runtime, daikindevice, *host, *bridgehost)

	select {}
}
