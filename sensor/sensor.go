package sensor

import (
	"encoding/binary"
	"github.com/karalabe/hid"
	"log"
	"time"
)

type Sensor struct {
	Channel     int
	Temperature float32
	Humidity    byte
	Absent      bool
}

func QueryAndPrintOnce(deviceInfo []hid.DeviceInfo) {
	sensors, err := QueryDeviceSensors(deviceInfo)
	if err != nil {
		log.Fatal("Querying sensor failed")
	}
	for i := range sensors {
		sensor := sensors[i]
		if sensor.Absent {
			log.Printf("No sensor on Channel %d", sensor.Channel)
		} else {
			log.Printf("Channel %d: %.1f°C  %d%%\n", sensor.Channel, sensor.Temperature, sensor.Humidity)
		}
	}
}

func QueryDeviceSensors(deviceInfo []hid.DeviceInfo) ([]*Sensor, error) {
	log.Printf("Opening device %v...\n", deviceInfo)

	device, err := deviceInfo[0].Open()
	if err != nil {
		log.Printf("Opening device failed: %v\n", err)
		return nil, err
	}
	defer closeDevice(device)

	requestBytes := getTempRequestBytes()
	log.Printf("Writing to device: %x", requestBytes)
	writtenBytes, err := device.Write(requestBytes)
	if err != nil {
		log.Printf("Writing to device failed: %v\n", err)
		return nil, err
	}
	log.Printf("Wrote %v bytes", writtenBytes)

	time.Sleep(500 * time.Millisecond)

	response := make([]byte, 64)

	readBytes, err := device.Read(response)
	if err != nil {
		log.Printf("Reading from device failed: %v\n", err)
		return nil, err
	}
	log.Printf("Read result (%v bytes): %x", readBytes, response)

	sensorData := getSensorDataFromBytes(response)
	return sensorData, nil
}

func getSensorDataFromBytes(response []byte) []*Sensor {
	sensors := make([]*Sensor, 7)
	for c := 0; c < 7; c++ {
		var sensor *Sensor
		channelDataStart := 1 + (c * 3)
		if response[channelDataStart+2] != 0xff {
			temperatureRaw := binary.BigEndian.Uint16(response[channelDataStart : channelDataStart+2])
			temperatureSigned := int16(temperatureRaw)
			temperature := float32(temperatureSigned) / 10
			humidity := response[channelDataStart+2]
			sensor = &Sensor{
				Channel:     c,
				Temperature: temperature,
				Humidity:    humidity,
			}
		} else {
			sensor = &Sensor{Channel: c, Absent: true}
		}
		sensors[c] = sensor
	}
	return sensors
}

func getTempRequestBytes() []byte {
	b := make([]byte, 64)
	b[0] = 0x7b
	b[1] = 0x03
	b[2] = 0x40
	b[3] = 0x7d
	return b
}

func closeDevice(device *hid.Device) {
	log.Println("Closing device...")
	err := device.Close()
	if err != nil {
		log.Fatalln("Closing device failed")
	}
}
