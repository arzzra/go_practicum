package main

import (
	"agent/internal/agent"
	"context"
	"os/signal"
	"syscall"
	"time"
)

const pollInterval = 2 * time.Second
const reportInterval = 10 * time.Second
const requesTimeout = 2 * time.Second
const serverHost = "127.0.0.1"
const serverPort = "8080"

var Metrics = []string{"Alloc", "BuckHashSys", "Frees",
	"GCCPUFraction", "GCSys", "HeapAlloc",
	"HeapIdle", "HeapInuse", "HeapObjects",
	"HeapReleased", "HeapSys", "LastGC",
	"Lookups", "MCacheInuse", "MCacheSys",
	"MSpanInuse", "MSpanSys", "Mallocs",
	"NextGC", "NumForcedGC", "NumGC",
	"OtherSys", "PauseTotalNs", "StackInuse",
	"StackSys", "Sys", "TotalAlloc"}

func initSettings() agent.Settings {
	settings := agent.Settings{
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
		Metrics:        nil,
		Host:           serverHost,
		Port:           serverPort,
		RequestTimeout: requesTimeout,
	}
	arrayString := make([]string, len(Metrics))
	copy(arrayString, Metrics)
	settings.Metrics = &arrayString
	return settings
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()
	agent := agent.MakeAgent(initSettings())
	agent.Start(ctx)
}
