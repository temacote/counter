package monitoring

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
	"go.uber.org/zap"
)

type (
	Monitoring interface {
		Counter(*Metric, float64) error
		Inc(*Metric) error
		ExecutionTime(*Metric, func() error) error
		Val(*Metric, float64) error
	}

	Metric struct {
		Namespace   string
		Subsystem   string
		Name        string
		ConstLabels prometheus.Labels
	}

	PrometheusMonitoring struct {
		logger *zap.Logger

		collectors map[string]prometheus.Collector

		lock sync.Mutex

		pushURL    string
		username   string
		password   string
		jobName    string
		subProcess string

		disableLogPushError bool
	}

	collectorType int16
)

const (
	counter collectorType = iota
	histogram
	gauge
)

func NewPrometheusMonitoring(
	logger *zap.Logger,
	pushURL,
	username,
	password,
	jobName,
	subProcess string,
	disableLogPushError bool) Monitoring {
	var m = &PrometheusMonitoring{
		logger:              logger,
		collectors:          map[string]prometheus.Collector{},
		pushURL:             pushURL,
		username:            username,
		password:            password,
		jobName:             jobName,
		subProcess:          subProcess,
		disableLogPushError: disableLogPushError,
	}

	go m.push()
	return m
}

func (m *PrometheusMonitoring) push() {
	var ticker = time.NewTicker(time.Second * 5)
	for range ticker.C {
		var pushJob = push.
			New(m.pushURL, m.jobName).
			Gatherer(prometheus.DefaultGatherer).
			BasicAuth(m.username, m.password).
			Grouping("sub_process", m.subProcess)
		if err := pushJob.Add(); err != nil && !m.disableLogPushError {
			m.logger.Error("error push metrics", zap.String("url", m.pushURL), zap.Error(err))
		}
	}
}

func (m *PrometheusMonitoring) Counter(metric *Metric, count float64) (err error) {
	var collector prometheus.Collector
	if collector, err = m.collector(counter, metric); err != nil {
		return err
	}

	var (
		counter prometheus.Counter
		ok      bool
	)
	if counter, ok = collector.(prometheus.Counter); !ok {
		return errors.New(fmt.Sprintf("incorrect collector type. Required prometheus.Counter, got %T", collector))
	}

	counter.Add(count)

	return nil
}

func (m *PrometheusMonitoring) Inc(metric *Metric) (err error) {
	return m.Counter(metric, 1)
}

func (m *PrometheusMonitoring) Val(metric *Metric, val float64) (err error) {
	var collector prometheus.Collector
	if collector, err = m.collector(gauge, metric); err != nil {
		return err
	}

	var (
		counter prometheus.Gauge
		ok      bool
	)
	if counter, ok = collector.(prometheus.Gauge); !ok {
		return errors.New(fmt.Sprintf("incorrect collector type. Required prometheus.Gauge, got %T", collector))
	}

	counter.Set(val)
	return nil
}

func (m *PrometheusMonitoring) ExecutionTime(metric *Metric, h func() error) (err error) {
	var collector prometheus.Collector
	if collector, err = m.collector(histogram, metric); err != nil {
		return h()
	}

	var (
		histogram prometheus.Histogram
		ok        bool
	)
	if histogram, ok = collector.(prometheus.Histogram); !ok {
		return h()
	}

	var start = time.Now()
	err = h()
	histogram.Observe(time.Since(start).Seconds())
	return err
}

func (m *PrometheusMonitoring) Handler() http.Handler {
	return promhttp.Handler()
}

func (m *PrometheusMonitoring) collector(t collectorType, metric *Metric) (collector prometheus.Collector, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	var name = metric.String()
	if _, ok := m.collectors[name]; !ok {
		switch t {
		case counter:
			m.collectors[name] = prometheus.NewCounter(prometheus.CounterOpts{
				Namespace:   metric.Namespace,
				Subsystem:   metric.Subsystem,
				Name:        metric.Name,
				ConstLabels: metric.ConstLabels,
				Help:        "counter " + metric.Name,
			})
		case histogram:
			m.collectors[name] = prometheus.NewHistogram(prometheus.HistogramOpts{
				Namespace:   metric.Namespace,
				Subsystem:   metric.Subsystem,
				Name:        metric.Name,
				ConstLabels: metric.ConstLabels,
				Help:        "histogram " + metric.Name,
			})
		case gauge:
			m.collectors[name] = prometheus.NewGauge(
				prometheus.GaugeOpts{
					Namespace:   metric.Namespace,
					Subsystem:   metric.Subsystem,
					Name:        metric.Name,
					ConstLabels: metric.ConstLabels,
					Help:        "gauge " + metric.Name,
				},
			)
		}

		if err = prometheus.Register(m.collectors[name]); err != nil {
			return nil, err
		}
	}
	return m.collectors[name], nil
}

func (m Metric) String() string {
	var (
		keys = make([]string, len(m.ConstLabels))
		i    int
	)
	for k := range m.ConstLabels {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	var h = sha1.New()
	h.Write([]byte(m.Namespace + m.Subsystem + m.Name))

	for _, k := range keys {
		h.Write([]byte(k + m.ConstLabels[k]))
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
