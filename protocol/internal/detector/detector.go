package detector

import (
	"math"

	"github.com/devicehub-go/cnpem-qds/protocol/config"
	"github.com/devicehub-go/cnpem-qds/protocol/internal/decoder"
)

type QuenchResult struct {
	IsQuench  bool
	QuenchIdx int // -1 indicates no quench
}

type Result struct {
	Channels       map[int]QuenchResult
	FirstQuenchIdx int
	QuenchCount    int
}

// Checks if a quench exists in the provided data and where it occurs
func checkThreshold(voltage []float64, cfg config.ComparisonConfig) QuenchResult {
	result := QuenchResult{QuenchIdx: -1}
	if cfg.VoltThreshold <= 0 {
		return result
	}

	threshold := cfg.VoltThreshold * 1e-3
	count := 0

	for i, v := range voltage {
		if math.Abs(v) >= threshold {
			count++
		} else {
			count = 0
		}
		if count >= cfg.DebounceTime {
			result.IsQuench = true
			result.QuenchIdx = i - count + 1
			break
		}
	}
	return result
}

// Updates the quench index according to the current frame result
func (r *Result) accumulate(result QuenchResult) {
	if !result.IsQuench {
		return
	}
	r.QuenchCount++
	if r.FirstQuenchIdx == -1 || result.QuenchIdx < r.FirstQuenchIdx {
		r.FirstQuenchIdx = result.QuenchIdx
	}
}

// Indicates that exists an quench in the result
func (r *Result) HasQuench() bool {
	return r.QuenchCount > 0
}

// Detects imbalance and overvoltage quenches
func Detect(frame decoder.Frame, cfg config.QDSConfig) Result {
	result := Result{
		FirstQuenchIdx: -1,
		Channels:       make(map[int]QuenchResult),
	}

	for _, pair := range cfg.ImbalanceChannels {
		comparison := cfg.Comparison[pair.Comparison]
		v1 := frame.Channels[pair.Channel1].Voltage
		v2 := frame.Channels[pair.Channel2].Voltage

		imbalance := make([]float64, len(v1))
		for i := range v1 {
			imbalance[i] = math.Abs(math.Abs(v1[i]) - math.Abs(v2[i]))
		}

		check := checkThreshold(imbalance, comparison)
		result.Channels[pair.Channel1] = check
		result.Channels[pair.Channel2] = check
		result.accumulate(check)
	}

	for _, single := range cfg.OvervoltageChannels {
		comparison := cfg.Comparison[single.Comparison]
		voltage := frame.Channels[single.Channel].Voltage

		check := checkThreshold(voltage, comparison)
		result.Channels[single.Channel] = check
		result.accumulate(check)
	}

	return result
}
