package main

import (
	"log/slog"
	"os"
	"sync"
)

type Server struct {
	Logger     *slog.Logger
	GrpcServer *GrpcServer
}

func (s *Server) Run() error {
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		s.Logger.Info("start")

		if err := s.GrpcServer.ListenAndServe(); err != nil {
			os.Exit(1)
		}
	}()

	wg.Wait()

	return nil
}
