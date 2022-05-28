package mqtts

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
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
	policy := "(1 OR 2 OR 3 OR 4)"
	err = client.Publish("test/simola", policy, "SIMOLA DOVE E IL SUGO")
	if err != nil {
		log.Fatal(err)
	}

	//msg := "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"
	//policies := []string{"(1 AND 2 AND 3 AND 4 AND 5)", "(1 AND 2 AND 3 AND 4 AND 5 AND 6 AND 7 AND 8 AND 9 AND 10)", "(1 AND 2 AND 3 AND 4 AND 5 AND 6 AND 7 AND 8 AND 9 AND 10 AND 11 AND 12 AND 13 AND 14 AND 15)", "(1 AND 2 AND 3 AND 4 AND 5 AND 6 AND 7 AND 8 AND 9 AND 10 AND 11 AND 12 AND 13 AND 14 AND 15 AND 16 AND 17 AND 18 AND 19 AND 20)", "(1 AND 2 AND 3 AND 4 AND 5 AND 6 AND 7 AND 8 AND 9 AND 10 AND 11 AND 12 AND 13 AND 14 AND 15 AND 16 AND 17 AND 18 AND 19 AND 20 AND 21 AND 22 AND 23 AND 24 AND 25)", "(1 AND 2 AND 3 AND 4 AND 5 AND 6 AND 7 AND 8 AND 9 AND 10 AND 11 AND 12 AND 13 AND 14 AND 15 AND 16 AND 17 AND 18 AND 19 AND 20 AND 21 AND 22 AND 23 AND 24 AND 25 AND 26 AND 27 AND 28 AND 29 AND 30)"}
	//testORPolicy(nil, client, msg, policies)
	//testANDPolicy(&wg, client, msg, policies)
	//testSizeMsgIncrement(client)
	defer client.Disconnect()

}

func testSizeMsgIncrement(client MQTT) {
	policy := "(1 AND 2 OR 3 AND 4 OR 5 AND 6 OR 7 AND 8 OR 9 OR 10)"
	dir, err := os.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dir {
		var sum uint64 = 0
		var size int = 0
		var msg string = ""
		if strings.Contains(file.Name(), ".txt") {
			f, err := os.ReadFile(file.Name())
			if err != nil {
				log.Fatal(err)
			}
			msg = string(f)
			for i := 0; i < 3; i++ {

				start := time.Now()
				err = client.Publish("test/simola", policy, msg)
				elapsed := time.Since(start)
				intElapsed := elapsed.Milliseconds()
				sum += uint64(intElapsed)
				if err != nil {
					log.Fatal(err)
				}

			}
			fmt.Printf("%d;%d;%d\n", len(msg), sum/3, size/3)
		}

	}
	log.Println("FINISHED")
}

func testANDPolicy(wg *sync.WaitGroup, client MQTT, msg string, policies []string) {
	var sum uint64 = 0
	var size int = 0

	log.Println("Started AND policy")
	fmt.Println("literals;time,msgsize")

	for j := 0; j < len(policies); j++ {
		//policy := strings.Replace(policies[j], "AND", "OR", (strings.Count(policies[j], "AND")+1)/2)
		size = 0
		sum = 0
		for i := 0; i < 3; i++ {
			start := time.Now()
			err := client.Publish("test/simola", policies[j], msg)
			elapsed := time.Since(start)
			intElapsed := elapsed.Milliseconds()
			sum += uint64(intElapsed)
			if err != nil {
				log.Fatal(err)
			}
		}
		val := strings.Count(policies[j], "AND")
		fmt.Printf("%d;%d;%d\n", val+1, sum/3, size/3)
	}
	log.Println("FINISHED")

}

func testORPolicy(wg *sync.WaitGroup, client MQTT, msg string, policies []string) {
	log.Println("Started OR policy")
	fmt.Println("literals;time,msgsize")

	var sum int64 = 0
	size := 0
	for j := 0; j < len(policies); j++ {
		policy := strings.ReplaceAll(policies[j], "AND", "OR")
		sum = 0
		size = 0
		for i := 0; i < 3; i++ {
			start := time.Now()
			err := client.Publish("test/simola", policy, msg)
			elapsed := time.Since(start)
			intElapsed := elapsed.Milliseconds()
			sum += intElapsed

			if err != nil {
				log.Fatal(err)
			}
		}
		val := strings.Count(policy, "OR")
		fmt.Printf("%d;%d;%d\n", val+1, sum/3, size/3)
	}
	log.Println("FINISHED")

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
		//Do something with the message
		log.Printf("%s: %s", topic, string(payload))
	})

	//defer client.Disconnect()
}

func Main() {
	keepAlive := make(chan os.Signal, 1)
	signal.Notify(keepAlive, os.Interrupt, syscall.SIGTERM)
	go startSubscriber()

	startPublisher()
	<-keepAlive

}
