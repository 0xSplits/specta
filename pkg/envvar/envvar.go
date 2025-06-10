package envvar

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	Environment string `split_words:"true" required:"true"`
	HttpHost    string `split_words:"true" required:"true"`
	HttpPort    string `split_words:"true" required:"true"`
	LogLevel    string `split_words:"true" required:"true"`
}

func Load(kin string) Env {
	var err error

	var env Env

	for {
		{
			err = godotenv.Load(fmt.Sprintf(".env.%s", kin))
			if err != nil {
				fmt.Printf("could not load %s (%s)\n", kin, err)
				time.Sleep(5 * time.Second)
				continue
			}
		}

		{
			err = envconfig.Process("SPECTA", &env)
			if err != nil {
				fmt.Printf("could not process envfile %s (%s)\n", kin, err)
				time.Sleep(5 * time.Second)
				continue
			}
		}

		return env
	}
}
