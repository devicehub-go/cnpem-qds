package protocol

import (
	"fmt"
	"log"
	"time"

	"github.com/devicehub-go/cnpem-qds/protocol/internal/decoder"
	"github.com/devicehub-go/cnpem-qds/protocol/internal/detector"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type FrameEvent struct {
	Timestamp time.Time
	Channels  map[int]ChannelData
}

type ChannelData struct {
	Voltage   []float64
	IsQuench  bool
	QuenchIdx int
}

// Handles quench data messages
func (m *Middleware) onQDSData(client paho.Client, message paho.Message) {
	fmt.Println(message.Topic())
	m.queue.Enqueue(message.Payload())
}

// onProcessing is the core pipeline: decode → detect → buffer → publish
func (m *Middleware) onProcessing(payload []byte) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	frame, err := m.decoder.DecodeFrame(payload)
	if err != nil {
		log.Printf("Invalid quench data: %+v", err)
		return
	}

	result := detector.Detect(frame, m.config)

	go func(f decoder.Frame, r detector.Result) {
		out := FrameEvent{
			Timestamp: frame.Timestamp,
			Channels:  make(map[int]ChannelData),
		}
		for ch, chFrame := range frame.Channels {
			v := make([]float64, len(chFrame.Voltage))
			copy(v, chFrame.Voltage)
			out.Channels[ch] = ChannelData{
				Voltage:   v,
				IsQuench:  result.Channels[ch].IsQuench,
				QuenchIdx: result.Channels[ch].QuenchIdx,
			}
		}
		m.FrameChan <- out
	}(frame, result)

	if m.buffer.Feed(frame, result) {
		snap, err := m.buffer.Dispatch()
		if err != nil {
			return
		}
		out := FrameEvent{
			Timestamp: time.Now(),
			Channels:  make(map[int]ChannelData),
		}
		for ch, s := range snap {
			v := make([]float64, len(s.Voltage))
			copy(v, s.Voltage)
			out.Channels[ch] = ChannelData{
				Voltage:   v,
				IsQuench:  true,
				QuenchIdx: s.QuenchIdx,
			}
		}
		m.SnapshotChan <- out
	}
}
