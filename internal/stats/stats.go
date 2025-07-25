package stats

import (
	"slices"
	"time"
)

type Statistics struct {
	Mean  time.Duration
	Min   time.Duration
	Max   time.Duration
	Range time.Duration
}

func CalculateStatistics(durations []time.Duration) Statistics {
	if len(durations) == 0 {
		return Statistics{}
	}

	slices.Sort(durations)

	var total time.Duration
	for _, d := range durations {
		total += d
	}

	mean := total / time.Duration(len(durations))
	min := durations[0]
	max := durations[len(durations)-1]
	rang := max - min

	return Statistics{
		Mean:  mean,
		Min:   min,
		Max:   max,
		Range: rang,
	}
}
