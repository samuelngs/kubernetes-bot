package client

import (
	"time"

	"github.com/samuelngs/env-go"
	"k8s.io/kubernetes/pkg/client/restclient"
)

// Options struct
type Options struct {
	Host, Username, Password string
	Insecure                 bool
	Interval                 time.Duration
}

func newOptions(opts ...Option) *Options {
	o := new(Options)
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// RestOpt option
func (v *Options) RestOpt() *restclient.Config {
	return &restclient.Config{
		Host:     v.Host,
		Username: v.Username,
		Password: v.Password,
		Insecure: v.Insecure,
	}
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
		env.Raw("K8S_HOST", "https://localhost:443"),
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
		env.Raw("K8S_USER"),
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
		env.Raw("K8S_PASS"),
	)
}

// Insecure option
func Insecure(v bool) Option {
	return func(o *Options) {
		o.Insecure = v
	}
}

// EnvInsecure environent option
func EnvInsecure() Option {
	return Insecure(
		env.Bool("K8S_INSECURE", true),
	)
}

// Interval option
func Interval(i time.Duration) Option {
	return func(o *Options) {
		o.Interval = i
	}
}

// EnvInterval environment option
func EnvInterval() Option {
	return func(o *Options) {
		t := env.I64("K8S_INTERVAL", int64(time.Minute)*30)
		o.Interval = time.Duration(t)
	}
}
