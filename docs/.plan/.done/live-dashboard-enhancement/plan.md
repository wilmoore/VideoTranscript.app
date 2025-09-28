# Live Dashboard Enhancement - Implementation Plan

## Objective
Enhance the web dashboard with real-time Server-Sent Events (SSE) for live statistics and metrics, implementing full live reload functionality, and improving the overall user experience.

## Current State Analysis
- Multiple web dashboard instances running on different ports (8769+)
- Web dashboard basic functionality present in `web-dashboard.go`
- Encore.dev auth handler issues resolved ✅
- Database persistence implemented with PostgreSQL
- Legacy file-based artifacts cleaned up ✅
- Comprehensive documentation structure completed ✅

## Implementation Steps

### 1. Implement Server-Sent Events (SSE) for Live Data
- **File**: `web-dashboard.go`
- Add SSE endpoint `/events` for real-time dashboard updates
- Stream live job status changes, processing metrics, and system stats
- Implement JSON event format for structured data updates

### 2. Enhance Dashboard Frontend with Live Updates
- **File**: HTML template in `web-dashboard.go`
- Add JavaScript EventSource for SSE connection
- Implement real-time DOM updates for job statistics
- Add live charts for processing metrics (jobs/minute, success rate, etc.)
- Remove manual refresh requirements

### 3. Implement Live Reload Development Features
- Add file watching capabilities for Go source files
- Automatic server restart on code changes
- Port management to avoid constant port switching
- Hot reload for frontend assets

### 4. Database Metrics Collection
- **Files**: `transcribe/service.go`, `transcribe/db.go`
- Implement real-time job metrics aggregation
- Add system performance tracking (processing times, queue depth)
- Create metrics endpoints for dashboard consumption

### 5. Clean Up Multiple Dashboard Instances
- Investigate and resolve multiple running dashboard processes
- Implement proper process management
- Add graceful shutdown handling

## Technical Requirements

### SSE Implementation
```go
// Add to web-dashboard.go
func handleSSE(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    // Stream metrics every 2 seconds
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            metrics := collectDashboardMetrics()
            fmt.Fprintf(w, "data: %s\n\n", metrics)
            w.(http.Flusher).Flush()
        case <-r.Context().Done():
            return
        }
    }
}
```

### Frontend JavaScript Updates
```javascript
// Add to dashboard template
const eventSource = new EventSource('/events');
eventSource.onmessage = function(event) {
    const metrics = JSON.parse(event.data);
    updateDashboardMetrics(metrics);
};
```

### Live Reload Implementation
```go
// File watcher for development mode
func watchFiles() {
    watcher, _ := fsnotify.NewWatcher()
    defer watcher.Close()

    go func() {
        for {
            select {
            case event := <-watcher.Events:
                if event.Op&fsnotify.Write == fsnotify.Write {
                    restartServer()
                }
            }
        }
    }()
}
```

## Success Criteria
- ✅ Real-time dashboard updates without manual refresh
- ✅ Live job statistics and processing metrics
- ✅ Single dashboard instance running on consistent port
- ✅ Development live reload functionality
- ✅ Clean process management and graceful shutdowns
- ✅ SSE connection stability and error handling

## Dependencies
- Database connection for metrics queries
- File system watching capabilities (`fsnotify` package)
- HTTP streaming support for SSE
- JavaScript EventSource API

## Implementation Priority
1. **High Priority**: SSE implementation and real-time metrics
2. **Medium Priority**: Live reload development features
3. **Low Priority**: Advanced dashboard visualizations

## Notes
- This builds on the completed documentation and legacy cleanup work
- Addresses the user's original request for "live stats and metrics via SSE"
- Improves development workflow with proper live reload
- Resolves multiple dashboard instance issues