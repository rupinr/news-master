package repository

import (
	"fmt"
	"news-master/app"
	"news-master/datamodels/dto"
	"news-master/datamodels/entity"
	"sync"
	"time"

	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbError struct {
	Message string
	Code    int
}

func (e *DbError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func CreateTopic(topicData dto.Topic) (entity.Topic, error) {
	topic := entity.Topic{Visible: false}
	err := db().Where(entity.Topic{Name: topicData.Name}).FirstOrCreate(&topic).Error
	return topic, err
}

func UpdateTopic(name string, visibility bool) error {
	err = db().Transaction(func(tx *gorm.DB) error {
		var topic entity.Topic
		if err := tx.Where(&entity.Topic{Name: name}).First(&topic).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("topic not found")
			}
			return err
		}
		topic.Visible = visibility
		if err := tx.Save(&topic).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func CreateSite(siteData dto.Site) {
	var siteDb entity.Site
	db().Where(entity.Site{Url: siteData.Url}).FirstOrCreate(&siteDb)
}

func CreateResult(result dto.Article) {
	article := entity.Article{
		Title:       result.Title,
		Link:        result.Link,
		Description: result.Description,
		Content:     result.Content,
		ImageURL:    result.ImageURL,
		Language:    result.Language,
		Country:     pq.StringArray(result.Country),
		Category:    pq.StringArray(result.Category),
		Site:        result.SourceUrl,
	}
	db().Create(&article)
}

func GetActiveSites() []entity.Site {
	var sites []entity.Site
	db().Where(entity.Site{Active: true}).Find(&sites)
	return sites
}

func CreateUser(userData dto.User) entity.User {
	userDb := entity.User{Email: userData.Email}
	var user entity.User
	db().Where(userDb).FirstOrCreate(&user, userDb)
	return user
}

func GetUser(userData dto.User) entity.User {
	userDb := entity.User{Email: userData.Email}
	var user entity.User
	db().Where(userDb).First(&user, userDb)
	return user
}

func IncrementAndGetLoginAttempt(userData dto.User) entity.User {
	user := entity.User{}
	db().Where(entity.User{Email: userData.Email}).Find(&user)
	user.LoginAttemptCount++
	db().Save(&user)
	return user
}

func GetSubscriptionByEmail(email string) (entity.Subscription, error) {
	var user entity.User
	result := db().First(&user, entity.User{Email: email})
	fmt.Println(result.Error)
	var subscription entity.Subscription

	db().Joins("SubscriptionSchedule").Joins("User").Joins("SubscriptionSchedule").Find(&subscription, entity.Subscription{UserID: user.ID})

	return subscription, result.Error
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
	conditions := entity.Subscription{Confirmed: true}
	r := db().Joins("SubscriptionSchedule").Joins("User").Find(&subscriptions, conditions)
	fmt.Printf("Query is %v \n", r.Statement.SQL.String())

	return subscriptions
}

func SetLastProcessedAt(subscriptionId uint) {
	time := time.Now()
	fmt.Printf("setting last processed at as %v for sub %v\n", time, subscriptionId)
	var sub entity.Subscription
	db().Find(&sub, subscriptionId)
	sub.LastProcessedAt = time
	db().Save(&sub)
}

func GetArticlesFrom(fromDate time.Time, sites []string) []entity.Article {
	var articles []entity.Article
	db().Where("created_at > ?", fromDate).Where("site IN ?", sites).Find(&articles)
	return articles

}

func CreateSubscriptionSchedule(subscriptionScheduleData dto.SubscriptionSchedule) entity.SubscriptionSchedule {

	subscriptionScheduleDb := entity.SubscriptionSchedule{
		Monday:    *subscriptionScheduleData.DailyFrequency.Monday,
		Tuesday:   *subscriptionScheduleData.DailyFrequency.Tuesday,
		Wednesday: *subscriptionScheduleData.DailyFrequency.Wednesday,
		Thursday:  *subscriptionScheduleData.DailyFrequency.Thursday,
		Friday:    *subscriptionScheduleData.DailyFrequency.Friday,
		Saturday:  *subscriptionScheduleData.DailyFrequency.Saturday,
		Sunday:    *subscriptionScheduleData.DailyFrequency.Sunday,
		TimeSlot:  subscriptionScheduleData.TimeSlot,
		TimeZone:  subscriptionScheduleData.TimeZone,
	}

	// Map to hold conditions
	conditions := map[string]interface{}{
		"monday":    subscriptionScheduleDb.Monday,
		"tuesday":   subscriptionScheduleDb.Tuesday,
		"wednesday": subscriptionScheduleDb.Wednesday,
		"thursday":  subscriptionScheduleDb.Thursday,
		"friday":    subscriptionScheduleDb.Friday,
		"saturday":  subscriptionScheduleDb.Saturday,
		"sunday":    subscriptionScheduleDb.Sunday,
		"time_slot": subscriptionScheduleDb.TimeSlot,
		"time_zone": subscriptionScheduleDb.TimeZone,
	}

	var subscriptionSchedule entity.SubscriptionSchedule

	if err := db().Where(conditions).FirstOrCreate(&subscriptionSchedule, subscriptionScheduleDb).Error; err != nil {
		fmt.Println("Error creating or finding record:", err)
	} else {
		fmt.Printf("Record found or created: %+v\n", subscriptionSchedule)
	}

	return subscriptionSchedule
}

func CreateSubscription(user entity.User, sites []string, subscriptionScheduleID uint) entity.Subscription {
	attrs := entity.Subscription{
		UserID: user.ID,
	}
	values := entity.Subscription{
		Sites:                  pq.StringArray(sites),
		SubscriptionScheduleID: subscriptionScheduleID,
		Confirmed:              false,
	}
	var subscription entity.Subscription
	db().Where(attrs).Assign(values).FirstOrCreate(&subscription)
	return subscription
}

func UpdateSubscriptionConfirmation(id uint, confirmed bool) error {
	err = db().Transaction(func(tx *gorm.DB) error {
		var subscription entity.Subscription
		if err := tx.Where(&entity.Subscription{UserID: id}).First(&subscription).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("topic not found")
			}
			return err
		}
		subscription.Confirmed = confirmed
		if err := tx.Save(&subscription).Error; err != nil {
			return err
		}
		return nil
	})
	return err

}

func Migrate() {
	db().AutoMigrate(&entity.Topic{}, &entity.Subscription{}, &entity.Site{}, &entity.User{}, &entity.SubscriptionSchedule{}, &entity.Article{})
}

var (
	dataBase *gorm.DB
	once     sync.Once
	err      error
)

func db() *gorm.DB {
	once.Do(func() {
		dbUser := app.Config.DbUser
		dbPassword := app.Config.DbPassword
		dbHost := app.Config.DbHost
		dbPort := app.Config.DbPort
		dbName := app.Config.DbName
		dbSslMode := app.Config.DbSslMode
		cStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", dbHost, dbUser, dbPassword, dbName, dbPort, dbSslMode)
		dataBase, err = gorm.Open(postgres.Open(cStr), &gorm.Config{Logger: logger.Default.LogMode(getLogMode())})
		if err != nil {
			panic("Unable to connect to db")
		}
		sqlDB, _ := dataBase.DB()
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetMaxIdleConns(5)
	})

	return dataBase
}

func getLogMode() logger.LogLevel {
	switch app.Config.GormLogMode {
	case "error":
		return logger.Error
	case "info":
		return logger.Info
	default:
		return logger.Error
	}
}
