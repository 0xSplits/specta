package envvar

import (
	"fmt"
	"slices"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/xh3b4sd/tracer"
)

type Env struct {
	Environment string `split_words:"true" required:"true"`
	HttpHost    string `split_words:"true" required:"true"`
	HttpPort    string `split_words:"true" required:"true"`
	LogLevel    string `split_words:"true" required:"true"`
}

func Load(pat string) Env {
	var err error
	var env Env

	{
		err = godotenv.Load(pat)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	{
		err = envconfig.Process("SPECTA", &env)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}

	{
		if !slices.Contains([]string{"development", "testing", "staging", "production"}, env.Environment) {
			tracer.Panic(tracer.Mask(fmt.Errorf("SPECTA_ENVIRONMENT must be one of [development testing staging production]")))
		}
	}

	return env
}
