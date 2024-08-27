package testtransport

import (
	"context"
	"fmt"
	"testing"

	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/prysmaticlabs/prysm/v5/testing/require"
)

type testcase struct {
	Name             string
	ServerPreference string
	ClientPreference string
	Expected         string
	Error            string
}

// TODO: Use Host or Network instead of the low level TcpTransport
func TestStreamMuxerNegociation(t *testing.T) {
	testcases := []testcase{
		{
			Name:             "server and client have the same preference",
			ServerPreference: yamux.ID,
			ClientPreference: yamux.ID,
			Expected:         yamux.ID,
		},
		{
			Name:             "server does not support yamux",
			ServerPreference: "unsupported_proto",
			ClientPreference: yamux.ID,
			Error:            "unsupported_proto",
		},
		{
			Name:             "client does not support yamux",
			ServerPreference: yamux.ID,
			ClientPreference: "unsupported_proto",
			Error:            "unsupported_proto",
		},
		{
			Name:             "client and server does not support yamux",
			ServerPreference: "unsupported_proto",
			ClientPreference: "unsupported_proto",
			Error:            fmt.Sprintf("stream multiplexing protocol not supported, only %s is supported", yamux.ID),
		},
	}

	for _, tc := range testcases {
		tc := tc

		t.Run(tc.Name, func(t *testing.T) {
			soptions := tcp.DefaultTcpTransportOptions
			soptions.StreamSupported = tc.ServerPreference
			server, err := tcp.NewTCPTransportWithOptions(soptions)
			require.NoError(t, err)

			coptions := tcp.DefaultTcpTransportOptions
			coptions.StreamSupported = tc.ClientPreference
			client, err := tcp.NewTCPTransportWithOptions(coptions)
			require.NoError(t, err)

			addr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/0")
			require.NoError(t, err)

			l, err := server.Listen(addr)
			require.NoError(t, err)
			defer l.Close()

			go func() {
				l.Accept()
			}()

			_, err = client.Dial(context.Background(), l.Multiaddr(), "")

			if tc.Error != "" {
				require.ErrorContains(t, tc.Error, err)
			} else {
				require.NoError(t, err)
				// TODO: Find a way to check if the connection has the correct muxer
			}
		})
	}
}
