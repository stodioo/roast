package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"module": "main"})

var (
	revisionID     = "unknown"
	buildTimestamp = "unknown"
)

func main() {
	_, err := fmt.Fprintf(os.Stdout, "Blog service revisionID %s, built at %s", revisionID, buildTimestamp)

	if err != nil {
		panic(err)
	}

	for _, flag := range os.Args[1:] {
		if flag == "--version" {
			return
		}
	}

	logLevelStr := os.Getenv("LOG_LEVEL")
	if logLevelStr != "" {
		logLevel, err := logrus.ParseLevel(logLevelStr)
		if err != nil {
			panic(err)
		}

		logrus.SetLevel(logLevel)
	}

	apiServicePort := os.Getenv("API_PORT")
	if apiServicePort == "" {
		apiServicePort = "8080"
	}

	// Initiating router...
	log.Infof("Initializing router")
	router := mux.NewRouter()

	router.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("OK"))
	})

	var shuttingDown bool
	shutdownSignal := make(chan os.Signal)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	apiServer := http.Server{
		Addr:    ":" + apiServicePort,
		Handler: handleCORS(router, ""),
	}

	go func() {
		log.Infof("Service is ready")
		if err := apiServer.ListenAndServe(); err != nil && (err != http.ErrServerClosed || !shuttingDown) {
			log.Fatalf("REST Api Server error: %v", err)
		}
	}()

	<-shutdownSignal
	shuttingDown = true
	log.Infof("Shutting down the server")

	go func() {
		<-shutdownSignal
		os.Exit(0)
	}()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	apiServer.Shutdown(shutdownCtx)
	log.Infof("Done.")
}

func handleCORS(router *mux.Router, envPrefix string) http.Handler {
	var corsAllowedHeaders []string
	if strVal, specified := os.LookupEnv(envPrefix + "CORS_ALLOWED_HEADERS"); specified {
		if strVal != "" {
			parts := strings.Split(strVal, ",")
			for _, str := range parts {
				corsAllowedHeaders = append(corsAllowedHeaders, strings.TrimSpace(str))
			}
		}
	} else {
		corsAllowedHeaders = []string{"Content-Type", "Accept", "Authorization", "x-api-key"}
	}

	var corsAllowedMethods []string
	if strVal := os.Getenv(envPrefix + "CORS_ALLOWED_METHODS"); strVal != "" {
		parts := strings.Split(strVal, ",")
		for _, str := range parts {
			corsAllowedMethods = append(corsAllowedMethods, strings.TrimSpace(str))
		}
	}
	var corsAllowedDomains []string
	if strVal := os.Getenv(envPrefix + "CORS_ALLOWED_DOMAINS"); strVal != "" {
		parts := strings.Split(strVal, ",")
		for _, str := range parts {
			corsAllowedDomains = append(corsAllowedDomains, strings.TrimSpace(str))
		}
	}
	allowedHeaders := handlers.AllowedHeaders(corsAllowedHeaders)
	allowedMethods := handlers.AllowedMethods(corsAllowedMethods)
	allowedDomains := handlers.AllowedOrigins(corsAllowedDomains)
	return handlers.CORS(allowedDomains, allowedHeaders, allowedMethods)(router)
}
