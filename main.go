package mqtts

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func startPublisher() {
	config := Config{
		Brokers:              []string{"localhost:1883"},
		ClientID:             "publisherNero",
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
		Publisher:            true,
	}

	client, err := config.CreateConnection()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		start := time.Now()
		err = client.Publish("test/simola", "((0 AND 1) OR (2 AND 3)) AND 5", "Hello World!")
		elapsed := time.Since(start)
		log.Printf("Published in %s", elapsed)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(2 * time.Second)
	}
	defer client.Disconnect()

}

func startSubscriber() {
	config := Config{
		Brokers:              []string{"localhost:1883"},
		ClientID:             "subscriberNero",
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
		Publisher:            false,
	}
	client, err := config.CreateConnection()
	if err != nil {
		log.Fatal(err)
	}

	client.Handle(func(topic string, payload []byte) {
		log.Printf("%s: %s", topic, string(payload))
	})

	//defer client.Disconnect()
}

func Main() {
	keepAlive := make(chan os.Signal, 1)
	signal.Notify(keepAlive, os.Interrupt, syscall.SIGTERM)
	startSubscriber()
	go startPublisher()
	<-keepAlive

}
