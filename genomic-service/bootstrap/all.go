package bootstrap

import (
	"github.com/golibs-starter/golib"
	golibcron "github.com/golibs-starter/golib-cron"
	golibgin "github.com/golibs-starter/golib-gin"
	golibsec "github.com/golibs-starter/golib-security"
	"github.com/minhngoc274/genomic-system/genomic-service/adapters/blockchain"
	"github.com/minhngoc274/genomic-system/genomic-service/adapters/repositories"
	"github.com/minhngoc274/genomic-system/genomic-service/adapters/tee"
	"github.com/minhngoc274/genomic-system/genomic-service/config"
	"github.com/minhngoc274/genomic-system/genomic-service/controllers"
	"github.com/minhngoc274/genomic-system/genomic-service/jobs"
	"github.com/minhngoc274/genomic-system/genomic-service/routers"
	"go.uber.org/fx"
)

// All register all constructors for fx container
func All() fx.Option {
	return fx.Options(
		golib.AppOpt(),
		golib.PropertiesOpt(),
		golib.LoggingOpt(),
		golib.EventOpt(),
		golib.BuildInfoOpt(Version, CommitHash, BuildTime),
		golib.ActuatorEndpointOpt(),
		golibgin.GinHttpServerOpt(),
		fx.Invoke(routers.RegisterHandlers),
		fx.Invoke(routers.RegisterGinRouters),
		golibcron.Opt(),

		golibsec.BasicAuthOpt(),

		golib.ProvideProps(config.NewChainProperties),

		golibcron.ProvideJob(jobs.NewFilterLogsJob),

		fx.Provide(tee.NewTEEService),
		fx.Provide(blockchain.NewBlockchainService),

		fx.Provide(repositories.NewUserRepository),
		fx.Provide(repositories.NewGeneticDataRepository),

		fx.Provide(controllers.NewUserController),


		// Graceful shutdown.
		// OnStop hooks will run in reverse order.
		golib.OnStopEventOpt(),
		golibgin.OnStopHttpServerOpt(),
	)
}
