package protocol

import (
	"fmt"
)

func (m *Middleware) SetGain(channel int, gain float64) error {
	config, ok := m.config.Channels[channel]
	if !ok {
		return fmt.Errorf("error on set gain: invalid channel '%d'", channel)
	}
	config.Gain = gain
	m.config.Channels[channel] = config
	return nil
}

func (m *Middleware) GetGain(channel int) (float64, error) {
	config, ok := m.config.Channels[channel]
	if !ok {
		return 0.0, fmt.Errorf("error on get gain: invalid channel '%d'", channel)
	}
	return config.Gain, nil
}

func (m *Middleware) SetOffset(channel int, offset float64) error {
	config, ok := m.config.Channels[channel]
	if !ok {
		return fmt.Errorf("error on set offset: invalid channel '%d'", channel)
	}
	config.Offset = offset
	m.config.Channels[channel] = config
	return nil
}

func (m *Middleware) GetOffset(channel int) (float64, error) {
	config, ok := m.config.Channels[channel]
	if !ok {
		return 0.0, fmt.Errorf("error on get offset: invalid channel '%d'", channel)
	}
	return config.Offset, nil
}

func (m *Middleware) getComparisonChannel(channel int) (int, error) {
	for _, config := range m.config.OvervoltageChannels {
		if config.Channel == channel {
			return config.Comparison, nil
		}
	}
	for _, config := range m.config.ImbalanceChannels {
		if config.Channel1 == channel || config.Channel2 == channel {
			return config.Comparison, nil
		}
	}
	return 0, fmt.Errorf("QDS channel %d not found as overvoltage or imbalance channel", channel)
}

func (m *Middleware) SetThreshold(channel int, threshold float64) error {
	comparisonChannel, err := m.getComparisonChannel(channel)
	if err != nil {
		return err
	}
	config, ok := m.config.Comparison[comparisonChannel]
	if !ok {
		return fmt.Errorf("error on set threshold: invalid channel '%d'", channel)
	}
	config.VoltThreshold = threshold
	m.config.Comparison[channel] = config
	return nil
}

func (m *Middleware) GetThreshold(channel int) (float64, error) {
	comparisonChannel, err := m.getComparisonChannel(channel)
	if err != nil {
		return 0.0, err
	}
	config, ok := m.config.Comparison[comparisonChannel]
	if !ok {
		return 0.0, fmt.Errorf("error on get threshold: invalid channel '%d'", channel)
	}
	return config.VoltThreshold, nil
}

func (m *Middleware) SetDebounceTime(channel int, time int) error {
	comparisonChannel, err := m.getComparisonChannel(channel)
	if err != nil {
		return err
	}
	config, ok := m.config.Comparison[comparisonChannel]
	if !ok {
		return fmt.Errorf("error on set debounce time: invalid channel '%d'", channel)
	}
	config.DebounceTime = time
	m.config.Comparison[channel] = config
	return nil
}

func (m *Middleware) GetDebounceTime(channel int) (int, error) {
	comparisonChannel, err := m.getComparisonChannel(channel)
	if err != nil {
		return 0, err
	}
	config, ok := m.config.Comparison[comparisonChannel]
	if !ok {
		return 0, fmt.Errorf("error on get debounce time: invalid channel '%d'", channel)
	}
	return config.DebounceTime, nil
}
