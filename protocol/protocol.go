package protocol

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/devicehub-go/cnpem-qds/protocol/config"
	"github.com/devicehub-go/cnpem-qds/protocol/internal/buffer"
	"github.com/devicehub-go/cnpem-qds/protocol/internal/decoder"
	"github.com/devicehub-go/cnpem-qds/protocol/internal/queue"
	"github.com/devicehub-go/cnpem-qds/protocol/internal/watchdog"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type Options struct {
	Host          string
	Port          int
	Timeout       time.Duration
	TopicPrefix   string
	Configuration config.QDSConfig
}

type Middleware struct {
	mutex   sync.Mutex
	options Options
	client  paho.Client

	config  config.QDSConfig
	buffer  *buffer.Buffer
	decoder *decoder.Decoder

	queue *queue.Queue

	FrameChan    chan FrameEvent
	SnapshotChan chan FrameEvent
}

// Creates a new middleware instance, that provides methods
// to interface communication with QDS through MQTT
func New(options Options, config config.QDSConfig) *Middleware {
	url := fmt.Sprintf("tcp://%s:%d", options.Host, options.Port)
	if options.Timeout == 0 {
		options.Timeout = 100 * time.Millisecond
	}

	clientOptions := paho.NewClientOptions()
	clientOptions.AddBroker(url)

	prefix := strings.ReplaceAll(options.TopicPrefix, "/", "_")
	clientId := fmt.Sprintf("qds_processing_go_%s", prefix)
	clientOptions.SetClientID(clientId)

	return &Middleware{
		options:      options,
		client:       paho.NewClient(clientOptions),
		buffer:       buffer.New(config),
		decoder:      decoder.New(config),
		queue:        queue.New(),
		FrameChan:    make(chan FrameEvent),
		SnapshotChan: make(chan FrameEvent),
	}
}

// Establishes a connection with MQTT broker
func (m *Middleware) Connect() error {
	var topicsCallback = map[string]paho.MessageHandler{
		m.options.TopicPrefix + "/quench_data": m.onQDSData,
	}

	token := m.client.Connect()
	if !token.WaitTimeout(m.options.Timeout) {
		return errors.New("MQTT broker connection timed out")
	}
	if err := token.Error(); err != nil {
		return err
	}
	for topic, callback := range topicsCallback {
		log.Printf("[DEBUG] Subscribing to %s\n", topic)
		token := m.client.Subscribe(topic, 0, callback)
		if token.Error() != nil {
			return token.Error()
		}
	}

	topic := m.options.TopicPrefix + "/qds-processing-connection"
	token = m.client.Publish(topic, 2, true, "1")
	return token.Error()
}

// Closes the connection with MQTT Broker
func (m *Middleware) Close() error {
	m.client.Disconnect(100)
	return nil
}

// Starts to process data queue messages
func (m *Middleware) Run(ctx context.Context) error {
	watchdog := watchdog.New(ctx, 10*time.Second)
	watchdog.Start()

	for {
		select {
		case <-ctx.Done():
			return nil

		case err := <-watchdog.Err():
			return err

		case <-m.queue.Notify:
			watchdog.Kick()
			for {
				fmt.Println(time.Now())
				payload, ok := m.queue.Dequeue()
				if !ok {
					break
				}
				m.onProcessing(payload)
			}
		}
	}
}
