package instrumentation

import (
	"github.com/davecheney/gmx"
)

var TimerFail *gmx.Counter

func init() {
	TimerFail = gmx.NewCounter("timerFail")
}
