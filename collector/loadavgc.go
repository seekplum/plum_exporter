// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collector

// #include <stdlib.h>
import "C"
import (
	"fmt"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type loadavgCCollector struct {
	metric []prometheus.Gauge
}

func init() {
	log.Debugln("load avg init")
	Factories["loadavgc"] = NewLoadavgCCollector
}

// Take a prometheus registry and return a new Collector exposing load average.
func NewLoadavgCCollector() (Collector, error) {
	return &loadavgCCollector{
		metric: []prometheus.Gauge{
			prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "cload1",
				Help:      "1m load average.",
			}),
			prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "cload5",
				Help:      "5m load average.",
			}),
			prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: Namespace,
				Name:      "cload15",
				Help:      "15m load average.",
			}),
		},
	}, nil
}

func getLoadC() ([]float64, error) {
	var loadavg [3]C.double
	samples := C.getloadavg(&loadavg[0], 3)
	if samples > 0 {
		return []float64{float64(loadavg[0]), float64(loadavg[1]), float64(loadavg[2])}, nil
	} else {
		return nil, errors.New("failed to get load average")
	}
}

func (c *loadavgCCollector) Update(ch chan<- prometheus.Metric) (err error) {
	loads, err := getLoadC()
	if err != nil {
		return fmt.Errorf("couldn't get load: %s", err)
	}
	for i, load := range loads {
		log.Debugf("Set load %d: %f", i, load)
		c.metric[i].Set(load)
		c.metric[i].Collect(ch)
	}
	return err
}
