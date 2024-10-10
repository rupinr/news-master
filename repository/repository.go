package repository

import (
	"fmt"
	"news-master/datamodels/dto"
	"news-master/datamodels/entity"
	"sync"

	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CreateTopic(topicData dto.Topic) {
	var topicDB entity.Topic
	db().Where(entity.Topic{Name: topicData.Name}).FirstOrCreate(&topicDB)
}

func CreateSite(siteData dto.Site) {
	var siteDb entity.Site
	db().Where(entity.Site{Url: siteData.Url}).FirstOrCreate(&siteDb)
}

func CreateUser(userData dto.User) entity.User {
	userDb := entity.User{Email: userData.Email}
	var user entity.User
	db().Where(userDb).FirstOrCreate(&user, userDb)
	return user
}

func GetSubscriptionByEmail(email string) entity.Subscription {
	var user entity.User
	db().Find(&user, entity.User{Email: email})

	fmt.Printf("user id is %v \n", user.Email)

	var subscription entity.Subscription

	r :=
		db().Joins("SubscriptionSchedule").Joins("User").Joins("SubscriptionSchedule").Find(&subscription, entity.Subscription{UserID: user.ID})
	fmt.Printf("Query is %v \n", r.Statement.SQL.String())

	fmt.Printf("sub id is %v \n", subscription.ID)

	return subscription
}

func GetSubscriptionByID(id int) entity.Subscription {

	var subscription entity.Subscription

	r := db().Joins("SubscriptionSchedule").Joins("User").Find(&subscription, id)
	fmt.Printf("Query is %v \n", r.Statement.SQL.String())

	fmt.Printf("sub id is %v \n", subscription.ID)

	return subscription
}

func GetAllSubscriptions() []entity.Subscription {

	var subscriptions []entity.Subscription

	r := db().Joins("SubscriptionSchedule").Joins("User").Find(&subscriptions)
	fmt.Printf("Query is %v \n", r.Statement.SQL.String())

	return subscriptions
}
func CreateSubscriptionSchedule(subscriptionScheduleData dto.SubscriptionSchedule) entity.SubscriptionSchedule {
	subscriptionScheduleDb := entity.SubscriptionSchedule{
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
	var subscriptionSchedule entity.SubscriptionSchedule
	db().Where(subscriptionScheduleDb).FirstOrCreate(&subscriptionSchedule, subscriptionScheduleDb)
	return subscriptionSchedule
}

func CreateSubscription(subscriptionData dto.Subscription, user entity.User, subscriptionSchedule entity.SubscriptionSchedule) entity.Subscription {
	attrs := entity.Subscription{
		UserID: user.ID,
	}
	values := entity.Subscription{
		Topics:                 pq.StringArray(subscriptionData.Topics),
		Sites:                  pq.StringArray(subscriptionData.Sites),
		SubscriptionScheduleID: subscriptionSchedule.ID,
	}
	var subscription entity.Subscription
	db().Where(attrs).Assign(values).FirstOrCreate(&subscription)
	return subscription
}

func Migrate() {
	db().AutoMigrate(&entity.Topic{}, &entity.Subscription{}, &entity.Site{}, &entity.User{}, &entity.SubscriptionSchedule{})
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
