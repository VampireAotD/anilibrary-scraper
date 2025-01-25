package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "scraper_requests_counter",
		Help: "Counter of HTTP requests made by scraper",
	})

	failedRequestCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "scraper_failed_requests_counter",
		Help: "Counter of failed HTTP requests made by scraper",
	})

	failedImageScrapeCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "scraper_failed_image_scrape_counter",
		Help: "Counter of failed image scrapes",
	})
)

func IncrScraperRequestCounter() {
	requestCounter.Inc()
}

func IncrScraperFailedRequestCounter() {
	failedRequestCounter.Inc()
}

func IncrScraperFailedImageScrapeCounter() {
	failedImageScrapeCounter.Inc()
}
