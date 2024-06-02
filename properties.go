package daikin2mqtt

import (
	"log"
	"strings"

	"github.com/balazsgrill/hass"
)

type propertyContext struct {
	hass.IPubSubRuntime
	base string
	id   string
}

func (pc *propertyContext) DefineString(name string) stringProperty {
	return stringProperty{
		property: property{
			propertyContext: pc,
			name:            name,
		},
	}
}

func (pc *propertyContext) DefineFloat(name string) floatProperty {
	return floatProperty{
		property: property{
			propertyContext: pc,
			name:            name,
		},
	}
}

type property struct {
	*propertyContext
	name string
}

func (p *property) StateTopic() string {
	return strings.Join([]string{p.base, p.id, p.name}, "/")
}

func (p *property) CommandTopic() string {
	return strings.Join([]string{p.base, p.id, p.name, "set"}, "/")
}

type stringProperty struct {
	property
}

func (p *stringProperty) NotifyState(value string) {
	hass.SendString(p, p.StateTopic(), value)
}

func (p *stringProperty) OnCommand(callback func(value string)) {
	hass.ReceiveString(p, p.CommandTopic(), func(topic, payload string) {
		callback(payload)
	})
}

type floatProperty struct {
	property
}

func (p *floatProperty) NotifyState(value float64) {
	hass.SendFloat(p, p.StateTopic(), value)
}

func (p *floatProperty) OnCommand(callback func(value float64)) {
	hass.ReceiveFloat(p, p.CommandTopic(), func(topic string, payload float64, err error) {
		if err == nil {
			callback(payload)
		} else {
			log.Printf("Float value error received on %s: %v\n", topic, err)
		}
	})
}
