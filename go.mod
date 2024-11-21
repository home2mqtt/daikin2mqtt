module github.com/home2mqtt/daikin2mqtt

go 1.21.0

require (
	github.com/home2mqtt/hass v0.0.8
	github.com/samthor/daikin-go v1.0.0
)

require (
	github.com/eclipse/paho.mqtt.golang v1.4.3 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
)

replace github.com/home2mqtt/hass => ../hass
