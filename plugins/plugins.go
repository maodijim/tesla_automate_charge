package plugins

import log "github.com/sirupsen/logrus"

var (
	ps = []Plugins{
		&SensePlugin{},
	}
)

type Plugins interface {
	NewPlugin(pluginConfigFile string) error
	GetSolarWatt() float64
	GetGridWatt() float64
	Close()
}

func LoadPlugins(confFile string) (pis []Plugins) {
	for _, p := range ps {
		err := p.NewPlugin(confFile)
		if err != nil {
			log.Error(err)
			continue
		}
		pis = append(pis, p)
	}
	return pis
}
