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

func createTransport(t *testing.T) *tcp.TcpTransport {
	t.Helper()

	tcp, err := tcp.NewTCPTransport()
	require.NoError(t, err)
	return tcp
}

func createTransportWithOptions(t *testing.T, options tcp.TcpTransportOptions) *tcp.TcpTransport {
	t.Helper()

	tcp, err := tcp.NewTCPTransportWithOptions(options)
	require.NoError(t, err)
	return tcp
}

func createListener(t *testing.T, tcp *tcp.TcpTransport) transport.Listener {
	addr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/0")
	require.NoError(t, err)

	l, err := tcp.Listen(addr)
	require.NoError(t, err)
	return l
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
	tcp := createTransport(t)
	l := createListener(t, tcp)
	defer l.Close()

	cconn, err := tcp.Dial(context.Background(), l.Multiaddr(), "")
	require.NoError(t, err)

	sconn, err := l.Accept()
	require.NoError(t, err)

	testConn(t, cconn, sconn)
}

func TestAcceptMultipleConns(t *testing.T) {
	tcp := createTransport(t)
	l := createListener(t, tcp)
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
	var timeout = 100 * time.Millisecond
	options := tcp.DefaultTcpTransportOptions
	options.AcceptTimeout = timeout

	tcp := createTransportWithOptions(t, options)
	l := createListener(t, tcp)
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

// Check that upgrade error is correctly returned to the dialer and listen
func TestFailedUpgradeOnListen(t *testing.T) {
	ttcp := createTransport(t)
	l := createListener(t, ttcp)
	defer l.Close()

	errCh := make(chan error)
	go func() {
		_, err := l.Accept()
		errCh <- err
	}()

	options := tcp.DefaultTcpTransportOptions
	options.StreamSupported = "testErrorMuxer"
	tcp := createTransportWithOptions(t, options)

	_, err := tcp.Dial(context.Background(), l.Multiaddr(), "dialer")
	l.Close()

	// Check dialer error
	require.ErrorContains(t, "protocols not supported: [testErrorMuxer]", err)
	// Check listener error
	require.ErrorContains(t, "closed", <-errCh)
}

func TestListenerClose(t *testing.T) {
	tcp := createTransport(t)
	l := createListener(t, tcp)

	errCh := make(chan error)
	defer close(errCh)
	go func() {
		_, err := l.Accept()
		errCh <- err
	}()

	select {
	case err := <-errCh:
		t.Fatalf("connection closed earlier than expected. expected nothing on channel, got: %v", err)
	case <-time.After(200 * time.Millisecond):
		// nothing in 200ms.
	}

	// unblocks Accept when it is closed.
	require.NoError(t, l.Close())
	require.ErrorContains(t, "use of closed network connection", <-errCh)

	// doesn't accept new connections when it is closed
	_, err := tcp.Dial(context.Background(), l.Multiaddr(), "")
	require.ErrorContains(t, "refused", err)
}

func TestListenerCloseClosesQueued(t *testing.T) {
	tcp := createTransport(t)
	l := createListener(t, tcp)

	var conns []transport.UpgradedConn
	for i := 0; i < 10; i++ {
		conn, err := tcp.Dial(context.Background(), l.Multiaddr(), "")
		require.NoError(t, err)
		conns = append(conns, conn)
	}

	// wait for all the dials to happen.
	time.Sleep(100 * time.Millisecond)

	// all the connections are opened.
	for _, c := range conns {
		require.Equal(t, false, c.IsClosed())
	}

	// expect that all the connections will be closed.
	err := l.Close()
	require.NoError(t, err)

	for _, c := range conns {
		_ = c.Close()
	}
	// wait for all the connections to close.
	time.Sleep(500 * time.Millisecond)
	for _, c := range conns {
		require.Equal(t, true, c.IsClosed())
	}
}
