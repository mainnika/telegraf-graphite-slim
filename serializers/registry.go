package serializers

import (
	"fmt"
	"time"

	"github.com/influxdata/telegraf"

	"github.com/influxdata/telegraf/serializers/graphite"
)

// SerializerOutput is an interface for output plugins that are able to
// serialize telegraf metrics into arbitrary data formats.
type SerializerOutput interface {
	// SetSerializer sets the serializer function for the interface.
	SetSerializer(serializer Serializer)
}

// Serializer is an interface defining functions that a serializer plugin must
// satisfy.
type Serializer interface {
	// Serialize takes a single telegraf metric and turns it into a byte buffer.
	// separate metrics should be separated by a newline, and there should be
	// a newline at the end of the buffer.
	Serialize(metric telegraf.Metric) ([]byte, error)

	// SerializeBatch takes an array of telegraf metric and serializes it into
	// a byte buffer.  This method is not required to be suitable for use with
	// line oriented framing.
	SerializeBatch(metrics []telegraf.Metric) ([]byte, error)
}

// Config is a struct that covers the data types needed for all serializer types,
// and can be used to instantiate _any_ of the serializers.
type Config struct {
	// Dataformat can be one of: influx, graphite, or json
	DataFormat string

	// Support tags in graphite protocol
	GraphiteTagSupport bool

	// Maximum line length in bytes; influx format only
	InfluxMaxLineBytes int

	// Sort field keys, set to true only when debugging as it less performant
	// than unsorted fields; influx format only
	InfluxSortFields bool

	// Support unsigned integer output; influx format only
	InfluxUintSupport bool

	// Prefix to add to all measurements, only supports Graphite
	Prefix string

	// Template for converting telegraf metrics into Graphite
	// only supports Graphite
	Template string

	// Timestamp units to use for JSON formatted output
	TimestampUnits time.Duration

	// Include HEC routing fields for splunkmetric output
	HecRouting bool
}

// NewSerializer a Serializer interface based on the given config.
func NewSerializer(config *Config) (Serializer, error) {
	var err error
	var serializer Serializer
	switch config.DataFormat {
	case "graphite":
		serializer, err = NewGraphiteSerializer(config.Prefix, config.Template, config.GraphiteTagSupport)
	default:
		err = fmt.Errorf("Invalid data format: %s", config.DataFormat)
	}
	return serializer, err
}

func NewGraphiteSerializer(prefix, template string, tag_support bool) (Serializer, error) {
	return &graphite.GraphiteSerializer{
		Prefix:     prefix,
		Template:   template,
		TagSupport: tag_support,
	}, nil
}
