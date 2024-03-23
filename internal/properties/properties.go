package properties

import (
	"log"
	"time"

	env "github.com/Netflix/go-env"
	"github.com/guilhermealvess/guicpay/internal/token"
	_ "github.com/joho/godotenv/autoload"
)

type props struct {
	TransactionTimeout     time.Duration `env:"TRANSACTION_TIMEOUT,default=5s"`
	RestPort               int           `env:"APP_PORT,default=3000"`
	GRPCPort               int           `env:"GRPC_PORT,default=5000"`
	RedisAddress           string        `env:"REDIS_ADDRESS"`
	AuthorizeServiceURL    string        `env:"AUTHORIZE_SERVICE_URL"`
	NotificationServiceURL string        `env:"NOTIFICATION_SERVICE_URL"`
	SnapshotWalletSize     int           `env:"SNAPSHOT_WALLET_SIZE,default=10"`
	DatabaseURL            string        `env:"DATABASE_URL"`
	JWT                    struct {
		Secret string        `env:"JWT_SECRET"`
		Expire time.Duration `env:"JWT_TOKEN_EXPIRE,default=3600s"`
	}
	TraceCollectorURL string `env:"TRACE_COLLECTOR_URL"`
	DatabaseMaxConn   int    `env:"DATABASE_MAX_CONN,default=15"`
	DatabaseMaxIdle   int    `env:"DATABASE_MAX_IDLE,default=15"`
}

var Props props

func init() {
	if _, err := env.UnmarshalFromEnviron(&Props); err != nil {
		log.Fatal(err)
	}

	token.InitJWT(Props.JWT.Secret, Props.JWT.Expire)
}
