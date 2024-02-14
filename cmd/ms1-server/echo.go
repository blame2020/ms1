package main

import (
	"context"
	"ms1/pbgen/echopb"
)

type EchoServiceServer struct {
	echopb.UnimplementedEchoServiceServer
}

func NewEchoServiceServer() *EchoServiceServer {
	return &EchoServiceServer{}
}

var _ echopb.EchoServiceServer = (*EchoServiceServer)(nil)

func (p *EchoServiceServer) Echo(_ context.Context, in *echopb.EchoRequest) (*echopb.EchoResponse, error) {
	return &echopb.EchoResponse{
		Message: in.Message,
	}, nil
}
