package main

import (

	"log"
	"net/http"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"
	"context"
	
	"http-server-projeto-korp/internal/handler"

)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/projeto-korp", handler.ProjetoKorpHandler)
	mux.HandleFunc("/health", handler.HealthCheckHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
	
	go func() {
		
		log.Println("Capsule Korp - Iniciando Servidor na Porta :8080")
		
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Epic Fail :8080: %v\n", err)
		
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Println("Capsule Korp - Encerrando Servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Epic Fail ao Encerrar Servidor: %v\n", err)
	}
	
	log.Println("Capsule Korp - Servidor Encerrado com Sucesso")

}