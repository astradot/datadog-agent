package aggregator

import (
	"github.com/DataDog/datadog-agent/pkg/metrics"
	"github.com/DataDog/datadog-agent/pkg/metrics/percentile"
)

// FIXME(Jee) : This should be integrated with time_sampler.go since it
// duplicates code/logic.

// DistSampler creates sketches of metrics by buckets of 'interval' seconds
type DistSampler struct {
	interval            int64
	contextResolver     *ContextResolver
	sketchesByTimestamp map[int64]ContextSketch
	defaultHostname     string
}

// NewDistSampler returns a newly initialized DistSampler
func NewDistSampler(interval int64, defaultHostname string) *DistSampler {
	return &DistSampler{
		interval:            interval,
		contextResolver:     newContextResolver(),
		sketchesByTimestamp: map[int64]ContextSketch{},
		defaultHostname:     defaultHostname,
	}
}

func (d *DistSampler) calculateBucketStart(timestamp float64) int64 {
	return int64(timestamp) - int64(timestamp)%d.interval
}

// Add the metricSample to the correct sketch
func (d *DistSampler) addSample(metricSample *metrics.MetricSample, timestamp float64) {
	contextKey := d.contextResolver.trackContext(metricSample, timestamp)
	bucketStart := d.calculateBucketStart(timestamp)
	sketch, ok := d.sketchesByTimestamp[bucketStart]
	if !ok {
		sketch = makeContextSketch()
		d.sketchesByTimestamp[bucketStart] = sketch
	}
	sketch.addSample(contextKey, metricSample, timestamp, d.interval)
}

func (d *DistSampler) flush(timestamp float64) []*percentile.SketchSeries {
	var result []*percentile.SketchSeries

	sketchesByContext := make(map[string]*percentile.SketchSeries)

	cutoffTime := d.calculateBucketStart(timestamp)
	for bucketTimestamp, ctxSketch := range d.sketchesByTimestamp {
		if cutoffTime <= bucketTimestamp {
			continue
		}

		sketches := ctxSketch.flush(float64(bucketTimestamp))
		for _, sketchSeries := range sketches {
			contextKey := sketchSeries.ContextKey

			if existingSeries, ok := sketchesByContext[contextKey]; ok {
				existingSeries.Sketches = append(existingSeries.Sketches, sketchSeries.Sketches[0])
			} else {
				context := d.contextResolver.contextsByKey[contextKey]
				sketchSeries.Name = context.Name
				sketchSeries.Tags = context.Tags
				if context.Host != "" {
					sketchSeries.Host = context.Host
				} else {
					sketchSeries.Host = d.defaultHostname
				}
				sketchSeries.Interval = d.interval

				sketchesByContext[contextKey] = sketchSeries
				result = append(result, sketchSeries)
			}
		}
		delete(d.sketchesByTimestamp, bucketTimestamp)
	}
	d.contextResolver.expireContexts(timestamp - defaultExpiry)

	return result
}
