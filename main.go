package main

import (
	"context"
	"os"

	"github.com/asecurityteam/awsconfig-filterd/pkg/filter"
	"github.com/asecurityteam/awsconfig-filterd/pkg/handlers/v1"
	"github.com/asecurityteam/runhttp"
	serverfull "github.com/asecurityteam/serverfull/pkg"
	serverfulldomain "github.com/asecurityteam/serverfull/pkg/domain"
	"github.com/asecurityteam/settings"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	configFilterHandler := v1.ConfigFilterHandler{
		LogFn:  runhttp.LoggerFromContext,
		StatFn: runhttp.StatFromContext,
		Filter: &filter.ResourceTypeFilter{},
	}
	handlers := map[string]serverfulldomain.Handler{
		"filter": lambda.NewHandler(configFilterHandler.Handle),
	}

	source, err := settings.NewEnvSource(os.Environ())
	if err != nil {
		panic(err.Error())
	}
	ctx := context.Background()
	rt, err := serverfull.NewStatic(ctx, source, handlers)
	if err != nil {
		panic(err.Error())
	}
	if err := rt.Run(); err != nil {
		panic(err.Error())
	}
}
