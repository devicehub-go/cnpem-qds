package buffer

import (
	"fmt"

	"github.com/devicehub-go/cnpem-qds/protocol/internal/decoder"
	"github.com/devicehub-go/cnpem-qds/protocol/internal/detector"
	"github.com/devicehub-go/cnpem-qds/protocol/internal/utils"

	"github.com/devicehub-go/cnpem-qds/protocol/config"
)

type ChannelSnapshot struct {
	Voltage   []float64
	QuenchIdx int
}

type Buffer struct {
	cfg        config.QDSConfig
	samples    [][]float64
	quenchIdx  int
	dispatched bool
}

// Creates a new instance of quench buffer
func New(cfg config.QDSConfig) *Buffer {
	return &Buffer{
		cfg:       cfg,
		quenchIdx: -1,
		samples:   make([][]float64, cfg.NumChannels()),
	}
}

// Replace the current configuration
func (b *Buffer) UpdateConfig(cfg config.QDSConfig) {
	b.cfg = cfg
}

// Calculates the size of a half window
func (b *Buffer) halfWindow() int {
	max := 0
	for _, ch := range b.cfg.Channels {
		if h := ch.Window / 2; h > max {
			max = h
		}
	}
	return utils.RoundMultiple(max, b.cfg.Common.DataSizeMs)
}

// Clears the buffer and reset returns the state to idle
func (b *Buffer) reset() {
	b.quenchIdx = -1
	b.dispatched = false
}

func (b *Buffer) Feed(frame decoder.Frame, result detector.Result) bool {
	if b.dispatched && !result.HasQuench() {
		b.reset()
	}

	halfWindow := b.halfWindow()
	fullWindow := halfWindow * 2

	for ch := range b.cfg.NumChannels() {
		voltage := frame.Channels[ch].Voltage
		b.samples[ch] = append(b.samples[ch], voltage...)

		if b.quenchIdx == -1 && len(b.samples[ch]) > fullWindow {
			b.samples[ch] = b.samples[ch][len(b.samples[ch])-fullWindow:]
		}
	}

	if b.quenchIdx == -1 && result.HasQuench() && result.FirstQuenchIdx != -1 {
		b.quenchIdx = len(b.samples[0]) - len(frame.Channels[0].Voltage) + result.FirstQuenchIdx
	}

	if b.quenchIdx != -1 && !b.dispatched {
		if len(b.samples[0]) >= b.quenchIdx+halfWindow {
			return true
		}
	}

	return false
}

// Slices and centers the buffer around the quench index
func (b *Buffer) Dispatch() ([]ChannelSnapshot, error) {
	snapshot := make([]ChannelSnapshot, b.cfg.NumChannels())

	for ch := range b.cfg.NumChannels() {
		window := b.cfg.Channels[ch].Window
		halfWindow := window / 2
		index := halfWindow

		start := b.quenchIdx - halfWindow
		end := b.quenchIdx + halfWindow

		if start < 0 {
			index = b.quenchIdx
			start = 0
		}

		if end > len(b.samples[ch]) {
			return snapshot, fmt.Errorf("missing data in snapshot")
		}

		slice := b.samples[ch][start:end]
		out := make([]float64, len(slice))
		copy(out, slice)

		snapshot[ch] = ChannelSnapshot{
			Voltage:   out,
			QuenchIdx: index,
		}
	}

	b.dispatched = true
	return snapshot, nil
}

// Reports whether the current event has already been dispatched
func (b *Buffer) IsDispatched() bool {
	return b.dispatched
}
