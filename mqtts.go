package mqtts

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sacca97/unitn-mqtts/crypto"
)

type mqttConfig struct {
	Brokers              []string
	ClientID             string
	Username             string
	Password             string
	Topics               []string
	QoS                  int
	Retained             bool
	AutoReconnect        bool
	MaxReconnectInterval time.Duration
	PersistentSession    bool
	KeepAlive            uint16
	TLSCA                string
	TLSCert              string
	TLSKey               string
	Version              int
}

func Init() {
	cfg := mqtt.NewClientOptions()
	cfg.AddBroker("tcp://broker.emqx.io:1883")
	cfg.SetClientID("go_mqtt_client")
	cfg.SetUsername("emqx")
	cfg.SetPassword("public")
	cfg.SetDefaultPublishHandler(messageHandler)
	cfg.OnConnect = connectionHandler
	cfg.OnConnectionLost = connectionLostHandler
	client := mqtt.NewClient(cfg)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	Subscribe(client, "topic/test", 0)
	Publish(client, "topic/test", []byte("Hello World!"))

}

func HandleEncrypted(msg mqtt.Message) {
	ciphertext := msg.Payload()
	plaintext, err := crypto.NewCPABE(false, "").DecryptDecode(ciphertext)
	if err != nil {
		fmt.Println(plaintext)
	}
	mqtt.NewClientOptions()
	//Do something with the plaintext
}

func Subscribe(client mqtt.Client, topic string, qos byte) {
	token := client.Subscribe(topic, qos, nil)
	token.Wait()
}

func Publish(client mqtt.Client, topic string, payload []byte) error {
	token := client.Publish(topic, 0, false, payload)
	token.Wait()
	return nil
}

func isEncrypted(msg mqtt.Message) bool {
	//TODO: Check if the message is encrypted
	return false
}

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	if isEncrypted(msg) {
		HandleEncrypted(msg)
	} else {
		fmt.Println(msg.Topic(), string(msg.Payload()))
	}
}

var connectionHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Println("Connection lost")
}
