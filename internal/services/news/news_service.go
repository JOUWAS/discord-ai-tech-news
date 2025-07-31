package news

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	newsRepo "discord-ai-tech-news/internal/repositories/news"
)

type News = newsRepo.News

type NewsResponse struct {
	News []News `json:"news"`
}

type NewsService interface {
	FetchTechNews(ctx context.Context) (*NewsResponse, error)
	SearchNews(ctx context.Context, keyword string) ([]News, error)
	ValidateNewsSource(source string) bool
	FormatNewsForDiscord(news []News) string
	TimeAgo(t time.Time) string
}

type ExternalNewsService struct {
	repository newsRepo.NewsRepository
}

func NewExternalNewsService(repo newsRepo.NewsRepository) *ExternalNewsService {
	return &ExternalNewsService{
		repository: repo,
	}
}

func (s *ExternalNewsService) FetchTechNews(ctx context.Context) (*NewsResponse, error) {
	// Ambil berita teknologi dari 24 jam terakhir
	since := time.Now().Add(-24 * time.Hour)
	news, err := s.repository.GetLatestNewsSince(since)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch news: %w", err)
	}

	// Filter tech-related news
	techNews := s.filterTechNews(news)

	// Limit ke 5 berita terbaru untuk performa
	if len(techNews) > 5 {
		techNews = techNews[:5]
	}

	return &NewsResponse{News: techNews}, nil
}

func (s *ExternalNewsService) SearchNews(ctx context.Context, keyword string) ([]News, error) {
	log.Printf("🔍 DEBUG: Service searching for: %s", keyword)

	// Call repository search
	results, err := s.repository.SearchNews(keyword)
	if err != nil {
		return nil, fmt.Errorf("failed to search news: %w", err)
	}

	// Filter and validate results
	validResults := s.filterSearchResults(results, keyword)

	log.Printf("✅ DEBUG: Search returned %d valid results", len(validResults))
	return validResults, nil
}

func (s *ExternalNewsService) ValidateNewsSource(source string) bool {
	validSources := []string{
		"techcrunch", "wired", "ars technica", "ieee spectrum",
		"the verge", "engadget", "gizmodo", "cnet",
	}

	sourceLower := strings.ToLower(source)
	for _, valid := range validSources {
		if strings.Contains(sourceLower, valid) {
			return true
		}
	}
	return false
}

func (s *ExternalNewsService) FormatNewsForDiscord(news []News) string {
	if len(news) == 0 {
		return "📰 **Tech News Update**\n\nTidak ada berita teknologi terbaru saat ini."
	}

	var result strings.Builder
	result.WriteString("📰 **Tech News Update - Berita Teknologi Terbaru**\n\n")

	for i, article := range news {
		if i >= 3 { // Limit ke 3 berita untuk Discord
			break
		}

		// Format waktu yang user-friendly
		timeAgo := s.TimeAgo(article.PublishedAt)

		result.WriteString(fmt.Sprintf("**%d. %s**\n", i+1, article.Title))
		if article.Description != "" {
			description := article.Description
			if len(description) > 150 {
				description = description[:150] + "..."
			}
			result.WriteString(fmt.Sprintf("📝 %s\n", description))
		}
		result.WriteString(fmt.Sprintf("🔗 [Baca Selengkapnya](%s)\n", article.URL))
		result.WriteString(fmt.Sprintf("📅 %s • 📰 %s\n\n", timeAgo, article.Source))
	}

	result.WriteString("---\n💡 *Ketik `help` untuk melihat command lainnya*")
	return result.String()
}

func (s *ExternalNewsService) filterTechNews(news []News) []News {
	techKeywords := []string{
		"ai", "artificial intelligence", "technology", "tech", "software",
		"programming", "computer", "mobile", "app", "startup", "innovation",
		"quantum", "blockchain", "crypto", "web3", "machine learning", "cloud",
	}

	var filtered []News
	for _, article := range news {
		if s.containsTechKeywords(article, techKeywords) {
			filtered = append(filtered, article)
		}
	}

	// Jika tidak ada yang match dengan keywords, return semua (assume semuanya tech news)
	if len(filtered) == 0 {
		return news
	}

	return filtered
}

func (s *ExternalNewsService) containsTechKeywords(article News, keywords []string) bool {
	content := strings.ToLower(article.Title + " " + article.Description)

	for _, keyword := range keywords {
		if strings.Contains(content, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func (s *ExternalNewsService) TimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "Baru saja"
	case diff < time.Hour:
		minutes := int(diff.Minutes())
		return fmt.Sprintf("%d menit yang lalu", minutes)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		return fmt.Sprintf("%d jam yang lalu", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d hari yang lalu", days)
	default:
		return t.Format("2 Jan 2006")
	}
}

func (s *ExternalNewsService) filterSearchResults(results []News, keyword string) []News {
	var filtered []News

	keywordLower := strings.ToLower(keyword)

	for _, article := range results {
		// Check if article is relevant
		titleLower := strings.ToLower(article.Title)
		descLower := strings.ToLower(article.Description)

		// Must contain the keyword and have valid content
		if (strings.Contains(titleLower, keywordLower) || strings.Contains(descLower, keywordLower)) &&
			article.Title != "" && article.URL != "" && article.Title != "[Removed]" {
			filtered = append(filtered, article)
		}
	}

	// Sort by published date (newest first)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].PublishedAt.After(filtered[j].PublishedAt)
	})

	return filtered
}
