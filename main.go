package main

import (
	"github.com/maodijim/tesla-api-go"
	"github.com/maodijim/tesla_automated_charge_control/plugins"
	log "github.com/sirupsen/logrus"
	"math"
	"net/http"
	_ "net/http/pprof"
	"time"
)

const (
	ver = "v1.0.0"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	configs := NewConfigs("")
	var ps []plugins.Plugins
	ps = plugins.LoadPlugins("")
	if ps == nil {
		log.Fatalf("no plugin found")
	}
	teslaApi := tesla.NewTeslaApi("", "", configs.TeslaRefreshToken, true)
	err := teslaApi.Login()
	if err != nil {
		log.Fatal(err)
	}
	vehicles, err := teslaApi.ListVehicles()
	if err != nil {
		log.Fatal(err)
	}
	err = teslaApi.SetActiveVehicle(vehicles[0])
	if err != nil {
		log.Fatalf("failed to set active tesla vehicle: %s", err)
	}
	isCharging := false
	for {
		for _, p := range ps {
			solarWatt := p.GetSolarWatt()
			gridWatt := p.GetGridWatt()
			p.Close()
			log.Infof("current solar watt: %f", solarWatt)
			log.Infof("current grid watt: %f", gridWatt)

			// Extra watt generated from solar available for charging
			extraWatt := solarWatt - gridWatt

			if inOffPeakTou(configs.Tou) || extraWatt > 0 {

				if extraWatt > float64(configs.WattGap) {
					chargeState, err := teslaApi.ChargeState()
					isCharging = teslaApi.IsCharging()
					powerConnected := isPowerConnected(*chargeState)
					chargeReached := chargeLimitReached(*chargeState)
					// skip if battery full or power not connected
					if chargeReached || !powerConnected {
						time.Sleep(10 * time.Minute)
						continue
					}
					log.Infof("%s is charging: %t", teslaApi.GetVehicleName(), isCharging)
					if err != nil {
						log.Errorf("failed to get charge state: %s", err)
						continue
					}
					currentVolt := chargeState.ChargerVoltage
					currentAmp := chargeState.ChargeAmps
					chargeAmp := int(math.Floor(extraWatt/float64(currentVolt))) + currentAmp
					if isCharging && !teslaApi.IsFastCharging() && currentAmp != chargeAmp {
						log.Infof("increasing charging amp to %d", chargeAmp)
						_, err = teslaApi.SetChargeAmps(chargeAmp)
						if err != nil {
							log.Errorf("failed to change amp: %s", err)
						}
					} else if configs.AutoChargeStart && !isCharging && !teslaApi.IsFastCharging() {
						log.Infof("start charging %s", teslaApi.GetVehicleName())
						_, err = teslaApi.ChargeStart()
						if err != nil {
							log.Errorf("failed to start charging %s", teslaApi.GetVehicleName())
						}
						time.Sleep(5 * time.Second)
						continue
					}
				} else if extraWatt < 0 && solarWatt > 200 {
					chargeState, err := teslaApi.ChargeState()
					isCharging = teslaApi.IsCharging()
					if isCharging && !teslaApi.IsFastCharging() {
						if err != nil {
							log.Errorf("failed to get charge state: %s", err)
						}
						currentVolt := chargeState.ChargerVoltage
						currentAmp := chargeState.ChargeAmps
						currentChargeWatt := currentAmp * currentVolt
						chargeAmp := currentAmp - int(math.Ceil(math.Abs(extraWatt)/float64(currentVolt)))
						chargeAmp = int(math.Max(float64(chargeAmp), 1))
						if int(solarWatt) > currentChargeWatt && currentAmp != chargeAmp {
							log.Infof("lowering charging amp to %d", chargeAmp)
							_, err = teslaApi.SetChargeAmps(chargeAmp)
							if err != nil {
								log.Errorf("failed to change amp: %s", err)
							}
						}
					}
				} else if extraWatt < 0 && configs.AutoChargeStop {
					if isCharging {
						_ = stopCharging(teslaApi)
					} else {
						log.Infof("%s is not charging", teslaApi.GetVehicleName())
					}
				}
				time.Sleep(time.Minute * 3)
			} else {
				if isCharging {
					err = stopCharging(teslaApi)
					if err == nil {
						isCharging = false
					}
				}
				log.Infof("currently not in off peak time")
				time.Sleep(time.Minute * 15)
			}
		}
	}
}

func stopCharging(api *tesla.TeslaApi) (err error) {
	_, err = api.ChargeStop()
	if err != nil {
		log.Errorf("failed to stop %s charge: %s", api.GetVehicleName(), err)
	} else {
		log.Infof("charging stopped")
	}
	return err
}

func isPowerConnected(state tesla.ChargeState) bool {
	return state.ChargingState != "Disconnected"
}

func chargeLimitReached(state tesla.ChargeState) bool {
	return state.BatteryLevel >= state.ChargeLimitSoc
}

func inOffPeakTou(tou map[string]ConfigTou) bool {
	if op, ok := tou["off-peak"]; ok {
		now := time.Now()
		a, err := time.Parse("15:04", op.Start)
		b, err := time.Parse("15:04", op.End)
		if err != nil {
			log.Errorf("%s", err)
		}
		if a == b {
			return true
		}
		start := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			a.Hour(),
			a.Minute(),
			0,
			0,
			now.Location(),
		)
		end := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			b.Hour(),
			b.Minute(),
			0,
			0,
			now.Location(),
		)
		if a.After(b) && now.After(end) {
			end = end.Add(time.Hour * 24)
		} else if a.After(b) && now.Before(start) {
			start = start.Add(time.Hour * -24)
		}
		if now.After(start) && now.Before(end) {
			return true
		}
		return false
	} else {
		return true
	}
}
