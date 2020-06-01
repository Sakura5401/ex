package utils

import (
	"os"
	"strconv"
)

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)

	if ok {
		return value
	}

	return fallback
}

type Opts struct {
	RabbitmqDsn string
	WorkerCount int
	PostgresURL string
}

func GetEnvs() Opts {
	workerCount := getEnv("WORKER_COUNT", "1")
	workers, err := strconv.Atoi(workerCount)

	if err != nil {
		workers = 1
	}

	return Opts{
		RabbitmqDsn: getEnv("AMQP_URL", "amqp://guest:guest@localhost:5672/"),
		WorkerCount: workers,
		PostgresURL: os.Getenv("DB_URL"),
	}
}
