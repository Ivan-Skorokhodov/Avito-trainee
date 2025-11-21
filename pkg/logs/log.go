package logs

import (
	"context"
	"fmt"
	"time"
)

func PrintLog(ctx context.Context, funcName, message string) {
	ctxLog, ok := ctx.Value(LogsKey).(*CtxLog)
	if !ok || ctxLog == nil {
		// fallback в обычный лог:
		fmt.Printf("[no-ctx-log] %s: %s\n", funcName, message)
		return
	}

	entry := LogEntry{
		Time:     time.Now(),
		Function: funcName,
		Message:  message,
	}

	ctxLog.Lock()
	ctxLog.Data = append(ctxLog.Data, entry)
	ctxLog.Unlock()
}
