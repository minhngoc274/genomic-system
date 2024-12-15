package tests

import (
	"context"
	"github.com/golibs-starter/golib"
	golibtest "github.com/golibs-starter/golib-test"
	"github.com/minhngoc274/genomic-system/genomic-service/bootstrap"
	"go.uber.org/fx"
)

func init() {
	err := fx.New(
		bootstrap.All(),
		golib.ProvidePropsOption(golib.WithActiveProfiles([]string{"local"})),
		golib.ProvidePropsOption(golib.WithPaths([]string{"../config/"})),
		golibtest.EnableWebTestUtil(),
	).Start(context.Background())
	if err != nil {
		panic(err)
	}
}
