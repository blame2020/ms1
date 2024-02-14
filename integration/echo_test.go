//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"ms1/pbgen/echopb"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestEcho(t *testing.T) {
	tests := map[string]struct {
		req     *echopb.EchoRequest
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

	ip := ServiceIP(t, "ms1-server")

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			conn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, 50051), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

			require.NoError(t, err)

			defer conn.Close()

			got, err := echopb.NewEchoServiceClient(conn).Echo(context.Background(), tt.req)

			assert.Equal(t, err != nil, tt.wantErr)

			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("unexpected diff:\n%v", diff)
			}
		})
	}
}
