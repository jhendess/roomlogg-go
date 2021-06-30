package sensor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getSensorDataFromBytes_positive(t *testing.T) {
	rawData := make([]byte, 64)
	rawData[0] = 0x7b
	rawData[1] = 0x00
	rawData[2] = 0xea
	rawData[3] = 0x38

	parsed := getSensorDataFromBytes(rawData)
	assert.Equal(t, float32(23.4), parsed[0].Temperature)
	assert.Equal(t, 56, int(parsed[0].Humidity))
}

func Test_getSensorDataFromBytes_negative(t *testing.T) {
	rawData := make([]byte, 64)
	rawData[1] = 0xff
	rawData[2] = 0xee
	rawData[3] = 0x26

	parsed := getSensorDataFromBytes(rawData)
	assert.Equal(t, float32(-1.8), parsed[0].Temperature)
	assert.Equal(t, 38, int(parsed[0].Humidity))
}
