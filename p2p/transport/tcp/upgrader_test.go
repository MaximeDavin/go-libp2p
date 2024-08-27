package tcp_test

import (
	"context"
	"testing"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"

	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"github.com/prysmaticlabs/prysm/v5/testing/require"
)

func TestIncompatibleOutboundSecurity(t *testing.T) {
	ltcp := createTransport(t)
	l := createListener(t, ltcp)
	defer l.Close()
	go func() {
		l.Accept()
	}()

	options := tcp.DefaultTcpTransportOptions
	options.SecuritySupported = "/invalidSecurity/1.0"
	dtcp := createTransportWithOptions(t, options)
	conn, err := manet.Dial(l.Multiaddr())
	require.NoError(t, err)

	_, err = tcp.Upgrade(context.Background(), dtcp, conn, "", network.DirOutbound)
	require.ErrorContains(t, "/invalidSecurity/1.0", err)
}

func TestIncompatibleOutboundStream(t *testing.T) {
	ltcp := createTransport(t)
	l := createListener(t, ltcp)
	defer l.Close()
	go func() {
		l.Accept()
	}()

	options := tcp.DefaultTcpTransportOptions
	options.StreamSupported = "/invalidStream/1.0"
	dtcp := createTransportWithOptions(t, options)
	conn, err := manet.Dial(l.Multiaddr())
	require.NoError(t, err)

	_, err = tcp.Upgrade(context.Background(), dtcp, conn, "", network.DirOutbound)
	require.ErrorContains(t, "/invalidStream/1.0", err)
}

// TODO(max): fix this test. Upgrade does not return when protocol negociation
// fails for an inbound connection. This could mean that each time a connection
// is accepted with unsupported security or stream we have a go routine that
// wait for ever ?
func TestIncompatibleInboundSecurity(t *testing.T) {
	addr, _ := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/0")
	l, _ := manet.Listen(addr)
	defer l.Close()

	dtcp := createTransport(t)
	go func() {
		dtcp.Dial(context.Background(), l.Multiaddr(), "")
	}()

	conn, err := l.Accept()
	require.NoError(t, err)

	options := tcp.DefaultTcpTransportOptions
	options.SecuritySupported = "/invalid/1.0"
	ltcp := createTransportWithOptions(t, options)
	_, err = tcp.Upgrade(context.Background(), ltcp, conn, "", network.DirInbound)
	require.ErrorContains(t, "/invalid/1.0", err)
}

// TODO(max): replicate TestIncompatibleInboundSecurity for streaming once fixed

func TestUpgradeSuccessfullyOutbound(t *testing.T) {
	ltcp := createTransport(t)
	l := createListener(t, ltcp)
	defer l.Close()
	go func() {
		l.Accept()
	}()

	dtcp := createTransport(t)
	conn, err := manet.Dial(l.Multiaddr())
	require.NoError(t, err)

	_, err = tcp.Upgrade(context.Background(), dtcp, conn, "", network.DirOutbound)
	require.NoError(t, err)
}

func TestUpgradeSuccessfullyInbound(t *testing.T) {
	addr, _ := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/0")
	l, _ := manet.Listen(addr)
	defer l.Close()

	dtcp := createTransport(t)
	go func() {
		dtcp.Dial(context.Background(), l.Multiaddr(), "")
	}()

	conn, err := l.Accept()
	require.NoError(t, err)

	ltcp := createTransport(t)
	_, err = tcp.Upgrade(context.Background(), ltcp, conn, "", network.DirInbound)
	require.NoError(t, err)
}
