// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package app

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/axetroy/wsm/internal/app/config"
	"github.com/axetroy/wsm/internal/app/db"
)

func Serve() error {
	port := config.Http.Port

	s := &http.Server{
		Addr:           ":" + port,
		Handler:        UserRouter,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 10M
	}

	go func() {
		if config.Http.TLS != nil {
			TLSConfig := &tls.Config{
				MinVersion:               tls.VersionTLS11,
				CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
				PreferServerCipherSuites: true,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
					tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				},
			}

			TLSProto := make(map[string]func(*http.Server, *tls.Conn, http.Handler))

			s.TLSConfig = TLSConfig
			s.TLSNextProto = TLSProto

			log.Printf("Listen on:  %s\n", s.Addr)

			if err := s.ListenAndServeTLS(config.Http.TLS.Cert, config.Http.TLS.Key); err != nil {
				log.Fatalln(err)
			}
		} else {
			log.Printf("Listen on:  %s\n", s.Addr)

			if err := s.ListenAndServe(); err != nil {
				log.Fatalln(err)
			}
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	config.Common.Exiting = true

	log.Println("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("Timeout of 1 seconds.")
	}

	_ = db.Db.Close()

	log.Println("Server exiting")

	return nil
}
