package di

import (
	"github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/providers"

	"go.uber.org/fx"
)

var ProviderModule = fx.Module(
	"providers",
	fx.Provide(
		providers.NewRedisProvider,
		providers.NewKafkaProvider,
		providers.NewValidator,
	),

	// Tracer is invoked because it does not consider as dependency, but rather as global provider
	// and if no part of application needs it as dependency, fx won't provide it to anything, so it needs to be invoked
	// once the app is started
	fx.Invoke(
		providers.NewLoggerProvider,
		providers.NewTraceProvider,
	),
)
