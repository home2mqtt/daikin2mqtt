package daikin2mqtt

import (
	"fmt"
	"log"
	"time"

	"github.com/home2mqtt/hass"
	"github.com/home2mqtt/hass/bridge"
)

type powerSensor struct {
	ac         Daikin
	descriptor *hass.Sensor
	runtime    hass.IPubSubRuntime
}

func PowerSensorDescriptor(id string) *hass.Sensor {
	return &hass.Sensor{
		BasicConfig: hass.BasicConfig{
			UniqueID: "daikin_" + id + "_energy",
			Device: &hass.Device{
				Name:         "daikin_" + id + "_energy",
				Manufacturer: "Home",
				Model:        "Virtual power sensor",
				SwVersion:    "0.0.1",
				Identifiers: []string{
					"daikin_ac_virtual_power_sensor_" + id,
				},
			},
		},
		UnitOfMeasurement: "kWh",
		Name:              id + "_energy",
		Topic:             "daikin/" + id + "/energy",
		StateClass:        "total_increasing",
		DeviceClass:       "energy",
		Icon:              "mdi:energy",
	}
}

func (s *powerSensor) init(bridgehost string) {
	log.Println("Announcing " + s.descriptor.UniqueID)
	bridge.AnnounceDevice(s.runtime, "homeassistant", bridgehost, s.descriptor.UniqueID, s.descriptor)
	go s.tick()
}

func (s *powerSensor) tick() {
	p, err := s.ac.GetMonthPowerEx()
	if err != nil {
		log.Printf("Couldn't read energy consumption: %v", err)
		return
	}
	energy := float64(p.CurrCool.Sum()+p.CurrHeat.Sum()) * 0.1
	s.runtime.Send(s.descriptor.Topic, []byte(fmt.Sprintf("%f", energy)))
}

func ProvideEnergy(runtime hass.IPubSubRuntime, ac Daikin, id string, bridgehost string) {
	item := &powerSensor{
		descriptor: PowerSensorDescriptor(id),
		ac:         ac,
		runtime:    runtime,
	}
	item.init(bridgehost)

	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for range ticker.C {
			go item.tick()
		}
	}()
}
