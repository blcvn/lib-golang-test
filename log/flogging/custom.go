package flogging

import (
	"time"

	"go.uber.org/zap"
)

func currentMillis() int64 {
	now := time.Now()
	return now.UnixNano() / int64(time.Millisecond)
}

// StartFunction must be called at the begin of function
// func (f *FabricLogger) StartFunction(traceNo string) {
// 	f.s.Infof("[%s] StartFunction at %d", traceNo, currentMillis())
// }

func (f *FabricLogger) GetRootLogger() *zap.SugaredLogger {
	return f.s
}
