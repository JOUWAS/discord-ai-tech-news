package response

import (
	"time"
)

// BaseResponse represents the basic structure for all API responses
type BaseResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// ErrorInfo contains detailed error information
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// NewsResponse represents the response for news-related requests
type NewsResponse struct {
	BaseResponse
	News []NewsItem `json:"news,omitempty"`
	Meta *MetaInfo  `json:"meta,omitempty"`
}

// NewsItem represents a single news article in the response
type NewsItem struct {
	ID          string    `json:"id,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"published_at"`
	Source      string    `json:"source"`
	Score       int       `json:"score,omitempty"`
	Category    string    `json:"category,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	TimeAgo     string    `json:"time_ago,omitempty"`
}

// SearchResponse represents the response for search requests
type SearchResponse struct {
	BaseResponse
	Query       string     `json:"query"`
	Results     []NewsItem `json:"results,omitempty"`
	ResultCount int        `json:"result_count"`
	Meta        *MetaInfo  `json:"meta,omitempty"`
}

// BotResponse represents general bot command responses
type BotResponse struct {
	BaseResponse
	Command     string            `json:"command,omitempty"`
	User        *UserInfo         `json:"user,omitempty"`
	Channel     *ChannelInfo      `json:"channel,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	DisplayText string            `json:"display_text,omitempty"`
}

// StatusResponse represents bot status information
type StatusResponse struct {
	BaseResponse
	Status      string            `json:"status"`
	Uptime      string            `json:"uptime,omitempty"`
	Version     string            `json:"version,omitempty"`
	Services    map[string]string `json:"services,omitempty"`
	Performance *PerformanceInfo  `json:"performance,omitempty"`
}

// MetaInfo contains metadata about the response
type MetaInfo struct {
	Page         int    `json:"page,omitempty"`
	PerPage      int    `json:"per_page,omitempty"`
	Total        int    `json:"total,omitempty"`
	Source       string `json:"source,omitempty"`
	CacheHit     bool   `json:"cache_hit,omitempty"`
	ResponseTime string `json:"response_time,omitempty"`
}

// UserInfo contains information about the Discord user
type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	IsBot    bool   `json:"is_bot,omitempty"`
}

// ChannelInfo contains information about the Discord channel
type ChannelInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type,omitempty"`
}

// PerformanceInfo contains performance metrics
type PerformanceInfo struct {
	MemoryUsage   string `json:"memory_usage,omitempty"`
	ResponseTime  string `json:"response_time,omitempty"`
	ActiveUsers   int    `json:"active_users,omitempty"`
	RequestsToday int    `json:"requests_today,omitempty"`
}

// HelpResponse represents help command response
type HelpResponse struct {
	BaseResponse
	Commands []CommandInfo `json:"commands,omitempty"`
	Examples []string      `json:"examples,omitempty"`
}

// CommandInfo represents information about a bot command
type CommandInfo struct {
	Name        string   `json:"name"`
	Aliases     []string `json:"aliases,omitempty"`
	Description string   `json:"description"`
	Usage       string   `json:"usage,omitempty"`
	Examples    []string `json:"examples,omitempty"`
	Category    string   `json:"category,omitempty"`
}

// WebhookResponse represents webhook endpoint responses
type WebhookResponse struct {
	BaseResponse
	WebhookID   string                 `json:"webhook_id,omitempty"`
	EventType   string                 `json:"event_type,omitempty"`
	ProcessedAt time.Time              `json:"processed_at"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status      string                 `json:"status"`
	Timestamp   time.Time              `json:"timestamp"`
	Version     string                 `json:"version,omitempty"`
	Services    map[string]ServiceInfo `json:"services,omitempty"`
	Environment string                 `json:"environment,omitempty"`
}

// ServiceInfo represents individual service health information
type ServiceInfo struct {
	Status       string `json:"status"`
	LastChecked  string `json:"last_checked,omitempty"`
	ResponseTime string `json:"response_time,omitempty"`
	Error        string `json:"error,omitempty"`
}
