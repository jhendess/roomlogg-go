package main

import (
	"github.com/jhendess/roomlogg-go/sensor"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// Test when sensor result is returned twice
func Test_checkForLostSensors_ok(t *testing.T) {
	lastSensorResponseMap = make(map[int]*LastSensorResponse, 0)

	testSensors1 := make([]*sensor.Sensor, 0)
	testSensors1 = append(testSensors1, &sensor.Sensor{
		Channel:     0,
		Temperature: 2,
		Humidity:    2,
		Absent:      false,
	})
	checkForLostSensors(testSensors1)
	require.NotNil(t, lastSensorResponseMap[0])
	require.False(t, lastSensorResponseMap[0].absent)
	require.Truef(t, lastSensorResponseMap[0].lastResponse.Add(time.Second).After(time.Now()), "Last response time must be ~ %v but was %v", time.Now(), lastSensorResponseMap[0].lastResponse)

	testSensors2 := make([]*sensor.Sensor, 0)
	testSensors2 = append(testSensors2, &sensor.Sensor{
		Channel:     0,
		Temperature: 2,
		Humidity:    2,
		Absent:      false,
	})
	checkForLostSensors(testSensors2)
	require.NotNil(t, lastSensorResponseMap[0])
	require.False(t, lastSensorResponseMap[0].absent)
}

func Test_checkForLostSensors_lost(t *testing.T) {
	lastSensorResponseMap = make(map[int]*LastSensorResponse, 0)

	testSensors1 := make([]*sensor.Sensor, 0)
	testSensors1 = append(testSensors1, &sensor.Sensor{
		Channel:     0,
		Temperature: 2,
		Humidity:    2,
		Absent:      false,
	})
	checkForLostSensors(testSensors1)
	require.NotNil(t, lastSensorResponseMap[0])
	require.False(t, lastSensorResponseMap[0].absent)
	lastSensorResponseMap[0].lastResponse = lastSensorResponseMap[0].lastResponse.Add(-1 * time.Second * (secondsBeforeLost + 1))
	lastResponseTime := lastSensorResponseMap[0].lastResponse

	time.Sleep(time.Millisecond * 500)

	// Sensor is now lost
	testSensors2 := make([]*sensor.Sensor, 0)
	testSensors2 = append(testSensors2, &sensor.Sensor{
		Channel:     0,
		Temperature: 0,
		Humidity:    0,
		Absent:      true,
	})
	checkForLostSensors(testSensors2)
	require.True(t, lastSensorResponseMap[0].absent)
	require.Equal(t, lastResponseTime, lastSensorResponseMap[0].lastResponse)

	// Recover the sensor
	testSensors3 := make([]*sensor.Sensor, 0)
	testSensors3 = append(testSensors3, &sensor.Sensor{
		Channel:     0,
		Temperature: 2,
		Humidity:    2,
		Absent:      false,
	})
	checkForLostSensors(testSensors3)
	require.False(t, lastSensorResponseMap[0].absent)
	lastSensorResponseMap[0].lastResponse = lastSensorResponseMap[0].lastResponse.Add(-1 * time.Second * (secondsBeforeLost + 1))
}
