package config

type ChannelConfig struct {
	Gain   float64 `json:"gain"`
	Offset float64 `json:"offset"`
	Window int     `json:"quench_window"`
}

type ComparisonConfig struct {
	VoltThreshold float64 `json:"threshold"`     // millivolts
	DebounceTime  int     `json:"debounce_time"` // minimum consecutive samples
}

type CommonConfig struct {
	ADCResolution float64 `json:"ADC_resolution"`
	DataSizeMs    int
}

type QDSConfig struct {
	Common              CommonConfig
	Channels            map[int]ChannelConfig
	Comparison          map[int]ComparisonConfig
	ImbalanceChannels   []ImbalanceChannel
	OvervoltageChannels []OvervoltageChannel
}

func (c QDSConfig) NumChannels() int {
	return len(c.Channels)
}
func (c QDSConfig) IsChannelValid(channel int) bool {
	_, ok := c.Channels[channel]
	return ok
}

type ImbalanceChannel struct {
	Channel1   int
	Channel2   int
	Comparison int
}

type OvervoltageChannel struct {
	Channel    int
	Comparison int
}
