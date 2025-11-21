package logs

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type LogEntry struct {
	Time     time.Time
	Function string
	Message  string
}

type CtxLog struct {
	sync.Mutex
	Data []LogEntry
}

type key int

const LogsKey key = 1

func logContext(ctx context.Context, path string, start time.Time) {
	logs, ok := ctx.Value(LogsKey).(*CtxLog)
	if !ok || logs == nil {
		return
	}

	duration := time.Since(start)
	logs.Lock()
	entries := make([]LogEntry, len(logs.Data))
	copy(entries, logs.Data)
	logs.Unlock()

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s\n", path)
	fmt.Fprintf(&buf, "Request duration: %v\n", duration)

	for _, e := range entries {
		fmt.Fprintf(
			&buf,
			"\t[%s] %s: %s\n",
			e.Time.Format(time.RFC3339),
			e.Function,
			e.Message,
		)
	}

	fmt.Println(buf.String())
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		start := time.Now()

		ctx = context.WithValue(ctx, LogsKey, &CtxLog{
			Data: make([]LogEntry, 0, 4), // небольшой запас, чтобы меньше реаллоцировать
		})

		defer logContext(ctx, r.URL.Path, start)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//TODO: реализовать контекст с логами для отдельных горутин, не внтури request.Context
//TODO: может быть реализовать перенаправление логов кодом
