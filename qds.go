package swlsqds

import (
	"github.com/devicehub-go/cnpem-qds/protocol"
	"github.com/devicehub-go/cnpem-qds/protocol/config"
)

// Creates a new middleware instance, that provides methods
// to interface communication with QDS through MQTT
func New(options protocol.Options, config config.QDSConfig) *protocol.Middleware {
	return protocol.New(options, config)
}
