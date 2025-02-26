package sxfad

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestAPI(t *testing.T) {
	// TODO: Implement test cases
	a := Api{
		Url:      "https://192.168.167.9",
		User:     "readonly",
		Password: "",
	}
	level, _ := logrus.ParseLevel("info")
	logrus.SetLevel(level)
	a.Monitor("/api/lb/current-version/stat/sys/system", "system")

}
