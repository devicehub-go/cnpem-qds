package protocol

import "fmt"

func (m *Middleware) SetGain(channel int, gain float64) error {
	config, ok := m.config.Channels[channel]
	if !ok {
		return fmt.Errorf("invalid channel, got: %d", channel)
	}
	config.Gain = gain
	m.config.Channels[channel] = config
	return nil
}

func (m *Middleware) GetGain(channel int) (float64, error) {
	config, ok := m.config.Channels[channel]
	if !ok {
		return 0.0, fmt.Errorf("invalid channel, got: %d", channel)
	}
	return config.Gain, nil
}

func (m *Middleware) SetOffset(channel int, offset float64) error {
	config, ok := m.config.Channels[channel]
	if !ok {
		return fmt.Errorf("invalid channel, got: %d", channel)
	}
	config.Offset = offset
	m.config.Channels[channel] = config
	return nil
}

func (m *Middleware) GetOffset(channel int) (float64, error) {
	config, ok := m.config.Channels[channel]
	if !ok {
		return 0.0, fmt.Errorf("invalid channel, got: %d", channel)
	}
	return config.Offset, nil
}
