package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

type ConfigTou struct {
	Start string `yaml:"start"`
	End   string `yaml:"end"`
}

type Configs struct {
	TeslaRefreshToken string               `yaml:"teslaRefreshToken"`
	Tou               map[string]ConfigTou `yaml:"tou"`
	AutoChargeStop    bool                 `yaml:"autoChargeStop"`
	AutoChargeStart   bool                 `yaml:"autoChargeStart"`
	WattGap           int                  `yaml:"wattGap"`
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
	err = yaml.Unmarshal(r, &c)
	if err != nil {
		log.Fatalln(err)
	}
	return c
}
