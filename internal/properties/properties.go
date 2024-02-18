package properties

import (
	"log"
	"time"

	env "github.com/Netflix/go-env"
	_ "github.com/joho/godotenv/autoload"
)

type props struct {
	TransactionTimeout     time.Duration `env:"TRANSACTION_TIMEOUT,default=5s"`
	Port                   int           `env:"APP_PORT,default=3000"`
	RedisAddress           string        `env:"REDIS_ADDRESS"`
	AuthorizeServiceURL    string        `env:"AUTHORIZE_SERVICE_URL"`
	NotificationServiceURL string        `env:"NOTIFICATION_SERVICE_URL"`
}

var Props props

func init() {
	if _, err := env.UnmarshalFromEnviron(&Props); err != nil {
		log.Fatal(err)
	}
}
