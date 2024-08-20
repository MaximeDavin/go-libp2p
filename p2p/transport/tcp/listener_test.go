package tcp_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/transport"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/prysmaticlabs/prysm/v5/testing/require"
)

func createTransportAndListener(t *testing.T) (*tcp.TcpTransport, transport.Listener) {
	t.Helper()

	tcp, err := tcp.NewTCPTransport()
	require.NoError(t, err)

	addr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/0")
	require.NoError(t, err)

	l, err := tcp.Listen(addr)
	require.NoError(t, err)

	return tcp, l
}

func testConn(t *testing.T, clientConn, serverConn transport.UpgradedConn) {
	t.Helper()

	cstr, err := clientConn.OpenStream(context.Background())
	require.NoError(t, err)

	_, err = cstr.Write([]byte("foobar"))
	require.NoError(t, err)

	sstr, err := serverConn.AcceptStream()
	require.NoError(t, err)

	b := make([]byte, 6)
	_, err = sstr.Read(b)
	require.NoError(t, err)
	require.DeepEqual(t, []byte("foobar"), b)
}

func TestAcceptSingleConn(t *testing.T) {
	tcp, l := createTransportAndListener(t)
	defer l.Close()

	cconn, err := tcp.Dial(context.Background(), l.Multiaddr(), "")
	require.NoError(t, err)

	sconn, err := l.Accept()
	require.NoError(t, err)

	testConn(t, cconn, sconn)
}

func TestAcceptMultipleConns(t *testing.T) {
	tcp, l := createTransportAndListener(t)
	defer l.Close()

	var toClose []io.Closer
	defer func() {
		for _, c := range toClose {
			_ = c.Close()
		}
	}()

	var ctx = context.Background()

	for i := 0; i < 10; i++ {
		cconn, err := tcp.Dial(ctx, l.Multiaddr(), "")
		require.NoError(t, err)
		toClose = append(toClose, cconn)

		sconn, err := l.Accept()
		require.NoError(t, err)
		toClose = append(toClose, sconn)

		testConn(t, cconn, sconn)
	}
}

func TestConnectionsClosedIfNotAccepted(t *testing.T) {
	var timeout = 2 * time.Second

	tcp, l := createTransportAndListener(t)
	defer l.Close()

	var ctx = context.Background()
	conn, err := tcp.Dial(ctx, l.Multiaddr(), "")
	require.NoError(t, err)

	errCh := make(chan error)
	go func() {
		defer conn.Close()
		str, err := conn.OpenStream(context.Background())
		if err != nil {
			errCh <- err
			return
		}
		// start a Read. It will block until the connection is closed
		_, _ = str.Read([]byte{0})
		errCh <- nil
	}()

	time.Sleep(timeout / 2)
	select {
	case err := <-errCh:
		t.Fatalf("connection closed earlier than expected. expected nothing on channel, got: %v", err)
	default:
	}

	time.Sleep(timeout)
	require.NoError(t, <-errCh)
}

func TestFailedUpgradeOnListen(t *testing.T) {
	// id, u := createUpgraderWithMuxers(t, []upgrader.StreamMuxer{{ID: "errorMuxer", Muxer: &errorMuxer{}}}, nil, nil)

	tcp, l := createTransportAndListener(t)

	errCh := make(chan error)
	go func() {
		_, err := l.Accept()
		errCh <- err
	}()

	_, err := tcp.Dial(context.Background(), l.Multiaddr(), "")
	require.ErrorContains(t, "errorMuxer", err)

	// close the listener.
	l.Close()
	require.ErrorIs(t, err, <-errCh)
}
