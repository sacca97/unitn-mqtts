# MQTTSecure

## Description
MQTTSecure is a MQTT wrapper that adds end-to-end encryption to the protocol. It Ciphertext Policy Attribute Based Encryption (CP-ABE) or Symmetric Encryption (AESGCM, CHACHA20) and it is compatible with MQTTv3 and MQTTv5.

This repo is part of the work done for:

S. Gupta, T. Sacchetti and B. Crispo, "End-to-End Encryption for Securing Communications in Industry 4.0," 2022 4th IEEE Middle East and North Africa COMMunications Conference (MENACOMM), Amman, Jordan, 2022, pp. 153-158, doi: 10.1109/MENACOMM57252.2022.9998272.

(https://ieeexplore.ieee.org/document/9998272)


## Installation
- Install the package in your Go environment
```bash
go get "github.com/sacca97/unitn-mqtts"
```

## Usage
- Import the package 
```go
import mqtt "github.com/sacca97/unitn-mqtts"
```

- Create a MQTTConfig variable with custom settings

```go
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
		Version:              V3,
		Publisher:            true,
		EncryptionKey:        "public.key",
	}
```

- Connect to a MQTT broker (or more)

```go
client, err := config.CreateConnection()
if err != nil {
    log.Fatal(err)
}
```

## Authors and acknowledgment

- The wrapper structure is based on [alihanyalcin/mqtt-wrapper library](https://github.com/alihanyalcin/mqtt-wrapper)

- The cryptographic modules are based on [flynn/noise library](https://github.com/flynn/noise)
