package entity

import (
	"news-master/datamodels/common"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Site struct {
	gorm.Model
	Url      string `gorm:"index:idx_url"`
	Name     string
	Language string
	Active   bool `gorm:"index:"`
}
type User struct {
	gorm.Model
	Email             string `gorm:"uniqueIndex"`
	LoginAttemptCount int
}

type Feedback struct {
	gorm.Model
	Content string
}

type Subscription struct {
	gorm.Model
	UserID                 uint `gorm:"uniqueIndex"`
	User                   User
	SubscriptionScheduleID uint
	SubscriptionSchedule   SubscriptionSchedule
	Confirmed              bool
	LastProcessedAt        time.Time
	Sites                  []Site `gorm:"many2many:subscription_sites;"`
}
type SubscriptionSchedule struct {
	gorm.Model
	Sunday    bool
	Monday    bool
	Tuesday   bool
	Wednesday bool
	Thursday  bool
	Friday    bool
	Saturday  bool
	TimeSlot  common.TimeSlot
	TimeZone  string
}

type Article struct {
	gorm.Model
	ArticleExternalId string `gorm:"index:idx_ext_id"`
	Title             string
	Link              string
	Description       string
	Content           string
	ImageURL          string
	Language          string
	Country           pq.StringArray `gorm:"type:text[]"`
	Category          pq.StringArray `gorm:"type:text[]"`
	Site              string         `gorm:"index:site"`
}
