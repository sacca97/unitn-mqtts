package mqtts

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Main() {
	keepAlive := make(chan os.Signal, 1)
	signal.Notify(keepAlive, os.Interrupt, syscall.SIGTERM)
	config := Config{
		Brokers:              []string{"localhost:1883"},
		ClientID:             "nero",
		Username:             "",
		Password:             "",
		Topics:               []string{"test/simola"},
		QoS:                  0,
		Retained:             false,
		AutoReconnect:        true,
		MaxReconnectInterval: 5,
		PersistentSession:    false,
		KeepAlive:            15,
		TLSCA:                "",
		TLSCert:              "",
		TLSKey:               "",
		Version:              V3, // use mqtt.V5 for MQTTv5 client.
	}
	client, err := config.CreateConnection()
	if err != nil {
		log.Fatal(err)
	}
	client.SetKeys()

	client.Handle(func(topic string, payload []byte) {
		plaintext, err := client.Decrypt(payload)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%s: %s", topic, plaintext)
	})
	err = client.Publish("test/simola", "((0 AND 1) OR (2 AND 3)) AND 5", "Hello World!")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()
	<-keepAlive

}
