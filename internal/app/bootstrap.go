package app

func Bootstrap() (*App, func()) {
	app, cleanup, err := WireApp()
	if err != nil {
		app.stopOnError("boostrap app", err)
	}

	app.logger.Info("Setting timezone")
	app.SetTimezone()

	app.logger.Info("Setting Jaeger tracing")
	app.JaegerTracer()

	return app, cleanup
}
