package response

import (
	"fmt"
	"strings"
	"time"

	"discord-ai-tech-news/internal/repository"
)

// Builder helps construct response objects with a fluent API
type Builder struct {
	response interface{}
}

// NewNewsResponse creates a new news response builder
func NewNewsResponse() *Builder {
	return &Builder{
		response: &NewsResponse{
			BaseResponse: BaseResponse{
				Success:   true,
				Timestamp: time.Now(),
			},
		},
	}
}

// NewSearchResponse creates a new search response builder
func NewSearchResponse(query string) *Builder {
	return &Builder{
		response: &SearchResponse{
			BaseResponse: BaseResponse{
				Success:   true,
				Timestamp: time.Now(),
			},
			Query: query,
		},
	}
}

// NewBotResponse creates a new bot response builder
func NewBotResponse(command string) *Builder {
	return &Builder{
		response: &BotResponse{
			BaseResponse: BaseResponse{
				Success:   true,
				Timestamp: time.Now(),
			},
			Command: command,
		},
	}
}

// NewStatusResponse creates a new status response builder
func NewStatusResponse() *Builder {
	return &Builder{
		response: &StatusResponse{
			BaseResponse: BaseResponse{
				Success:   true,
				Timestamp: time.Now(),
			},
		},
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code, message string) *Builder {
	return &Builder{
		response: &BaseResponse{
			Success:   false,
			Timestamp: time.Now(),
			Error: &ErrorInfo{
				Code:    code,
				Message: message,
			},
		},
	}
}

// WithNews adds news items to the response
func (b *Builder) WithNews(news []repository.News) *Builder {
	if resp, ok := b.response.(*NewsResponse); ok {
		resp.News = ConvertToNewsItems(news)
		resp.Meta = &MetaInfo{
			Total:  len(news),
			Source: "HackerNews API",
		}
	}
	return b
}

// WithSearchResults adds search results to the response
func (b *Builder) WithSearchResults(results []repository.News, count int) *Builder {
	if resp, ok := b.response.(*SearchResponse); ok {
		resp.Results = ConvertToNewsItems(results)
		resp.ResultCount = count
		resp.Meta = &MetaInfo{
			Total:   count,
			PerPage: len(results),
			Source:  "HackerNews API",
		}
	}
	return b
}

// WithMessage sets the response message
func (b *Builder) WithMessage(message string) *Builder {
	switch resp := b.response.(type) {
	case *NewsResponse:
		resp.Message = message
	case *SearchResponse:
		resp.Message = message
	case *BotResponse:
		resp.Message = message
	case *StatusResponse:
		resp.Message = message
	case *BaseResponse:
		resp.Message = message
	}
	return b
}

// WithDisplayText sets the display text for Discord
func (b *Builder) WithDisplayText(text string) *Builder {
	if resp, ok := b.response.(*BotResponse); ok {
		resp.DisplayText = text
	}
	return b
}

// WithUserInfo adds user information
func (b *Builder) WithUserInfo(id, username string, isBot bool) *Builder {
	if resp, ok := b.response.(*BotResponse); ok {
		resp.User = &UserInfo{
			ID:       id,
			Username: username,
			IsBot:    isBot,
		}
	}
	return b
}

// WithChannelInfo adds channel information
func (b *Builder) WithChannelInfo(id, name, channelType string) *Builder {
	if resp, ok := b.response.(*BotResponse); ok {
		resp.Channel = &ChannelInfo{
			ID:   id,
			Name: name,
			Type: channelType,
		}
	}
	return b
}

// WithStatus sets the status information
func (b *Builder) WithStatus(status string) *Builder {
	if resp, ok := b.response.(*StatusResponse); ok {
		resp.Status = status
	}
	return b
}

// WithServices adds service status information
func (b *Builder) WithServices(services map[string]string) *Builder {
	if resp, ok := b.response.(*StatusResponse); ok {
		resp.Services = services
	}
	return b
}

// WithMetadata adds metadata to bot response
func (b *Builder) WithMetadata(metadata map[string]string) *Builder {
	if resp, ok := b.response.(*BotResponse); ok {
		resp.Metadata = metadata
	}
	return b
}

// WithError sets error information
func (b *Builder) WithError(code, message, details string) *Builder {
	switch resp := b.response.(type) {
	case *NewsResponse:
		resp.Success = false
		resp.Error = &ErrorInfo{Code: code, Message: message, Details: details}
	case *SearchResponse:
		resp.Success = false
		resp.Error = &ErrorInfo{Code: code, Message: message, Details: details}
	case *BotResponse:
		resp.Success = false
		resp.Error = &ErrorInfo{Code: code, Message: message, Details: details}
	case *StatusResponse:
		resp.Success = false
		resp.Error = &ErrorInfo{Code: code, Message: message, Details: details}
	case *BaseResponse:
		resp.Success = false
		resp.Error = &ErrorInfo{Code: code, Message: message, Details: details}
	}
	return b
}

// Build returns the constructed response
func (b *Builder) Build() interface{} {
	return b.response
}

// ConvertToNewsItems converts repository.News to response.NewsItem
func ConvertToNewsItems(news []repository.News) []NewsItem {
	items := make([]NewsItem, len(news))
	for i, article := range news {
		items[i] = NewsItem{
			ID:          fmt.Sprintf("news_%d", i+1),
			Title:       article.Title,
			Description: article.Description,
			URL:         article.URL,
			PublishedAt: article.PublishedAt,
			Source:      article.Source,
			Score:       0, // Default score since repository.News doesn't have Score field
			Category:    "Technology",
			Tags:        extractTags(article.Title + " " + article.Description),
			TimeAgo:     TimeAgo(article.PublishedAt),
		}
	}
	return items
}

// extractTags extracts relevant tags from content
func extractTags(content string) []string {
	techKeywords := []string{
		"ai", "artificial intelligence", "machine learning", "ml",
		"blockchain", "crypto", "bitcoin", "web3",
		"startup", "tech", "technology", "software",
		"programming", "developer", "coding",
		"cloud", "aws", "google", "microsoft",
		"mobile", "app", "ios", "android",
		"security", "cybersecurity", "privacy",
		"quantum", "robotics", "iot", "5g",
	}

	var tags []string
	contentLower := strings.ToLower(content)

	for _, keyword := range techKeywords {
		if strings.Contains(contentLower, keyword) {
			tags = append(tags, keyword)
		}
	}

	// Limit to 5 tags
	if len(tags) > 5 {
		tags = tags[:5]
	}

	return tags
}

// TimeAgo returns a human-readable time difference
func TimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "Baru saja"
	case diff < time.Hour:
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 menit yang lalu"
		}
		return fmt.Sprintf("%d menit yang lalu", minutes)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 jam yang lalu"
		}
		return fmt.Sprintf("%d jam yang lalu", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 hari yang lalu"
		}
		return fmt.Sprintf("%d hari yang lalu", days)
	default:
		return t.Format("2 Jan 2006")
	}
}
