package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"net/http"
	"os"
	"sync"
	"time"
	//"github.com/seekplum/plum_exporter/collector" // 绝对路径导入
	"plum_exporter/collector" // 相对路径导入
)

var (
	requestMax = flag.Int(
		"config.request-max", 2, "max requestions coexist",
	)

	listenAddress = flag.String(
		"web.listen-address", ":10002",
		"Address to listen on for web interface and telemetry.",
	)
	showVersion = flag.Bool(
		"version", false,
		"Print version information.",
	)
	scrapeDurations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: collector.Namespace,
			Subsystem: "exporter",
			Name:      "scrape_duration_seconds",
			Help:      "plum_exporter: Duration of a scrape job.",
		},
		[]string{"collector", "result"},
	)
)

// PlumCollector implements the prometheus.Collector interface.
type PlumCollector struct {
	collectors map[string]collector.Collector
}

// Describe implements the prometheus.Collector interface.
func (r PlumCollector) Describe(ch chan<- *prometheus.Desc) {
	scrapeDurations.Describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (r PlumCollector) Collect(ch chan<- prometheus.Metric) {
	//a := prometheus.NewGauge(prometheus.GaugeOpts{
	//	Namespace: "node",
	//	Name:      "load1",
	//	Help:      "1m load average.",
	//})
	//a.Set(1)
	//a.Collect(ch)
	//delayDesc := prometheus.NewDesc(
	//	prometheus.BuildFQName("node", "", "delay"),
	//	"Get peers delay",
	//	[]string{},
	//	nil,
	//)
	//ch <- prometheus.MustNewConstMetric(delayDesc, prometheus.GaugeValue,1)
	//return

	log.Infoln("plum collector")
	wg := sync.WaitGroup{}
	wg.Add(len(r.collectors))
	for name, c := range r.collectors {
		log.Infof("collect name: %s", name)
		go func(name string, c collector.Collector) {
			execute(name, c, ch)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
	scrapeDurations.Collect(ch)
}

func execute(name string, c collector.Collector, ch chan<- prometheus.Metric) {
	begin := time.Now()
	err := c.Update(ch)
	duration := time.Since(begin)
	var result string

	if err != nil {
		log.Errorf("ERROR: %s collector failed after %fs: %s", name, duration.Seconds(), err)
		result = "error"
	} else {
		log.Infof("OK: %s collector succeeded after %fs.", name, duration.Seconds())
		result = "success"
	}
	scrapeDurations.WithLabelValues(name, result).Observe(duration.Seconds())
}

func loadCollectors() (map[string]collector.Collector, error) {
	collectors := map[string]collector.Collector{}
	log.Infof("collector Factories : %v", collector.Factories)
	for name, fn := range collector.Factories {
		log.Infof("load collect %s", name)
		c, err := fn()
		if err != nil {
			return nil, fmt.Errorf("collector '%s' not available", name)
		}
		collectors[name] = c
	}
	return collectors, nil
}

func handler(w http.ResponseWriter, r *http.Request) {

	registry := prometheus.NewRegistry()
	collectors, err := loadCollectors()
	if err != nil {
		log.Fatalf("Couldn't load collectors: %s", err)
	}
	c := PlumCollector{collectors}
	registry.MustRegister(c)
	// Delegate http serving to Promethues client library, which will call collector.Collect.
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	// send content
	h.ServeHTTP(w, r)
}

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("plum_exporter"))
		os.Exit(0)
	}
	if *requestMax < 1 {
		log.Fatal("config.request-max must be more than 0")
	}
	log.Infoln("Starting plum", version.Info())
	log.Infoln("Build context", version.BuildContext())
	defer func() {
		if err := recover(); err != nil {

		}
	}()

	http.HandleFunc("/metrics", handler)
	http.Handle("/profile_metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer, promhttp.HandlerOpts{ErrorLog: log.NewErrorLogger()},
	))

	log.Infoln("Listening on", *listenAddress)
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
