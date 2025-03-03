module github.com/home2mqtt/daikin2mqtt

go 1.21.0

require (
	github.com/home2mqtt/hass v0.1.0
	github.com/samthor/daikin-go v1.0.0
)

require (
	github.com/eclipse/paho.mqtt.golang v1.5.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
)

replace github.com/home2mqtt/hass => ../hass
