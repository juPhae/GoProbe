package mqtt

import (
	"container/list"
	"fmt"
	mqttC "github.com/eclipse/paho.mqtt.golang"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/hooks/auth"
	"github.com/mochi-co/mqtt/v2/listeners"
	"log"
)

const (
	authUsername = "admin"
	authPassword = "password"
)

func MqttServerStart() {
	// 创建一个新的 MQTT 服务器。
	server := mqtt.New(nil)

	// 允许所有连接。
	_ = server.AddHook(new(auth.AllowHook), nil)

	// 创建一个标准端口的 TCP 监听器。
	tcp := listeners.NewTCP("t1", ":1883", nil)
	err := server.AddListener(tcp)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Serve()
	if err != nil {
		log.Fatal(err)
	}
}

func MqttClientStart() {
	// 设置连接参数
	opts := mqttC.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("mqtt-example-client")

	// 创建客户端实例
	client := mqttC.NewClient(opts)

	// 连接 MQTT 服务器
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	// 订阅泛主题
	if token := client.Subscribe("device/status/topic/#", 0, messageHandler); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}

// Queue 定义一个最大长度为20的队列
var Queue = list.New()

const queueMaxLen = 24

// 处理接收到的消息
func messageHandler(client mqttC.Client, msg mqttC.Message) {
	fmt.Printf("收到消息 %s 来自主题 %s\n", string(msg.Payload()), msg.Topic())
	if msg.Topic() == "device/status/topic/info" {
		//存到队列中
		// 将消息添加到队列中
		if Queue.Len() >= queueMaxLen {
			// 队列已满，移除队头元素
			Queue.Remove(Queue.Front())
		}
		Queue.PushBack(msg.Payload())
	}
	if msg.Topic() == "device/status/topic/management" {

	}
}

func publish(client mqttC.Client) {

	// 发送消息
	topic := "test/topic"
	text := "PUSH, MQTT!"
	if token := client.Publish(topic, 0, false, text); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	fmt.Printf("已发送消息 %s 到主题 %s\n", text, topic)
}
