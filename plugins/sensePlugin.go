package plugins

import (
	"errors"
	"fmt"
	"github.com/maodijim/sense-api"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

type SenseConfig struct {
	SenseUser         string `yaml:"senseUser"`
	SensePass         string `yaml:"sensePass"`
	SenseRefreshToken string `yaml:"senseRefreshToken"`
}

type SensePlugin struct {
	Configs  SenseConfig
	senseApi sense.SenseApi
}

func (s *SensePlugin) GetSolarWatt() (watt float64) {
	msg, err := s.senseApi.ReadMessage()
	if err != nil || msg == nil {
		log.Errorf("error %s: %s", err, msg)
		return watt
	}
	for {
		switch msg.Type {
		case sense.PayloadDataChange, sense.PayloadHello, sense.PayloadMonitorInfo:
			msg, _ = s.senseApi.ReadMessage()
		case sense.PayloadRealTimeUpdate:
			watt = msg.Payload.SolarW
			return watt
		default:
			return watt
		}
	}
}

func (s *SensePlugin) Close() {
	s.senseApi.Close()
}

func (s *SensePlugin) GetGridWatt() (watt float64) {
	msg, err := s.senseApi.ReadMessage()
	if err != nil || msg == nil {
		log.Errorf("error %s: %s", err, msg)
		return watt
	}
	for {
		switch msg.Type {
		case sense.PayloadDataChange, sense.PayloadHello, sense.PayloadMonitorInfo:
			msg, _ = s.senseApi.ReadMessage()
		case sense.PayloadRealTimeUpdate:
			watt = msg.Payload.W
			return watt
		default:
			return watt
		}
	}
}

func (s *SensePlugin) NewPlugin(pluginConfigFile string) (err error) {
	sc := SenseConfig{}
	if pluginConfigFile == "" {
		pluginConfigFile = "configs.yml"
	}
	f, err := os.ReadFile(pluginConfigFile)
	if err != nil {
		errMsg := fmt.Sprintf("sense plugin failed to read plug in config %s: %s", pluginConfigFile, err)
		err = errors.New(errMsg)
		return err
	}
	err = yaml.Unmarshal(f, &sc)
	if err != nil {
		errMsg := fmt.Sprintf("sense plugin failed to parse plug in config: %s", err)
		err = errors.New(errMsg)
		return err
	}
	s.Configs = sc
	sApi, err := sense.NewSenseApi(s.Configs.SenseUser, s.Configs.SensePass)
	if err != nil {
		errMsg := fmt.Sprintf("sense plugin failed to create: %s", err)
		err = errors.New(errMsg)
		return err
	}
	s.senseApi = *sApi
	return err
}
