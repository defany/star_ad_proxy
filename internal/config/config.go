package config

import (
	"log"
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var cfg = koanf.New(".")

func MustLoad() {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		if err := cfg.Load(file.Provider(path), json.Parser()); err != nil {
			log.Fatalf("error loading config: %v", err)
		}
	}

	err := cfg.Load(env.Provider(".", env.Opt{
		TransformFunc: func(k, v string) (string, any) {
			k = strings.ReplaceAll(strings.ToLower(k), "_", ".")

			if strings.Contains(v, " ") {
				return k, strings.Split(v, " ")
			}

			return k, v
		},
	}), nil)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}
}

func IsProduction() bool {
	return cfg.String("env") == "production"
}

func IsDevelopment() bool {
	return cfg.String("env") == "dev"
}
