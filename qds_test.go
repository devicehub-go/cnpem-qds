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
	config := swlsqds.Config{
		Common: config.CommonConfig{
			ADCResolution: 152,
			DataSizeMs:    200,
		},
		Channels: map[int]config.ChannelConfig{
			0: {
				Gain:   1.0,
				Offset: 0.0,
				Window: 200,
			},
			1: {
				Gain:   1.0,
				Offset: 0.0,
				Window: 200,
			},
			2: {
				Gain:   1.0,
				Offset: 0.0,
				Window: 200,
			},
			3: {
				Gain:   1.0,
				Offset: 0.0,
				Window: 200,
			},
		},
		Comparison: map[int]config.ComparisonConfig{
			0: {
				VoltThreshold: 250,
				DebounceTime:  10,
			},
			1: {
				VoltThreshold: 250,
				DebounceTime:  10,
			},
			2: {
				VoltThreshold: 250,
				DebounceTime:  10,
			},
			3: {
				VoltThreshold: 250,
				DebounceTime:  10,
			},
		},
		OvervoltageChannels: []config.OvervoltageChannel{
			{Channel: 0, Comparison: 0},
			{Channel: 1, Comparison: 1},
			{Channel: 2, Comparison: 2},
			{Channel: 3, Comparison: 3},
		},
	}

	qds := swlsqds.New(options, config)
	if err := qds.Connect(); err != nil {
		log.Fatalf("failed on connect to QDS: %v", err)
	}

	ctx := context.Background()
	go qds.Run(ctx)

	fmt.Println(qds.GetGain(0))

	for {
		fmt.Println(qds.CurrentFrame)
		time.Sleep(1 * time.Second)
	}
}
