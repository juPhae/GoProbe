package main

import (
	"goprobe/internal/mqtt"
	"goprobe/internal/service"
)

func main() {
	go mqtt.MqttClientStart()
	go mqtt.MqttServerStart()
	service.GinStart()
	//time.Sleep(1000 * time.Second)
}
