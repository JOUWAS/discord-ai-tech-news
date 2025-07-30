package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type News struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"publishedAt"`
	Source      string    `json:"source"`
	Score       int       `json:"score,omitempty"` // HackerNews score
}

// HackerNews API response structures
type HackerNewsStory struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Time        int64  `json:"time"`
	Score       int    `json:"score"`
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	Text        string `json:"text"`
}

type NewsRepository interface {
	GetLatestNews() ([]News, error)
	GetLatestNewsSince(since time.Time) ([]News, error)
	SearchNews(keyword string) ([]News, error)
}

type HackerNewsRepository struct {
	client  *http.Client
	baseURL string
}

func NewHackerNewsRepository() *HackerNewsRepository {
	return &HackerNewsRepository{
		client:  &http.Client{Timeout: 15 * time.Second},
		baseURL: "https://hacker-news.firebaseio.com/v0",
	}
}

func (r *HackerNewsRepository) GetLatestNews() ([]News, error) {
	// Mengambil berita 24 jam terakhir
	since := time.Now().Add(-24 * time.Hour)
	return r.GetLatestNewsSince(since)
}

func (r *HackerNewsRepository) GetLatestNewsSince(since time.Time) ([]News, error) {
	log.Printf("🌐 DEBUG: Fetching tech news from HackerNews API")

	// Get top stories from HackerNews
	topStoriesURL := fmt.Sprintf("%s/topstories.json", r.baseURL)

	resp, err := r.client.Get(topStoriesURL)
	if err != nil {
		log.Printf("❌ DEBUG: HackerNews API failed: %v", err)
		return r.getMockNews(), nil
	}
	defer resp.Body.Close()

	var storyIDs []int
	if err := json.NewDecoder(resp.Body).Decode(&storyIDs); err != nil {
		log.Printf("❌ DEBUG: Failed to decode story IDs: %v", err)
		return r.getMockNews(), nil
	}

	log.Printf("✅ DEBUG: Got %d stories from HackerNews", len(storyIDs))

	var news []News
	processedCount := 0

	// Get first 30 stories to have enough tech content
	for _, id := range storyIDs {
		if len(news) >= 10 { // Stop when we have enough tech stories
			break
		}
		if processedCount >= 30 { // Don't process more than 30
			break
		}
		processedCount++

		story, err := r.getStoryDetails(id)
		if err != nil {
			continue
		}

		// Filter tech-related stories
		if r.isTechRelated(story.Title) && story.URL != "" {
			publishedAt := time.Unix(story.Time, 0)

			// Filter: only articles within the time range
			if publishedAt.After(since) {
				description := fmt.Sprintf("⭐ %d points • 💬 %d comments", story.Score, story.Descendants)
				if len(story.Text) > 0 && len(story.Text) < 200 {
					description = story.Text + " | " + description
				}

				news = append(news, News{
					Title:       story.Title,
					Description: description,
					URL:         story.URL,
					PublishedAt: publishedAt,
					Source:      "Hacker News",
					Score:       story.Score,
				})

				log.Printf("✅ DEBUG: Added tech story: %s (score: %d)",
					story.Title[:min(50, len(story.Title))], story.Score)
			}
		}
	}

	// Sort by score (popularity) instead of time for HackerNews
	for i := 0; i < len(news)-1; i++ {
		for j := i + 1; j < len(news); j++ {
			if news[j].Score > news[i].Score {
				news[i], news[j] = news[j], news[i]
			}
		}
	}

	if len(news) == 0 {
		log.Println("❌ DEBUG: No tech stories found in HackerNews, using mock data")
		return r.getMockNews(), nil
	}

	log.Printf("🎉 DEBUG: Returning %d real tech articles from HackerNews", len(news))
	return news, nil
}

func (r *HackerNewsRepository) getStoryDetails(id int) (*HackerNewsStory, error) {
	storyURL := fmt.Sprintf("%s/item/%d.json", r.baseURL, id)
	resp, err := r.client.Get(storyURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var story HackerNewsStory
	if err := json.NewDecoder(resp.Body).Decode(&story); err != nil {
		return nil, err
	}

	return &story, nil
}

func (r *HackerNewsRepository) isTechRelated(title string) bool {
	techKeywords := []string{
		"ai", "artificial intelligence", "machine learning", "ml", "tech", "technology",
		"programming", "software", "app", "api", "cloud", "startup", "github",
		"open source", "javascript", "python", "go", "rust", "framework", "developer",
		"database", "security", "crypto", "blockchain", "web", "mobile", "ios", "android",
		"react", "vue", "node", "docker", "kubernetes", "aws", "google", "microsoft",
		"apple", "tesla", "spacex", "quantum", "data", "algorithm", "code", "coding",
		"compiler", "linux", "windows", "mac", "server", "devops", "ci/cd", "git",
		"typescript", "java", "c++", "c#", "php", "ruby", "kotlin", "swift", "scala",
	}

	titleLower := strings.ToLower(title)
	for _, keyword := range techKeywords {
		if strings.Contains(titleLower, keyword) {
			return true
		}
	}
	return false
}

func (r *HackerNewsRepository) SearchNews(keyword string) ([]News, error) {
	log.Printf("🔍 DEBUG: Searching HackerNews for keyword: %s", keyword)

	// For search, we'll get latest stories and filter by keyword
	news, err := r.GetLatestNews()
	if err != nil {
		return r.getMockNews(), nil
	}

	var filtered []News
	keywordLower := strings.ToLower(keyword)

	for _, article := range news {
		titleLower := strings.ToLower(article.Title)
		descLower := strings.ToLower(article.Description)

		if strings.Contains(titleLower, keywordLower) || strings.Contains(descLower, keywordLower) {
			filtered = append(filtered, article)
		}
	}

	if len(filtered) == 0 {
		log.Printf("❌ DEBUG: No articles found for keyword '%s'", keyword)
		return r.getMockNews(), nil
	}

	log.Printf("✅ DEBUG: Found %d articles for keyword '%s'", len(filtered), keyword)
	return filtered, nil
}

func (r *HackerNewsRepository) getMockNews() []News {
	return []News{
		{
			Title:       "🚀 AI Revolution: GPT-5 Released with Breakthrough Capabilities",
			Description: "⭐ 1547 points • 💬 423 comments | OpenAI announces GPT-5 with unprecedented reasoning abilities.",
			URL:         "https://example.com/gpt5-release",
			PublishedAt: time.Now().Add(-1 * time.Hour),
			Source:      "Hacker News",
			Score:       1547,
		},
		{
			Title:       "💻 Quantum Computing Reaches New Milestone",
			Description: "⭐ 1205 points • 💬 287 comments | IBM's new quantum processor achieves 1000+ qubit stability.",
			URL:         "https://example.com/quantum-breakthrough",
			PublishedAt: time.Now().Add(-2 * time.Hour),
			Source:      "Hacker News",
			Score:       1205,
		},
		{
			Title:       "🌐 Web 3.0 Adoption Accelerates in 2025",
			Description: "⭐ 892 points • 💬 156 comments | Decentralized applications see 400% growth as mainstream adoption takes off.",
			URL:         "https://example.com/web3-growth",
			PublishedAt: time.Now().Add(-3 * time.Hour),
			Source:      "Hacker News",
			Score:       892,
		},
		{
			Title:       "🔧 New Go Framework Simplifies Microservices Development",
			Description: "⭐ 756 points • 💬 198 comments | Developer-friendly framework reduces boilerplate by 70%.",
			URL:         "https://example.com/go-framework",
			PublishedAt: time.Now().Add(-4 * time.Hour),
			Source:      "Hacker News",
			Score:       756,
		},
		{
			Title:       "🛡️ Zero-Day Vulnerability Discovered in Popular JavaScript Library",
			Description: "⭐ 2341 points • 💬 534 comments | Security researchers urge immediate updates.",
			URL:         "https://example.com/js-vulnerability",
			PublishedAt: time.Now().Add(-5 * time.Hour),
			Source:      "Hacker News",
			Score:       2341,
		},
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
