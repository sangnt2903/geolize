package contexts

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

type serverCtx struct{}
type serverData struct {
	// Headers
	IncomingHeaders map[string][]string
	OutgoingHeaders map[string][]string
	InternalHeaders map[string][]string

	// Request tracking
	RequestID string
}

// Add helper functions to get/set data
func GetServerData(ctx context.Context) *serverData {
	if d := ctx.Value(serverCtx{}); d != nil {
		return d.(*serverData)
	}
	return nil
}

func SetIncomingHeader(ctx context.Context, key string, values []string) {
	if d := GetServerData(ctx); d != nil {
		if d.IncomingHeaders == nil {
			d.IncomingHeaders = make(map[string][]string)
		}
		d.IncomingHeaders[key] = values
	}
}

func SetOutgoingHeader(ctx context.Context, key string, values []string) {
	if d := GetServerData(ctx); d != nil {
		if d.OutgoingHeaders == nil {
			d.OutgoingHeaders = make(map[string][]string)
		}
		d.OutgoingHeaders[key] = values
	}
}

func SetInternalHeader(ctx context.Context, key string, values []string) {
	if d := GetServerData(ctx); d != nil {
		if d.InternalHeaders == nil {
			d.InternalHeaders = make(map[string][]string)
		}
		d.InternalHeaders[key] = values
	}
}

func SetRequestID(ctx context.Context, requestID string) {
	if d := GetServerData(ctx); d != nil {
		d.RequestID = requestID
	}
}

func NewServerContext(ctx context.Context) context.Context {
	d := ctx.Value(serverCtx{})
	if d != nil {
		return ctx
	}

	// loading header from metadata
	serverData := &serverData{
		IncomingHeaders: make(map[string][]string),
		OutgoingHeaders: make(map[string][]string),
		InternalHeaders: make(map[string][]string),
	}

	// Get incoming metadata from context
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for key, values := range md {
			serverData.IncomingHeaders[key] = values
		}
	}

	// check robot request
	if userAgents := serverData.IncomingHeaders["user-agent"]; len(userAgents) > 0 {
		userAgent := userAgents[0]
		isRobot := false
		robotPatterns := []string{
			"bot", "crawler", "spider", "ping", "slurp",
			"google", "baidu", "bing", "yahoo",
		}

		userAgentLower := strings.ToLower(userAgent)
		for _, pattern := range robotPatterns {
			if strings.Contains(userAgentLower, pattern) {
				isRobot = true
				break
			}
		}

		if isRobot {
			serverData.InternalHeaders["is-robot"] = []string{"true"}
		}
	}

	// Initialize request tracking IDs from headers if present
	if requestIDs := serverData.IncomingHeaders["x-request-id"]; len(requestIDs) > 0 {
		serverData.RequestID = requestIDs[0]
	}

	return context.WithValue(ctx, serverCtx{}, serverData)
}

func IsRobot(ctx context.Context) bool {
	if d := GetServerData(ctx); d != nil {
		if values, exists := d.InternalHeaders["is-robot"]; exists && len(values) > 0 {
			return values[0] == "true"
		}
	}
	return false
}
