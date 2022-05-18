package main

import (
	"testing"
	"time"
)

func Test_inOffPeakTou(t *testing.T) {
	type args struct {
		tou map[string]ConfigTou
	}
	tests := []struct {
		name string
		args args
		want func() bool
	}{
		{
			name: "Test 23 to 15",
			args: args{
				tou: map[string]ConfigTou{
					"off-peak": {
						Start: "23:00",
						End:   "15:00",
					},
				},
			},
			want: func() bool {
				h := time.Now().Hour()
				switch {
				case h >= 23, h <= 15:
					return true
				default:
					return false
				}
			},
		},
		{
			name: "Test 3 to 18",
			args: args{
				tou: map[string]ConfigTou{
					"off-peak": {
						Start: "3:00",
						End:   "18:00",
					},
				},
			},
			want: func() bool {
				h := time.Now().Hour()
				switch {
				case h >= 3 && h <= 18:
					return true
				default:
					return false
				}
			},
		},
		{
			name: "Test 0 to 2",
			args: args{tou: map[string]ConfigTou{
				"off-peak": {
					Start: "0:00",
					End:   "2:00",
				},
			}},
			want: func() bool {
				h := time.Now().Hour()
				switch {
				case h >= 0 && h <= 2:
					return true
				default:
					return false
				}
			},
		},
		{
			name: "Test 0 30 to 0",
			args: args{tou: map[string]ConfigTou{
				"off-peak": {
					Start: "0:30",
					End:   "0:00",
				},
			}},
			want: func() bool {
				h := time.Now().Hour()
				switch {
				case h == 0 && h < 30 && h > 0:
					return false
				default:
					return true
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inOffPeakTou(tt.args.tou); got != tt.want() {
				t.Errorf("inOffPeakTou() = %v, want %v", got, tt.want())
			}
		})
	}
}

func Test_calculateDistance(t *testing.T) {
	type args struct {
		dsLat   float64
		dsLon   float64
		configs Configs
	}
	tests := []struct {
		name           string
		args           args
		withinDistance float64
	}{
		{
			name: "Test the distance from location",
			args: args{
				dsLat: 40.6892566,
				dsLon: -74.044766,
				configs: Configs{ChargeLocation: ChargeLocation{
					Lat: 40.6892568,
					Lon: -74.044776,
				}},
			},
			withinDistance: 300,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDistance := calculateDistance(tt.args.dsLat, tt.args.dsLon, tt.args.configs); gotDistance > tt.withinDistance {
				t.Errorf("isWithinChargeLocation() = %v, want less than %v", gotDistance, tt.withinDistance)
			}
		})
	}
}
