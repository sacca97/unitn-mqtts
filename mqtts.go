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

type MqttClient struct {
	client           mqtt.Client
	encryptionScheme *crypto.Cpabe
}

func NewClient(broker string) *MqttClient {
	c := &MqttClient{}
	cfg := mqtt.NewClientOptions()
	cfg.AddBroker(broker)
	cfg.SetClientID("test")

	cfg.SetDefaultPublishHandler(messageHandler)
	cfg.OnConnect = connectionHandler
	cfg.OnConnectionLost = connectionLostHandler
	c.client = mqtt.NewClient(cfg)
	c.encryptionScheme = crypto.NewCPABE()
	return c
}

func Connect(client *MqttClient) {
	if token := client.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func Init(s crypto.Cpabe) {

	public, secret, _ := s.PubKeygen()
	attribute, _ := s.PrivKeygen(secret, []string{"0", "1", "2", "3", "5"})
	s.SetAttribKey(attribute)
	s.SetPublicKey(public)

	/*Subscribe(client, "simola/test", 0)
	msg := "This is a test message"
	policy := "((0 AND 1) OR (2 AND 3)) AND 5"
	payload, _ := s.EncryptEncode(policy, msg)
	fmt.Println(len(payload))
	Publish(client, "simola/test", payload)*/
}

func HandleEncrypted(c crypto.Cpabe, payload []byte) {

	//If key not loaded, load it

	//Need to find a way to save the state and the keys

	plaintext, err := c.DecryptDecode(payload)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(plaintext)
	//Do something with the plaintext
}

func Subscribe(client mqtt.Client, topic string, qos byte) {
	token := client.Subscribe(topic, qos, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
}

func Publish(client mqtt.Client, topic string, payload []byte) {
	token := client.Publish(topic, 0, false, payload)
	token.Wait()
}

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

}

var connectionHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Println("Connection lost")
}
