package commands

import (
	"fmt"
	"time"

	"github.com/montanaflynn/stats"
)

var benchmark = app.Command("benchmark", "Repeat build for ten seconds. Implies --profile.")

// benchmarkCommand builds the site repeatedly until at least 10 seconds has elapsed,
// and reports the trial times. Empirically, it the same mean but low variance as using
// a separate benchmark runner that invokes a new gojekyll process each time.
func benchmarkCommand() (err error) {
	startTime := time.Now()
	samples := []float64{}
	for i := 0; time.Since(startTime) < 10*time.Second; i++ {
		sampleStart := time.Now()
		site, err := loadSite(*source, options)
		if err != nil {
			return err
		}
		_, err = site.Build()
		if err != nil {
			return err
		}
		dur := time.Since(sampleStart).Seconds()
		samples = append(samples, dur)
		quiet = true
		fmt.Printf("Run #%d; %.1fs elapsed\n", i+1, time.Since(commandStartTime).Seconds())
	}
	median, _ := stats.Median(samples)
	stddev, _ := stats.StandardDeviationSample(samples)
	fmt.Printf("%d samples @ %.2fs ± %.2fs\n", len(samples), median, stddev)
	return nil
}
