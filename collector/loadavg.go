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

import (
	"fmt"
	"errors"
	"regexp"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/seekplum/plum_exporter/utils"
	"github.com/seekplum/plum_exporter/config"
)

type loadavgCollector struct {
	metric []prometheus.Gauge
}

func init() {
	log.Debugln("load avg init")
	Factories["loadavg"] = NewLoadavgCollector
}

// Take a prometheus registry and return a new Collector exposing load average.
func NewLoadavgCollector() (Collector, error) {
	return &loadavgCollector{}, nil
}

func getLoad() (map[string]float64, error) {
	var loadInfo = map[string]float64{}
	output, err := utils.Cmd(config.UPTIME)
	if err == nil {
		pattern := regexp.MustCompile(`load averages?:\s+(\d+\.?\d+),?\s+(\d+\.?\d+),?\s+(\d+\.?\d+)`)
		match := pattern.FindStringSubmatch(output)
		load1, _ := strconv.ParseFloat(match[1], 64)
		load5, _ := strconv.ParseFloat(match[2], 64)
		load15, _ := strconv.ParseFloat(match[3], 64)
		loadInfo["load1"] = load1
		loadInfo["load5"] = load5
		loadInfo["load15"] = load15
		return loadInfo, nil
	} else {
		return nil, errors.New("failed to get load average")
	}
}

func (c *loadavgCollector) Update(ch chan<- prometheus.Metric) (err error) {
	loads, err := getLoad()
	if err != nil {
		return fmt.Errorf("couldn't get load: %s", err)
	}
	log.Debugf("Set uptime load: %v", loads)
	for k, v := range loads {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(Namespace, "", k),
				fmt.Sprintf("%s load averages.", k),
				nil, nil,
			),
			prometheus.GaugeValue, v,
		)
	}
	return err
}
