package main

import (
	"fmt"
	"testing"
	"wcws/dialogflow"

	"github.com/stretchr/testify/assert"
)

func TestGetThanksAnswer(t *testing.T) {
	name := "Hoang"
	rs := GetThanksAnswer(name)
	assert.Equal(t, rs, "")
}
func TestExtractAddressFromCoordinator(t *testing.T) {
	cor := dialogflow.Coordinates{
		Latitude:  16.074345,
		Longitude: 108.22385129999999,
	}
	str, err := ExtractAddressFromCoordinator(cor)
	assert.NoError(t, err)
	fmt.Sprintf(str)
}
