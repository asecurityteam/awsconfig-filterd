package main

import (
	"context"
	"os"

	"github.com/asecurityteam/awsconfig-filterd/pkg/decorators"
	"github.com/asecurityteam/awsconfig-filterd/pkg/filter"
	v1 "github.com/asecurityteam/awsconfig-filterd/pkg/handlers/v1"
	producer "github.com/asecurityteam/component-producer"
	"github.com/asecurityteam/runhttp"
	"github.com/asecurityteam/serverfull"
	"github.com/asecurityteam/settings"
)

type config struct {
	Decorators *decorators.ChainConfig
	Filter     *filter.FilterConfig
	Producer   *producer.Config
	LambdaMode bool `description:"Use the Lambda SDK to start the system."`
}

func (*config) Name() string {
	return "awsconfigfilterd"
}

type component struct {
	Decorators *decorators.ChainComponent
	Filter     *filter.FilterComponent
	Producer   *producer.Component
}

func newComponent() *component {
	return &component{
		Decorators: decorators.NewChainComponent(),
		Filter:     filter.NewFilterComponent(),
		Producer:   producer.NewComponent(),
	}
}

func (c *component) Settings() *config {
	return &config{
		Decorators: c.Decorators.Settings(),
		Filter:     c.Filter.Settings(),
		Producer:   c.Producer.Settings(),
	}
}

func (c *component) New(ctx context.Context, conf *config) (func(context.Context, settings.Source) error, error) {
	f, err := c.Filter.New(ctx, conf.Filter)
	if err != nil {
		return nil, err
	}
	p, err := c.Producer.New(ctx, conf.Producer)
	if err != nil {
		return nil, err
	}
	chain, err := c.Decorators.New(ctx, conf.Decorators)
	if err != nil {
		return nil, err
	}

	filterHandler := &v1.ConfigFilter{
		LogFn:          runhttp.LoggerFromContext,
		StatFn:         runhttp.StatFromContext,
		ConfigFilterer: f,
		Producer:       p,
	}
	decoratedFilter := chain.Decorate(filterHandler.Handle)
	handlers := map[string]serverfull.Function{
		"filter": serverfull.NewFunction(decoratedFilter),
	}
	fetcher := &serverfull.StaticFetcher{Functions: handlers}
	if conf.LambdaMode {
		return func(ctx context.Context, source settings.Source) error {
			return serverfull.StartLambda(ctx, source, fetcher, "filter")
		}, nil
	}
	return func(ctx context.Context, source settings.Source) error {
		return serverfull.StartHTTP(ctx, source, fetcher)
	}, nil
}

func main() {
	source, err := settings.NewEnvSource(os.Environ())
	if err != nil {
		panic(err.Error())
	}
	ctx := context.Background()
	runner := new(func(context.Context, settings.Source) error)
	cmp := newComponent()
	err = settings.NewComponent(ctx, source, cmp, runner)
	if err != nil {
		panic(err.Error())
	}
	if err := (*runner)(ctx, source); err != nil {
		panic(err.Error())
	}
}
