package main

import (
	"context"
	"ms1/pbgen/echopb"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestNewEchoServiceServer(t *testing.T) {
	tests := map[string]struct {
		want    *EchoServiceServer
		wantErr bool
	}{
		"basic": {&EchoServiceServer{}, false},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := NewEchoServiceServer()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEchoServiceServer_Echo(t *testing.T) {
	tests := map[string]struct {
		in      *echopb.EchoRequest
		want    *echopb.EchoResponse
		wantErr bool
	}{
		"basic": {
			&echopb.EchoRequest{
				Message: "echo",
			},
			&echopb.EchoResponse{
				Message: "echo",
			},
			false,
		},
	}

	s := grpc.NewServer()
	srv := NewEchoServiceServer()
	echopb.RegisterEchoServiceServer(s, srv)

	lis := bufconn.Listen(1024 * 1024)
	errs := make(chan error, 1)

	go func() { errs <- s.Serve(lis) }()

	t.Cleanup(func() {
		s.Stop()

		_ = lis.Close()

		require.NoError(t, <-errs)
	})

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			conn, err := grpc.DialContext(
				context.Background(),
				"bufnet",
				grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return lis.Dial() }),
				grpc.WithTransportCredentials(insecure.NewCredentials()))

			require.NoError(t, err)

			defer conn.Close()

			got, err := echopb.NewEchoServiceClient(conn).Echo(context.Background(), tt.in)

			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("unexpected diff:\n%v", diff)
			}

			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}
