package actions

import (
	"fmt"
	"sync"

	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Topic struct {
	gorm.Model
	Name string
}

type TopicData struct {
	Name string `form:"name"`
}

type Site struct {
	gorm.Model
	Url string
}

type SiteData struct {
	Url string `form:"url"`
}

type User struct {
	gorm.Model
	Email string
}

type UserData struct {
	Email string `form:"email"`
}

type Subscription struct {
	gorm.Model
	UserID                 uint
	User                   User
	Topics                 pq.StringArray `gorm:"type:text[]"`
	Sites                  pq.StringArray `gorm:"type:text[]"`
	SubscriptionScheduleID uint
	SubscriptionSchedule   SubscriptionSchedule
}

type SubscriptionData struct {
	Email                    string                   `form:"email" json:"email"`
	Topics                   pq.StringArray           `form:"topics" json:"topics"`
	Sites                    pq.StringArray           `form:"sites" json:"sites"`
	SubscriptionScheduleData SubscriptionScheduleData `form:"subscriptionSchedule" json:"subscriptionSchedule"`
}

type DailyFrequency struct {
	Monday    bool `form:"monday" json:"monday"`
	Tuesday   bool `form:"tuesday" json:"tuesday"`
	Wednesday bool `form:"wednesday" json:"wednesday"`
	Thursday  bool `form:"thursday" json:"thursday"`
	Friday    bool `form:"friday" json:"friday"`
	Saturday  bool `form:"staturday" json:"staturday"`
	Sunday    bool `form:"sunday" json:"sunday"`
}

type TimeSlot string

const (
	Morning   TimeSlot = "Morning"
	Afternoon TimeSlot = "Afternoon"
	Evening   TimeSlot = "Evening"
	Night     TimeSlot = "Night"
)

type SubscriptionScheduleData struct {
	DailyFrequency DailyFrequency `form:"dailyFrequency" json:"dailyFrequency"`
	TimeSlot       TimeSlot       `form:"timeSlot" json:"timeSlot"`
	TimeZone       string         `form:"timezone" json:"timezone"`
}

type SubscriptionSchedule struct {
	gorm.Model
	Sunday       bool
	Monday       bool
	Tuesday      bool
	Wednesday    bool
	Thursday     bool
	Friday       bool
	Saturday     bool
	TimeSlotEnum TimeSlot
	TimeZone     string
}

func CreateTopic(topicData TopicData) {
	var topicDB Topic
	db().Where(Topic{Name: topicData.Name}).FirstOrCreate(&topicDB)
}

func CreateSite(siteData SiteData) {
	var siteDb Site
	db().Where(Site{Url: siteData.Url}).FirstOrCreate(&siteDb)
}

func CreateUser(userData UserData) User {
	userDb := User{Email: userData.Email}
	var user User
	db().Where(userDb).FirstOrCreate(&user, userDb)
	return user
}

func GetSubscriptionByEmail(email string) Subscription {
	var user User
	db().Find(&user, User{Email: email})

	fmt.Printf("user id is %v \n", user.Email)

	var subscription Subscription

	r :=
		db().Joins("SubscriptionSchedule").Joins("User").Joins("SubscriptionSchedule").Find(&subscription, Subscription{UserID: user.ID})
	fmt.Printf("Query is %v \n", r.Statement.SQL.String())

	fmt.Printf("sub id is %v \n", subscription.ID)

	return subscription
}

func GetSubscriptionByID(id int) Subscription {

	var subscription Subscription

	r := db().Joins("SubscriptionSchedule").Joins("User").Find(&subscription, id)
	fmt.Printf("Query is %v \n", r.Statement.SQL.String())

	fmt.Printf("sub id is %v \n", subscription.ID)

	return subscription
}

func GetAllSubscriptions() []Subscription {

	var subscriptions []Subscription

	r := db().Joins("SubscriptionSchedule").Joins("User").Find(&subscriptions)
	fmt.Printf("Query is %v \n", r.Statement.SQL.String())

	return subscriptions
}
func CreateSubscriptionSchedule(subscriptionScheduleData SubscriptionScheduleData) SubscriptionSchedule {
	subscriptionScheduleDb := SubscriptionSchedule{
		Monday:       subscriptionScheduleData.DailyFrequency.Monday,
		Tuesday:      subscriptionScheduleData.DailyFrequency.Tuesday,
		Wednesday:    subscriptionScheduleData.DailyFrequency.Wednesday,
		Thursday:     subscriptionScheduleData.DailyFrequency.Thursday,
		Friday:       subscriptionScheduleData.DailyFrequency.Friday,
		Saturday:     subscriptionScheduleData.DailyFrequency.Saturday,
		Sunday:       subscriptionScheduleData.DailyFrequency.Sunday,
		TimeSlotEnum: subscriptionScheduleData.TimeSlot,
		TimeZone:     subscriptionScheduleData.TimeZone,
	}

	fmt.Printf("Schedule %v\n", subscriptionScheduleDb)
	var subscriptionSchedule SubscriptionSchedule
	db().Where(subscriptionScheduleDb).FirstOrCreate(&subscriptionSchedule, subscriptionScheduleDb)
	return subscriptionSchedule
}

func CreateSubscription(subscriptionData SubscriptionData, user User, subscriptionSchedule SubscriptionSchedule) Subscription {
	attrs := Subscription{
		UserID: user.ID,
	}
	values := Subscription{
		Topics:                 pq.StringArray(subscriptionData.Topics),
		Sites:                  pq.StringArray(subscriptionData.Sites),
		SubscriptionScheduleID: subscriptionSchedule.ID,
	}
	var subscription Subscription
	db().Where(attrs).Assign(values).FirstOrCreate(&subscription)
	return subscription
}

func Migrate() {
	db().AutoMigrate(&Topic{}, &Subscription{}, &Site{}, &User{}, &SubscriptionSchedule{})
}

var (
	dataBase *gorm.DB
	once     sync.Once
	err      error
)

func db() *gorm.DB {
	once.Do(func() {
		dsn := "host=localhost user=postgres password=password dbname=news-master port=5432 sslmode=disable TimeZone=Asia/Shanghai"
		dataBase, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
		if err != nil {
			panic("Unable to connect to db")
		}
		sqlDB, _ := dataBase.DB()
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetMaxIdleConns(5)
	})

	return dataBase
}
