package main

import (
	"context"
	"os"

	"github.com/asecurityteam/awsconfig-filterd/pkg/filter"
	v1 "github.com/asecurityteam/awsconfig-filterd/pkg/handlers/v1"
	"github.com/asecurityteam/runhttp"
	"github.com/asecurityteam/serverfull"
	"github.com/asecurityteam/settings"
)

func main() {
	source, err := settings.NewEnvSource(os.Environ())
	if err != nil {
		panic(err.Error())
	}
	ctx := context.Background()
	resourceFiltererComponent := &filter.ResourceTypeFiltererComponent{}
	resourceTypeFilterer := new(filter.ResourceTypeFilterer)
	err = settings.NewComponent(ctx, source, resourceFiltererComponent, resourceTypeFilterer)
	if err != nil {
		panic(err.Error())
	}
	configFilterHandler := v1.ConfigFilterHandler{
		LogFn:          runhttp.LoggerFromContext,
		StatFn:         runhttp.StatFromContext,
		ConfigFilterer: resourceTypeFilterer,
	}
	handlers := map[string]serverfull.Function{
		"filter": serverfull.NewFunction(configFilterHandler.Handle),
	}

	fetcher := &serverfull.StaticFetcher{Functions: handlers}
	if err := serverfull.Start(ctx, source, fetcher); err != nil {
		panic(err.Error())
	}
}
