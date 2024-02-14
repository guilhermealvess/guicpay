package properties

import (
	"log"
	"time"

	_ "github.com/joho/godotenv/autoload"
	env "github.com/Netflix/go-env"
)

type props struct {
	TransactionTimeout time.Duration `env:"TRANSACTION_TIMEOUT,default=5s"`
}

var Props props

func init() {
	if _, err := env.UnmarshalFromEnviron(&Props); err != nil {
		log.Fatal(err)
	}
}
