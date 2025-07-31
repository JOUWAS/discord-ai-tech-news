package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type News struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"publishedAt"`
	Source      string    `json:"source"`
}

type NewsAPIResponse struct {
	Status       string        `json:"status"`
	TotalResults int           `json:"totalResults"`
	Articles     []NewsArticle `json:"articles"`
}

type NewsArticle struct {
	Source      NewsSource `json:"source"`
	Author      string     `json:"author"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	URL         string     `json:"url"`
	URLToImage  string     `json:"urlToImage"`
	PublishedAt string     `json:"publishedAt"`
	Content     string     `json:"content"`
}

type NewsSource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type NewsRepository interface {
	GetLatestNews() ([]News, error)
	GetLatestNewsSince(since time.Time) ([]News, error)
	SearchNews(keyword string) ([]News, error)
}

type NewsApiRepository struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

func NewNewsApiRepository() *NewsApiRepository {
	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Fatal("⚠️ WARNING: NEWS_API_KEY is not set in environment variables")
	}

	return &NewsApiRepository{
		client:  &http.Client{Timeout: 15 * time.Second},
		baseURL: "https://newsapi.org/v2",
		apiKey:  apiKey,
	}
}

func (r *NewsApiRepository) GetLatestNews() ([]News, error) {
	since := time.Now().Add(-24 * time.Hour)
	return r.GetLatestNewsSince(since)
}

func (r *NewsApiRepository) GetLatestNewsSince(since time.Time) ([]News, error) {
	log.Printf("🌐 DEBUG: Fetching tech news from News API")

	fromDate := since.Format("2006-01-02")
	url := fmt.Sprintf("%s/everything?q=technology&from=%s&sortBy=popularity&pageSize=20&apiKey=%s",
		r.baseURL, fromDate, r.apiKey)

	log.Printf("🔗 DEBUG: NewsAPI URL (without API key): %s/everything?q=technology&from=%s&sortBy=popularity&pageSize=20&apiKey=***", r.baseURL, fromDate)

	resp, err := r.client.Get(url)
	if err != nil {
		log.Printf("❌ ERROR: NewsAPI failed: %v", err)
		return r.getMockNews(), nil
	}
	defer resp.Body.Close()

	log.Printf("📊 DEBUG: Response Status: %s", resp.Status)
	log.Printf("📊 DEBUG: Response Status Code: %d", resp.StatusCode)

	// Read response body as bytes first for debugging
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ ERROR: Failed to read response body: %v", err)
		return r.getMockNews(), nil
	}

	log.Printf("📋 DEBUG: Raw response body length: %d bytes", len(bodyBytes))
	if len(bodyBytes) == 0 {
		log.Printf("❌ ERROR: Empty response body from NewsAPI")
		return r.getMockNews(), nil
	}

	// Log first 500 characters of response for debugging
	if len(bodyBytes) > 500 {
		log.Printf("📋 DEBUG: Response preview: %s...", string(bodyBytes[:500]))
	} else {
		log.Printf("📋 DEBUG: Full response: %s", string(bodyBytes))
	}

	var apiResponse NewsAPIResponse
	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		log.Printf("❌ ERROR: Failed to decode NewsAPI response: %v", err)
		log.Printf("📋 DEBUG: Response body was: %s", string(bodyBytes))
		return r.getMockNews(), nil
	}

	if apiResponse.Status != "ok" {
		log.Printf("❌ ERROR: NewsAPI returned status %s", apiResponse.Status)
		return r.getMockNews(), nil
	}

	log.Printf("✅ SUCCESS: Got %d articles from NewsAPI", len(apiResponse.Articles))

	var news []News
	for _, article := range apiResponse.Articles {
		publishedAt, err := time.Parse(time.RFC3339, article.PublishedAt)
		if err != nil {
			publishedAt = time.Now()
		}

		if article.Title != "" && article.URL != "" {
			news = append(news, News{
				Title:       article.Title,
				Description: article.Description,
				URL:         article.URL,
				PublishedAt: publishedAt,
				Source:      article.Source.Name,
			})

			log.Printf("📰 DEBUG: Added article: %s from %s",
				article.Title, article.Source.Name)
		}
	}

	if len(news) == 0 {
		log.Println("⚠️ WARNING: No valid articles found in NewsAPI response, using mock data")
		return r.getMockNews(), nil
	}

	log.Printf("✅ DEBUG: Returning %d tech news articles from NewsAPI", len(news))
	return news, nil
}

func (r *NewsApiRepository) SearchNews(keyword string) ([]News, error) {
	log.Printf("🔍 DEBUG: Searching NewsAPI for keyword: %s", keyword)

	url := fmt.Sprintf("%s/everything?q=%s&sortBy=relevancy&pageSize=10&apiKey=%s",
		r.baseURL, keyword, r.apiKey)

	resp, err := r.client.Get(url)
	if err != nil {
		log.Printf("❌ ERROR: NewsAPI search failed: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	var apiResponse NewsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		log.Printf("❌ ERROR: Failed to decode NewsAPI search response: %v", err)
		return r.getMockNews(), nil
	}

	if apiResponse.Status != "ok" {
		log.Printf("❌ ERROR: NewsAPI search returned status %s", apiResponse.Status)
		return r.getMockNews(), nil
	}

	var news []News
	for _, article := range apiResponse.Articles {
		publishedAt, err := time.Parse(time.RFC3339, article.PublishedAt)
		if err != nil {
			publishedAt = time.Now()
		}

		if article.Title != "" && article.URL != "" {
			news = append(news, News{
				Title:       article.Title,
				Description: article.Description,
				URL:         article.URL,
				PublishedAt: publishedAt,
				Source:      article.Source.Name,
			})

			log.Printf("📰 DEBUG: Added search result: %s from %s",
				article.Title, article.Source.Name)
		}
	}

	if len(news) == 0 {
		log.Println("⚠️ WARNING: No valid articles found in NewsAPI response, using mock data")
		return r.getMockNews(), nil
	}

	log.Printf("✅ DEBUG: Returning %d tech news articles from NewsAPI", len(news))
	return news, nil
}

func (r *NewsApiRepository) getMockNews() []News {
	return []News{
		{
			Title:       "🚀 AI Revolution: GPT-5 Released with Breakthrough Capabilities",
			Description: "⭐ 1547 points • 💬 423 comments | OpenAI announces GPT-5 with unprecedented reasoning abilities.",
			URL:         "https://example.com/gpt5-release",
			PublishedAt: time.Now().Add(-1 * time.Hour),
			Source:      "Mock News",
		},
		{
			Title:       "💻 Quantum Computing Reaches New Milestone",
			Description: "⭐ 1205 points • 💬 287 comments | IBM's new quantum processor achieves 1000+ qubit stability.",
			URL:         "https://example.com/quantum-breakthrough",
			PublishedAt: time.Now().Add(-2 * time.Hour),
			Source:      "Mock News",
		},
		{
			Title:       "🌐 Web 3.0 Adoption Accelerates in 2025",
			Description: "⭐ 892 points • 💬 156 comments | Decentralized applications see 400% growth as mainstream adoption takes off.",
			URL:         "https://example.com/web3-growth",
			PublishedAt: time.Now().Add(-3 * time.Hour),
			Source:      "Mock News",
		},
		{
			Title:       "🔧 New Go Framework Simplifies Microservices Development",
			Description: "⭐ 756 points • 💬 198 comments | Developer-friendly framework reduces boilerplate by 70%.",
			URL:         "https://example.com/go-framework",
			PublishedAt: time.Now().Add(-4 * time.Hour),
			Source:      "Mock News",
		},
		{
			Title:       "🛡️ Zero-Day Vulnerability Discovered in Popular JavaScript Library",
			Description: "⭐ 2341 points • 💬 534 comments | Security researchers urge immediate updates.",
			URL:         "https://example.com/js-vulnerability",
			PublishedAt: time.Now().Add(-5 * time.Hour),
			Source:      "Mock News",
		},
	}
}
