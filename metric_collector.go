package honeybadger

import (
	"fmt"
	"time"
)

type counterMetric struct {
	name  string
	value int
}

type metricCollector struct {
	counter chan counterMetric
}

func (c *metricCollector) increment(metric string, value int) {
	c.counter <- counterMetric{name: metric, value: value}
}

func newMetricCollector(config *Configuration, worker worker) *metricCollector {
	counter := make(chan counterMetric)
	go func() {
		flush := make(chan int)
		go func() {
			for {
				time.Sleep(config.MetricsInterval)
				flush <- 1
			}
		}()

		counters := make(map[string]int)
		for {
			select {
			case metric := <-counter:
				counters[metric.name] += metric.value
			case <-flush:
				if len(counters) <= 0 {
					break
				}
				metrics := []string{}
				for name, value := range counters {
					metrics = append(metrics, fmt.Sprintf("%v %v", name, value))
				}
				counters = make(map[string]int)
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

	return &metricCollector{counter: counter}
}
