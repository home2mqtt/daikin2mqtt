package daikin2mqtt

import (
	"strconv"

	"github.com/balazsgrill/hass"
)

func HVACDescriptor(id string) *hass.HVAC {
	return &hass.HVAC{
		BasicConfig: hass.BasicConfig{
			UniqueID: "daikin_" + id + "_climate",
			Device: &hass.Device{
				Name:         "daikin_" + id + "_climate",
				Manufacturer: "Daikin",
				Model:        "Daikin AC",
				SwVersion:    "0.0.1",
				Identifiers: []string{
					"daikin_ac_" + id,
				},
			},
		},
		CurrentTemperatureTopic: "daikin/" + id + "/temp",
		TemperatureStateTopic:   "daikin/" + id + "/targettemp",
		TemperatreCommandTopic:  "daikin/" + id + "/targettemp/set",
		ModeCommandTopic:        "daikin/" + id + "/mode/set",
	}
}

func Attach(runtime hass.IPubSubRuntime, ac Daikin, config *hass.HVAC) {
	runtime.Receive(config.TemperatreCommandTopic, func(topic string, payload []byte) {
		value, err := strconv.ParseFloat(string(payload), 64)
		if err != nil {
			return
		}
		ac.GetAndSet(func(ds *DaikinState) bool {
			ds.Temp = value
			return true
		})
	})
}
