package outputs

import (
	"github.com/mainnika/telegraf-graphite-slim"
)

type Creator func() telegraf.Output

var Outputs = map[string]Creator{}

func Add(name string, creator Creator) {
	Outputs[name] = creator
}
