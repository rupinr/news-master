package dto

import (
	"news-master/datamodels/common"
	"news-master/datamodels/entity"
)

func MapToUrls(sites []entity.Site) []string {
	urls := make([]string, len(sites))
	for i, site := range sites {
		urls[i] = site.Url
	}
	return urls
}

type Topic struct {
	Name    string `json:"name" binding:"required"`
	Visible bool   `json:"visible"`
}

type TopicUpdate struct {
	Visibility *bool `json:"visibility" binding:"required"`
}

type Site struct {
	Url      string `json:"url" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Active   *bool  `json:"active"`
	Language string `json:"language"`
}
type Sites struct {
	Sites []Site `json:"sites,omitempty"`
}

type NewsletterData struct {
	AboutLink              string
	PrivacyLink            string
	ManageSubscriptionLink string
	Articles               []entity.Article
}

type User struct {
	Email string `json:"email" validate:"required,emailValidator"`
}

type Subscription struct {
	Confirmed                bool                 `json:"confirmed"`
	Sites                    []string             `json:"sites" binding:"required"`
	SubscriptionScheduleData SubscriptionSchedule `json:"subscriptionSchedule" binding:"required"`
}

type SubscriptionConfirmation struct {
	Confirmed *bool `json:"confirmed" binding:"required"`
}

type DailyFrequency struct {
	Monday    *bool `json:"monday" binding:"required"`
	Tuesday   *bool `json:"tuesday" binding:"required"`
	Wednesday *bool `json:"wednesday" binding:"required"`
	Thursday  *bool `json:"thursday" binding:"required"`
	Friday    *bool `json:"friday" binding:"required"`
	Saturday  *bool `json:"saturday" binding:"required"`
	Sunday    *bool `json:"sunday" binding:"required"`
}

type SubscriptionSchedule struct {
	DailyFrequency DailyFrequency  `json:"dailyFrequency" binding:"required"`
	TimeSlot       common.TimeSlot `json:"timeSlot" binding:"required"`
	TimeZone       string          `json:"timezone" binding:"required"`
}

type Feedback struct {
	Content string `json:"content" binding:"required"`
}

type Article struct {
	ArticleId   string   `json:"article_id"`
	Title       string   `json:"title"`
	Link        string   `json:"link"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	ImageURL    string   `json:"image_url"`
	Language    string   `json:"language"`
	Country     []string `json:"country"`
	Category    []string `json:"category"`
	SourceUrl   string   `json:"source_url"`
}

type NewsdataApiResponse struct {
	Results []Article `json:"results"`
}
