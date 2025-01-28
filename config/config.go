package config

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AWSIAMUser struct {
		AccessKeyID     string `yaml:"access_key_id" envconfig:"AWS_ACCESS_KEY_ID"`
		SecretAccessKey string `yaml:"secret_access_key" envconfig:"AWS_SECRET_ACCESS_KEY"`
		Region          string `yaml:"region" envconfig:"AWS_REGION"`
	} `yaml:"aws_iam_user"`
	JWT struct {
		FrontendSecretKey string `yaml:"frontend_secret_key" envconfig:"JWT_FRONTEND_SECRET_KEY"`
		BackendSecretKey  string `yaml:"backend_secret_key" envconfig:"JWT_BACKEND_SECRET_KEY"`
	} `yaml:"jwt"`
	Database struct {
		File     string `yaml:"file" envconfig:"DB_FILE"`
		Username string `yaml:"sql-user" envconfig:"DB_USERNAME"`
		Password string `yaml:"sql-pass" envconfig:"DB_PASSWORD"`
		Database string `yaml:"sql-database" envconfig:"DB_DATABASE"`
		Server   string `yaml:"sql-server" envconfig:"DB_SERVER"`
		Port     int    `yaml:"sql-port" envconfig:"DB_PORT"`
	} `yaml:"database"`
	S3Bucket   string `yaml:"s3_bucket" envconfig:"S3_BUCKET"`
	UserAPIURL string `yaml:"userapi_url" envconfig:"USERAPI_URL"`
	ClientURL  string `yaml:"client_url" envconfig:"CLIENT_URL"`
}

var AppConfig Config

func LoadConfig() {
	// Ladda miljövariabler från .env-filen i utvecklingsmiljö
	if os.Getenv("RUNENVIRONMENT") != "Production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	readFile(&AppConfig)
	readEnv(&AppConfig)
	readSecrets(&AppConfig)
	fmt.Printf("%+v", AppConfig)
}

func readFile(cfg *Config) {
	fileName := "config/config.yaml"
	s := os.Getenv("RUNENVIRONMENT")
	if len(s) > 0 {
		fileName = "config/config" + s + ".yaml"
	}

	f, err := os.Open(fileName)
	if err != nil {
		log.Printf("Error opening config file: %v", err)
		return
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(cfg); err != nil {
		log.Printf("Error decoding config file: %v", err)
	}
}

func readEnv(cfg *Config) {
	if err := envconfig.Process("", cfg); err != nil {
		log.Fatalf("Error processing environment variables: %v", err)
	}
}

func readSecrets(cfg *Config) {
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if awsAccessKeyID == "" || awsSecretAccessKey == "" {
		log.Fatal("Missing AWS credentials")
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.AWSIAMUser.Region),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	ssmClient := ssm.New(sess)

	// Hämta parametrar från Parameter Store
	params := map[string]*string{
		"/imagebook-imageapi/jwt/frontend_secret_key": &cfg.JWT.FrontendSecretKey,
		"/imagebook-imageapi/jwt/backend_secret_key":  &cfg.JWT.BackendSecretKey,
		"/imagebook-imageapi/database/sql-user":       &cfg.Database.Username,
		"/imagebook-imageapi/database/sql-pass":       &cfg.Database.Password,
		"/imagebook-imageapi/database/sql-database":   &cfg.Database.Database,
		"/imagebook-imageapi/database/sql-server":     &cfg.Database.Server,
	}

	if os.Getenv("RUNENVIRONMENT") == "Production" {
		params["/imagebook-imageapi/s3/bucketnameprod"] = &cfg.S3Bucket
		params["/imagebook-imageapi/userapi_url"] = &cfg.UserAPIURL
		params["/imagebook-imageapi/client_url"] = &cfg.ClientURL
	} else {
		params["/imagebook-imageapi/s3/bucketnamedev"] = &cfg.S3Bucket
	}

	for param, dest := range params {
		result, err := ssmClient.GetParameter(&ssm.GetParameterInput{
			Name:           aws.String(param),
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			log.Fatalf("Failed to get parameter %s: %v", param, err)
		}
		*dest = *result.Parameter.Value
	}
}
