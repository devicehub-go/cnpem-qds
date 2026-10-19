package swlsqds_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	swlsqds "github.com/devicehub-go/cnpem-qds"
	"github.com/devicehub-go/cnpem-qds/protocol/config"
)

func TestMain(t *testing.T) {
	options := swlsqds.Options{
		Host:            "dat-aro.cnpem.br",
		Port:            1883,
		Timeout:         1 * time.Second,
		TopicPrefix:     "weg/qds/r",
		WatchdogTimeout: 10 * time.Hour,
	}

	numChannels := 8

	channels := make(map[int]config.ChannelConfig, numChannels)
	comparison := make(map[int]config.ComparisonConfig, numChannels)
	overvoltage := make([]config.OvervoltageChannel, numChannels)

	for i := range numChannels {
		channels[i] = config.ChannelConfig{
			Gain:   1.0,
			Offset: 0.0,
			Window: 200,
		}
		overvoltage[i] = config.OvervoltageChannel{
			Channel:    i,
			Comparison: i,
		}
		comparison[i] = config.ComparisonConfig{
			VoltThreshold: 250,
			DebounceTime:  10,
		}
	}

	config := swlsqds.Config{
		Common: config.CommonConfig{
			ADCResolution: 152,
			DataSizeMs:    200,
		},
		Channels:            channels,
		Comparison:          comparison,
		OvervoltageChannels: overvoltage,
	}

	qds := swlsqds.New(options, config)
	if err := qds.Connect(); err != nil {
		log.Fatalf("failed on connect to QDS: %v", err)
	}

	ctx := context.Background()
	go qds.Run(ctx)

	for {
		fmt.Println(qds.GetDebounceTime(0))
		time.Sleep(1 * time.Second)
	}
}
