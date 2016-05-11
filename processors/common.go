package processors

import (
	"github.com/getlantern/golog"

	events "github.com/getlantern/events-pipeline"
)

type keyCountMap map[events.Key]int64

var (
	log = golog.LoggerFor("processors")
)
