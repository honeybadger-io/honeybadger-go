package honeybadger

import (
	"math"
	"sort"
	"time"
)

type durationSlice []time.Duration

func (a durationSlice) Len() int           { return len(a) }
func (a durationSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a durationSlice) Less(i, j int) bool { return int64(a[i]) < int64(a[j]) }

func mean(values []time.Duration) float64 {
	if len(values) < 1 {
		return 0.0
	}
	total := 0.0
	for _, v := range values {
		total += float64(v)
	}
	return (total / float64(len(values))) / float64(time.Millisecond)
}

func median(values []time.Duration) float64 {
	if len(values) < 1 {
		return 0.0
	}
	middle := len(values) / 2
	result := values[middle]
	if len(values)%2 == 0 {
		result = (result + values[middle-1]) / 2
	}
	return float64(result) / float64(time.Millisecond)
}

func stdDev(values []time.Duration, mean float64) float64 {
	// Standard deviation requires at least 2 values.
	if len(values) <= 1 {
		return 0.0
	}
	total := 0.0
	for _, value := range values {
		total += math.Pow((float64(value)/float64(time.Millisecond))-mean, 2)
	}
	variance := total / float64(len(values)-1)
	return math.Sqrt(variance)
}

func round(f float64) int {
	return int(f + math.Copysign(0.5, f))
}

func percentile(values []time.Duration, threshold int) float64 {
	if len(values) < 1 {
		return 0.0
	}
	if len(values) == 1 {
		return float64(values[0]) / float64(time.Millisecond)
	}
	sort.Sort(durationSlice(values))
	index := int(math.Floor((float64(threshold)/100.00)*float64(len(values)))) + 1
	values = values[0:index]
	return float64(values[len(values)-1]) / float64(time.Millisecond)
}

func min(values []time.Duration) float64 {
	if len(values) < 1 {
		return 0.0
	}
	sort.Sort(durationSlice(values))
	return float64(values[0]) / float64(time.Millisecond)
}

func max(values []time.Duration) float64 {
	if len(values) < 1 {
		return 0.0
	}
	sort.Sort(durationSlice(values))
	return float64(values[len(values)-1]) / float64(time.Millisecond)
}
