package webfinger

import (
	"encoding/json"
	"io/ioutil"
)

type fingerprint struct {
	CMS      string   `json:"cms"`
	Method   string   `json:"method"`
	Location string   `json:"location"`
	Keyword  []string `json:"keyword"`
}

type PackFinger struct {
	Fingerprint []fingerprint `json:"fingerprint"`
}

func LoadFingerPrint(finger string) (*PackFinger, error) {
	var config PackFinger
	cmsfile, err := ioutil.ReadFile(finger)
	if err != nil {
		println(err)
	}
	err1 := json.Unmarshal(cmsfile, &config)
	if err1 != nil {
		return nil, err1
	}
	return &config, nil
}
