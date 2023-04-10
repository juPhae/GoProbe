//go:build linux
// +build linux

package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type SystemStats struct {
	Device string  `json:"device"`
	CPU    float64 `json:"cpu" unit:"%"`
	Memory float64 `json:"memory" unit:"%"`
	Disk   float64 `json:"disk" unit:"%"`
	NetIn  float64 `json:"net_in" unit:"kbps"`
	NetOut float64 `json:"net_out" unit:"kbps"`
	Time   string  `json:"time"`
}

var deviceName string

func getSystemStatus() string {
	var prevNetIn float64
	var prevNetOut float64
	var jsonStr []byte

	hostInfo, err := host.Info()
	if err != nil {
		log.Fatal(err)
	}
	deviceName = hostInfo.Hostname

	cpuUsage, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Fatal(err)
	}

	memoryStats, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}

	diskStats, err := disk.Usage("/")
	if err != nil {
		log.Fatal(err)
	}

	netStats, err := net.IOCounters(false)
	if err != nil {
		log.Fatal(err)
	}

	// 将字节转换为比特
	netIn := float64(netStats[0].BytesRecv * 8)
	netOut := float64(netStats[0].BytesSent * 8)

	// 计算速度并将单位转换为 kbps
	elapsed := time.Since(time.Now().Add(-time.Second))
	netInSpeed := (netIn - prevNetIn) / elapsed.Seconds() / 1024
	netOutSpeed := (netOut - prevNetOut) / elapsed.Seconds() / 1024

	stats := SystemStats{
		Device: hostInfo.Hostname,
		CPU:    cpuUsage[0],
		Memory: memoryStats.UsedPercent,
		Disk:   diskStats.UsedPercent,
		NetIn:  netInSpeed,
		NetOut: netOutSpeed,
		Time:   time.Now().Format("2006-01-02 15:04:05"),
	}

	s, err := json.MarshalIndent(stats, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	prevNetIn = netIn
	prevNetOut = netOut

	jsonStr = s

	return string(jsonStr)
}

func main() {

	// 设置连接参数
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("mqtt-example-clientsss")

	// 创建客户端实例
	client := mqtt.NewClient(opts)

	// 连接 MQTT 服务器
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// 订阅主题
	if token := client.Subscribe("device/status/topic/management", 0, messageHandler); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	fmt.Println("已订阅主题 test/topic")

	// 发布系统状态信息
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		str := getSystemStatus()
		fmt.Println(str)

		for !client.IsConnected() {
			if token := client.Connect(); token.Wait() && token.Error() != nil {
				log.Println("连接失败，等待重新连接...")
				time.Sleep(1 * time.Second)
				continue
			}
		}
		//发送设备状态信息
		topic := "device/status/topic/info"
		text := str

		if token := client.Publish(topic, 0, false, text); token.Wait() && token.Error() != nil {
			log.Println(token.Error())
			client.Disconnect(250)
			time.Sleep(1 * time.Second)
			continue
		}

		fmt.Printf("已发送消息 %s 到主题 %s\n", text, topic)
	}
}

// 定义设备状态
type DeviceStatus struct {
	Command   string `json:"command"`
	Device    string `json:"device"`
	Parameter string `json:"parameter"`
}

// 处理接收到的消息
func messageHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("收到消息 %s 来自主题 %s\n", string(msg.Payload()), msg.Topic())

	if msg.Topic() == "device/status/topic/management" {
		// 解析 JSON 字符串
		var deviceStatus DeviceStatus
		err := json.Unmarshal(msg.Payload(), &deviceStatus)
		if err != nil {
			fmt.Println("解析 JSON 字符串失败：", err)
			return
		}
		if deviceName == deviceStatus.Device {
			// 执行命令
			switch deviceStatus.Command {
			case "start":
				fmt.Println("正在启动脚本...")
				startShell()
				fmt.Println("脚本启动成功")
			case "stop":
				fmt.Println("正在停止脚本...")
				stopShell()
				fmt.Println("脚本已停止")
			default:
				fmt.Println("未知命令：", deviceStatus.Command)
			}

			// 将执行结果返回给服务器
			result := map[string]string{
				"command": deviceStatus.Command,
				"result":  "success",
			}
			resultJson, err := json.Marshal(result)
			if err != nil {
				fmt.Println("生成 JSON 字符串失败：", err)
				return
			}
			client.Publish("device/status/topic/server", 0, false, resultJson)
		} else {
			fmt.Println("设备名称不匹配")
			return
		}

	}
}
func startShell() error {
	// 检查脚本是否存在
	scriptPath := "/root/keep/keepLive.sh"
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("脚本文件不存在：%s", scriptPath)
	}

	// 创建 cmd 对象
	cmd := exec.Command("/bin/bash", scriptPath)

	// 配置 cmd 对象
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// 启动脚本
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动脚本失败：%s", err)
	}

	return nil
}

func stopShell() error {
	// 查找脚本进程
	cmd := exec.Command("/usr/bin/pgrep", "-f", "/root/keep/keepLive.sh")
	pids, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("查找脚本进程失败：%s", err)
	}

	// 终止脚本进程
	for _, pid := range strings.Split(strings.TrimSpace(string(pids)), "\n") {
		pid = strings.TrimSpace(pid)
		if pid == "" {
			continue
		}
		pidInt, err := strconv.Atoi(pid)
		if err != nil {
			return fmt.Errorf("无效的进程号：%s", pid)
		}
		if err := unix.Kill(-pidInt, unix.SIGTERM); err != nil {
			return fmt.Errorf("终止脚本进程失败：%s", err)
		}
	}

	return nil
}
