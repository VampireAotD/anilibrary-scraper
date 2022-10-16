package app

func Bootstrap() *App {
	var app App

	app.InitLogger()
	app.logger.Info("Initialized logger")

	app.logger.Info("Reading config")
	app.ReadConfig()

	app.logger.Info("Setting timezone")
	app.SetTimezone()

	app.logger.Info("Setting redis connection")
	app.SetRedisConnection()

	app.logger.Info("Setting Jaeger tracing")
	app.JaegerTracer()

	return &app
}
