package main

import (
	"github.com/shirou/gopsutil/host"
	"goprobe/internal/mqtt"
	"goprobe/internal/service"
	"goprobe/internal/util"
)

func main() {
	//获取设备名称和cpuID
	hostInfo, _ := host.Info()
	util.DeviceName = hostInfo.Hostname
	util.CpuID = hostInfo.HostID
	//启动mqtt客户端和服务端
	go mqtt.MqttClientStart()
	go mqtt.MqttServerStart()
	//启动gin服务
	service.GinStart()
}
