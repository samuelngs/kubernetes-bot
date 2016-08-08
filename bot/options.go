package bots

import "github.com/samuelngs/env-go"

// Options struct
type Options struct {
	Host, Username, Password, Token, Secret, Channel string
}

func newOptions(opts ...Option) *Options {
	o := new(Options)
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Option func
type Option func(*Options)

// Host option
func Host(s string) Option {
	return func(o *Options) {
		o.Host = s
	}
}

// EnvHost environment option
func EnvHost() Option {
	return Host(
		env.Raw("BOT_HOST", "localhost"),
	)
}

// Username option
func Username(s string) Option {
	return func(o *Options) {
		o.Username = s
	}
}

// EnvUsername option
func EnvUsername() Option {
	return Username(
		env.Raw("BOT_USER"),
	)
}

// Password option
func Password(s string) Option {
	return func(o *Options) {
		o.Password = s
	}
}

// EnvPassword environment option
func EnvPassword() Option {
	return Password(
		env.Raw("BOT_PASS"),
	)
}

// Token option
func Token(s string) Option {
	return func(o *Options) {
		o.Token = s
	}
}

// EnvToken environment option
func EnvToken() Option {
	return Token(
		env.Raw("BOT_TOKEN"),
	)
}

// Secret option
func Secret(s string) Option {
	return func(o *Options) {
		o.Secret = s
	}
}

// EnvSecret environment option
func EnvSecret() Option {
	return Secret(
		env.Raw("BOT_SECRET"),
	)
}

// Channel option
func Channel(s string) Option {
	return func(o *Options) {
		o.Channel = s
	}
}

// EnvChannel environment option
func EnvChannel() Option {
	return Channel(
		env.Raw("BOT_CHANNEL"),
	)
}
