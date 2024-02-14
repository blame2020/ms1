package main

import (
	"ms1/pbgen/echopb"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
)

type GrpcAddr = string

type GrpcServer struct {
	Addr              GrpcAddr
	EchoServiceServer echopb.EchoServiceServer
}

func (gs *GrpcServer) ListenAndServe() error {
	s := grpc.NewServer()
	echopb.RegisterEchoServiceServer(s, gs.EchoServiceServer)

	lis, err := net.Listen("tcp", gs.Addr)
	if err != nil {
		return err
	}

	cancel := make(chan struct{})
	defer func() { close(cancel) }()

	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)

	go func() {
		defer wg.Done()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

		select {
		case <-sig:
			s.GracefulStop()

			return
		case <-cancel:
			return
		}
	}()

	return s.Serve(lis)
}
