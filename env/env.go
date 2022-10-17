package env

import "github.com/caarlos0/env"

func Load() (e *Env, err error) {
	e = new(Env)
	if err := env.Parse(e); err != nil {
		return nil, err
	}
	return e, nil
}

type Env struct {
	DatabaseURL string `env:"DATABASE_URL,required"`
	CSRFToken   string `env:"CSRF_TOKEN,required"`
	SessionKey  string `env:"SESSION_KEY,required"`
}
