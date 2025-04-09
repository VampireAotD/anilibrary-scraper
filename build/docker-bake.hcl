group "default" {
    targets = ["app", "kafka"]
}

variable "TIMEZONE" {
    default = "Europe/Kiev"
}

target "app" {
    tags = ["anilibrary-scraper:latest"]
    dockerfile = "build/app/Dockerfile"
    target = "final"
    args = {
        TIMEZONE = "${TIMEZONE}"
    }
}

target "kafka" {
    tags = ["anilibrary-scraper-kafka:latest"]
    context = "build/kafka"
    dockerfile = "Dockerfile"
}
