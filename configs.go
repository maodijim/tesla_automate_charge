package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type ConfigTou struct {
	Start string `yaml:"start"`
	End   string `yaml:"end"`
}

type ChargeLocation struct {
	Lon    float64 `yaml:"longitude"`
	Lat    float64 `yaml:"latitude"`
	Radius float64 `yaml:"radius"`
}

type Configs struct {
	TeslaRefreshToken string               `yaml:"teslaRefreshToken"`
	Tou               map[string]ConfigTou `yaml:"tou"`
	AutoChargeStop    bool                 `yaml:"autoChargeStop"`
	AutoChargeStart   bool                 `yaml:"autoChargeStart"`
	WattGap           int                  `yaml:"wattGap"`
	ChargeOnlyOffPeak bool                 `yaml:"chargeOnlyOffPeak"`
	ChargeLocation    ChargeLocation       `yaml:"chargeLocation"`
}

func (c Configs) IsLocationSet() bool {
	return c.ChargeLocation.Lat != 0 && c.ChargeLocation.Lon != 0 &&
		c.ChargeLocation.Lon >= -180 && c.ChargeLocation.Lon <= 180 &&
		c.ChargeLocation.Lat >= -90 && c.ChargeLocation.Lat <= 90
}

func NewConfigs(confPath string) (c Configs) {
	if confPath == "" {
		confPath = "configs.yml"
	}
	r, err := os.ReadFile(confPath)
	if err != nil {
		log.Fatalln(err)
	}
	c = Configs{}
	c.ChargeLocation.Radius = 300
	err = yaml.Unmarshal(r, &c)
	if err != nil {
		log.Fatalln(err)
	}
	return c
}
