package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jhendess/roomlogg-go/sensor"
	"github.com/karalabe/hid"
	"log"
	"strconv"
	"strings"
	"time"
)

// Seconds after which a previously available sensor is considered "lost"
const secondsBeforeLost = 120

// Keep track of when which sensor responded the last time to detect battery failures
var lastSensorResponseMap map[int]*LastSensorResponse

// LastSensorResponse is used to
type LastSensorResponse struct {
	lastResponse time.Time
	absent       bool
}

func main() {
	checkHidSupport()

	queryOncePtr := flag.Bool("query", false, "Query sensors only once and then exit")
	serverPtr := flag.Bool("server", false, "Start a prometheus-compatible exporter server")
	serverPortPtr := flag.Int("port", 8080, "Port to listen on when starting a server")
	flag.Parse()

	deviceInfo := detectDevice()

	if queryOnce := *queryOncePtr; queryOnce {
		sensor.QueryAndPrintOnce(deviceInfo)
	} else if server := *serverPtr; server {
		lastSensorResponseMap = make(map[int]*LastSensorResponse, 0)
		startServer(deviceInfo, *serverPortPtr)
	} else {
		flag.Usage()
	}
}

func checkHidSupport() {
	if !hid.Supported() {
		log.Fatal("HID is not supported on this platform :(")
	}
}

func detectDevice() []hid.DeviceInfo {
	deviceInfo := hid.Enumerate(0x483, 0x5750)
	if len(deviceInfo) == 0 {
		log.Fatalln("No device found")
	} else if len(deviceInfo) > 1 {
		log.Fatalln("Only one device at the time is supported")
	}
	return deviceInfo
}

var globalDeviceInfo []hid.DeviceInfo

func startServer(deviceInfo []hid.DeviceInfo, port int) {
	if port < 80 || port > 65535 {
		log.Fatal("Port number must be between 80 and 65535")
	}
	addr := ":" + strconv.FormatInt(int64(port), 10)
	r := gin.Default()

	globalDeviceInfo = deviceInfo
	log.Printf("Starting server on %s", addr)

	r.GET("/metrics", exporterFunc)

	err := r.Run(addr)
	if err != nil {
		log.Fatalf("Starting server failed: %v", err)
	}
}

func exporterFunc(c *gin.Context) {
	sensors, err := sensor.QueryDeviceSensors(globalDeviceInfo)
	if err != nil {
		log.Printf("Unexpected error: %s\n", err)
		c.String(500, "Internal server error")
	} else {
		value := buildExporterData(sensors)
		checkForLostSensors(sensors)
		c.String(200, value)
	}
}

func buildExporterData(sensors []*sensor.Sensor) string {
	var builder strings.Builder
	builder.Grow(1024)

	for _, s := range sensors {
		if !s.Absent {
			_, _ = fmt.Fprintf(&builder, "room_temperature{channel=\"%d\"} %.1f\n", s.Channel, s.Temperature)
			_, _ = fmt.Fprintf(&builder, "room_humidity{channel=\"%d\"} %d\n", s.Channel, s.Humidity)
		}
	}

	return builder.String()
}

func checkForLostSensors(sensors []*sensor.Sensor) {
	for _, s := range sensors {
		if _, ok := lastSensorResponseMap[s.Channel]; !ok {
			lastSensorResponseMap[s.Channel] = &LastSensorResponse{
				lastResponse: time.Time{},
				absent:       s.Absent,
			}
		}
		lastSensorResponse := lastSensorResponseMap[s.Channel]
		if !s.Absent {
			if lastSensorResponse.absent {
				log.Printf("Sensor %d was recovered", s.Channel)
			}
			lastSensorResponse.absent = false
			lastSensorResponse.lastResponse = time.Now()
		} else if !lastSensorResponse.absent && s.Absent &&
			lastSensorResponse.lastResponse.Add(time.Second*secondsBeforeLost).Before(time.Now()) {
			log.Printf("Sensor %d has been absent for %d seconds", s.Channel, secondsBeforeLost)
			lastSensorResponse.absent = true
		}
	}
}
