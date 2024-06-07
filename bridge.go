package daikin2mqtt

import (
	"time"

	"github.com/home2mqtt/hass"
	"github.com/samthor/daikin-go/api"
)

type daikin2MqttBridge struct {
	propertyContext
	mode       stringProperty
	temp       floatProperty
	outside    floatProperty
	targettemp floatProperty
	power      stringProperty
	//humidity   floatProperty
	fan   stringProperty
	swing stringProperty
}

func CreateBridge(runtime hass.IPubSubRuntime, host string) *daikin2MqttBridge {
	bridge := &daikin2MqttBridge{
		propertyContext: propertyContext{
			IPubSubRuntime: runtime,
			base:           "daikin",
			id:             host,
		},
	}
	bridge.mode = bridge.DefineString("mode")
	bridge.temp = bridge.DefineFloat("temp")
	bridge.outside = bridge.DefineFloat("outside")
	bridge.targettemp = bridge.DefineFloat("targettemp")
	bridge.power = bridge.DefineString("power")
	//bridge.humidity = bridge.DefineFloat("humidity")
	bridge.fan = bridge.DefineString("fan")
	bridge.swing = bridge.DefineString("swing")
	return bridge
}

func (bridge *daikin2MqttBridge) HVACDescriptor() *hass.HVAC {
	return &hass.HVAC{
		BasicConfig: hass.BasicConfig{
			UniqueID: "daikin_" + bridge.id + "_climate",
			Device: &hass.Device{
				Name:         "daikin_" + bridge.id + "_climate",
				Manufacturer: "Daikin",
				Model:        "Daikin AC",
				SwVersion:    "0.0.1",
				Identifiers: []string{
					"daikin_ac_" + bridge.id,
				},
			},
		},
		CurrentTemperatureTopic: bridge.temp.StateTopic(),
		TemperatureStateTopic:   bridge.targettemp.StateTopic(),
		TemperatureCommandTopic: bridge.targettemp.CommandTopic(),
		ModeCommandTopic:        bridge.mode.CommandTopic(),
		ModeStateTopic:          bridge.mode.StateTopic(),
		Modes:                   []string{"auto", "off", "cool", "heat", "dry", "fan_only"}, // "auto", "dehum", "cool", "heat", "fan"
		PowerCommandTopic:       bridge.power.CommandTopic(),
		PayloadOn:               "ON",
		PayloadOff:              "OFF",
		//CurrentHumidityTopic:    bridge.humidity.StateTopic(),
		FanModes:              []string{"auto", "quiet", "1", "2", "3", "4", "5"},
		FanModeCommandTopic:   bridge.fan.CommandTopic(),
		FanModeStateTopic:     bridge.fan.StateTopic(),
		SwingModes:            []string{"none", "horizontal", "vertical", "both"},
		SwingModeCommandTopic: bridge.swing.CommandTopic(),
		SwingModeStateTopics:  bridge.swing.StateTopic(),
	}
}

func mode_daikin2mqtt(power bool, daikinmode string) string {
	if !power {
		return "off"
	}
	switch daikinmode {
	case "auto", "cool", "heat":
		return daikinmode
	case "dehum":
		return "dry"
	case "fan":
		return "fan_only"
	}
	return ""
}

func mode_mqtt2daikin(mqttmode string) (bool, string) {
	switch mqttmode {
	case "auto", "cool", "heat":
		return true, mqttmode
	case "dry":
		return true, "dehum"
	case "fan_only":
		return true, "fan"
	case "off":
		return false, ""
	}
	return false, ""
}

func fan_daikin2mqtt(fanrate api.FanRate) string {
	switch fanrate {
	case api.FanRateAuto:
		return "auto"
	case api.FanRateQuiet:
		return "quiet"
	case api.FanRateOne:
		return "1"
	case api.FanRateTwo:
		return "2"
	case api.FanRateThree:
		return "3"
	case api.FanRateFour:
		return "4"
	case api.FanRateFive:
		return "5"
	}
	return "auto"
}

func fan_mqtt2daikin(fanmode string) api.FanRate {
	switch fanmode {
	case "auto":
		return api.FanRateAuto
	case "quiet":
		return api.FanRateQuiet
	case "1":
		return api.FanRateOne
	case "2":
		return api.FanRateTwo
	case "3":
		return api.FanRateThree
	case "4":
		return api.FanRateFour
	case "5":
		return api.FanRateFive
	}
	return api.FanRateUnset
}

func swing_mqtt2daikin(swingmode string) api.FanDir {
	switch swingmode {
	case "none":
		return api.FanDirNone
	case "horizontal":
		return api.FanDirHorizontal
	case "vertical":
		return api.FanDirVertical
	case "both":
		return api.FanDirBoth
	}
	return api.FanDirUnset
}

func swing_daikin2mqtt(fandir api.FanDir) string {
	switch fandir {
	case api.FanDirNone:
		return "none"
	case api.FanDirHorizontal:
		return "horizontal"
	case api.FanDirVertical:
		return "vertical"
	case api.FanDirBoth:
		return "both"
	}
	return "none"
}

func (bridge *daikin2MqttBridge) Attach(ac Daikin, nodeid string) {
	hvac := bridge.HVACDescriptor()
	hass.AnnounceDevice(bridge, "homeassistant", nodeid, hvac.GetBasic().UniqueID, hvac)

	bridge.mode.OnCommand(func(value string) {
		ac.GetAndSet(func(ds *DaikinState) bool {
			power, mode := mode_mqtt2daikin(value)
			if mode != "" {
				ds.Mode = (value)
			}
			ds.Power = power
			return true
		})
	})
	bridge.targettemp.OnCommand(func(value float64) {
		ac.GetAndSet(func(ds *DaikinState) bool {
			ds.Temp = value
			return true
		})
	})
	bridge.power.OnCommand(func(value string) {
		if value == "ON" {
			ac.GetAndSet(func(ds *DaikinState) bool {
				ds.Power = true
				return true
			})
		}
		if value == "OFF" {
			ac.GetAndSet(func(ds *DaikinState) bool {
				ds.Power = false
				return true
			})
		}
	})
	bridge.fan.OnCommand(func(value string) {
		ac.GetAndSet(func(ds *DaikinState) bool {
			ds.FanRate = fan_mqtt2daikin(value)
			return ds.FanRate != api.FanRateUnset
		})
	})
	bridge.swing.OnCommand(func(value string) {
		ac.GetAndSet(func(ds *DaikinState) bool {
			ds.FanDir = swing_mqtt2daikin(value)
			return ds.FanDir != api.FanDirUnset
		})
	})

	go func() {
		for range time.Tick(time.Minute) {
			sensor, err := ac.ReadSensor()
			if err == nil {
				bridge.temp.NotifyState(sensor.Temp)
				bridge.outside.NotifyState(*sensor.World)
			}
			ci := ac.State()
			bridge.targettemp.NotifyState(ci.Temp)
			bridge.mode.NotifyState(mode_daikin2mqtt(ci.Power, ci.Mode))
			bridge.fan.NotifyState(fan_daikin2mqtt(ci.FanRate))
			bridge.swing.NotifyState(swing_daikin2mqtt(ci.FanDir))
		}
	}()
}
