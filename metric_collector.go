package honeybadger

import (
	"fmt"
	"time"
)

type counterMetric struct {
	name  string
	value int
}

type timingMetric struct {
	name  string
	value time.Duration
}

type metricCollector struct {
	counter  chan counterMetric
	timingCh chan timingMetric
}

func (c *metricCollector) increment(metric string, value int) {
	c.counter <- counterMetric{metric, value}
}

func (c *metricCollector) timing(metric string, value time.Duration) {
	c.timingCh <- timingMetric{metric, value}
}

func newMetricCollector(config *Configuration, worker worker) *metricCollector {
	counter := make(chan counterMetric)
	timing := make(chan timingMetric)

	go func() {
		flush := make(chan int)

		go func() {
			for {
				time.Sleep(config.MetricsInterval)
				flush <- 1
			}
		}()

		counters := make(map[string]int)
		timings := make(map[string][]time.Duration)

		for {
			select {
			case metric := <-counter:
				counters[metric.name] += metric.value
			case metric := <-timing:
				timings[metric.name] = append(timings[metric.name], metric.value)
			case <-flush:
				if len(counters) <= 0 && len(timings) <= 0 {
					break
				}
				metrics := []string{}
				for name, value := range counters {
					metrics = append(metrics, fmt.Sprintf("%v %v", name, value))
				}
				for name, values := range timings {
					mean := mean(values)
					metrics = append(metrics, fmt.Sprintf("%v:mean %v", name, mean))
					metrics = append(metrics, fmt.Sprintf("%v:median %v", name, median(values)))
					metrics = append(metrics, fmt.Sprintf("%v:percentile_90 %v", name, percentile(values, 90)))
					metrics = append(metrics, fmt.Sprintf("%v:min %v", name, min(values)))
					metrics = append(metrics, fmt.Sprintf("%v:max %v", name, max(values)))
					metrics = append(metrics, fmt.Sprintf("%v:stddev %v", name, stdDev(values, mean)))
					metrics = append(metrics, fmt.Sprintf("%v %v", name, len(values)))
				}
				counters = make(map[string]int)
				timings = make(map[string][]time.Duration)
				payload := &hash{
					"metrics":     metrics,
					"environment": config.Env,
					"hostname":    config.Hostname,
				}
				workerErr := worker.Push(func() error {
					if err := config.Backend.Notify(Metrics, payload); err != nil {
						return err
					}
					return nil
				})
				if workerErr != nil {
					config.Logger.Printf("worker error: %v feature=metrics\n", workerErr)
				}
			}
		}
	}()

	return &metricCollector{counter, timing}
}
