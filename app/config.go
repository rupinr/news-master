package app

import (
	"encoding/json"
	"log/slog"
	"os"
	"reflect"

	"context"
	"log"

	"github.com/joho/godotenv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type EnvVars struct {
	Port                 string `json:"PORT"`
	AdminToken           string `json:"ADMIN_TOKEN"`
	DbPassword           string `json:"DB_PASSWORD"`
	DbHost               string `json:"DB_HOST"`
	DbUser               string `json:"DB_USER"`
	DbName               string `json:"DB_NAME"`
	DbPort               string `json:"DB_PORT"`
	DbSslMode            string `json:"DB_SSL_MODE"`
	GormLogMode          string `json:"GORM_LOG_MODE"`
	PrivateKeyPath       string `json:"PRIVATE_KEY_PATH"`
	PublicKeyPath        string `json:"PUBLIC_KEY_PATH"`
	NewsDataApiUrl       string `json:"NEWS_DATA_API_URL"`
	NewsDataApiKey       string `json:"NEWS_DATA_API_KEY"`
	SubscriptionMailCron string `json:"SUBSCRIPTION_MAIL_CRON"`
	NewsFetchCron        string `json:"NEWS_FETCH_CRON"`
	MaxLoginAttempt      string `json:"MAX_LOGIN_ATTEMPT"`
	AllowOrigin          string `json:"ALLOW_ORIGIN"`
	GinMode              string `json:"GIN_MODE"`
	EmailSender          string `json:"EMAIL_SENDER"`
	SiteUrl              string `json:"SITE_URL"`
	EmailSimulatorMode   string `json:"EMAIL_SIMULATOR_MODE"`
	AdminEmail           string `json:"ADMIN_EMAIL"`
}

var Config EnvVars

func Load() {

	if envFileExists() {
		slog.Info(".env.development exists, Running in Dev mode")
		loadFromEnvFile()
	} else {
		slog.Info("Running in Production mode")
		loadFromAws()
	}

}

func envFileExists() bool {
	info, err := os.Stat(".env.development")
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func loadFromEnvFile() {
	godotenv.Load(".env.development")
	t := reflect.TypeOf(Config)
	v := reflect.ValueOf(&Config).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		f := v.FieldByName(field.Name)
		f.SetString(os.Getenv(field.Tag.Get("json")))
	}
}

func loadFromAws() {
	const secretName = "quick-brew-secrets"
	const region = "eu-central-1"

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatal(err.Error())
	}

	var secretString string = *result.SecretString

	json.Unmarshal([]byte(secretString), &Config)

}
