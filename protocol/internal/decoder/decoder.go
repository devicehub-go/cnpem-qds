package decoder

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/devicehub-go/cnpem-qds/protocol/config"
)

type ChannelFrame struct {
	Voltage []float64
}

type Frame struct {
	Timestamp time.Time
	Channels  []ChannelFrame
}

type Decoder struct {
	cfg config.QDSConfig
}

func New(config config.QDSConfig) *Decoder {
	return &Decoder{cfg: config}
}

// Replace the current configuration
func (d *Decoder) UpdateConfig(cfg config.QDSConfig) {
	d.cfg = cfg
}

// Returns calibrated data based on gain and offset of the channel
func calibrate(raw []float64, adcResolution, gain, offset float64) []float64 {
	calibrated := make([]float64, len(raw))
	adcRes := adcResolution * 1e-6
	for i, v := range raw {
		calibrated[i] = (v*adcRes)*gain + offset
	}
	return calibrated
}

// Extracts channel from string of decoded quench data
func (d *Decoder) channelFromKey(key string) (int, error) {
	parts := strings.Split(key, "_")
	if len(parts) < 1 {
		return 0, fmt.Errorf("decoder: invalid key %q", key)
	}
	channel, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return 0, fmt.Errorf("decoder: non-numeric channel in key %q: %w", key, err)
	}
	if !d.cfg.IsChannelValid(channel) {
		return 0, fmt.Errorf("decoder: channel '%d' is not registered", channel)
	}
	return channel, nil
}

// Decodes the payload for quench data message
func (d *Decoder) DecodeFrame(payload []byte) (Frame, error) {
	var raw map[string]any

	if err := json.Unmarshal(payload, &raw); err != nil {
		return Frame{}, fmt.Errorf("decoder: unmarshal failed: %w", err)
	}
	frame := Frame{
		Timestamp: time.Now(),
		Channels:  make([]ChannelFrame, d.cfg.NumChannels()),
	}
	for key, value := range raw {
		if !strings.HasPrefix(key, "channel_data_") {
			continue
		}

		ch, err := d.channelFromKey(key)
		if err != nil {
			return Frame{}, err
		}

		encoded, ok := value.(string)
		if !ok {
			return Frame{}, fmt.Errorf("decoder: channel %d value is not a string", ch)
		}

		raw := Decode(DecodeOptions{Encoded: encoded})
		if len(raw) != d.cfg.Common.DataSizeMs {
			log.Printf(
				"[WARNING] decoder: channel %d length %d, want %d",
				ch, len(raw), d.cfg.Common.DataSizeMs,
			)
		}

		cfg := d.cfg.Channels[ch]
		adc := d.cfg.Common.ADCResolution
		frame.Channels[ch].Voltage = calibrate(raw, adc, cfg.Gain, cfg.Offset)
	}

	return frame, nil
}
