package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/KadimovRus/calc_go/internal/application"
)

func main() {
	agent, err := application.NewAgent()
	if err != nil {
		log.Fatal("Ошибка создания агента:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		log.Println("Агент запущен")
		<-sigChan
		cancel()
	}()

	if err := agent.Run(ctx); err != nil {
		log.Fatal("Ошибка запуска агента:", err)
	}
}
