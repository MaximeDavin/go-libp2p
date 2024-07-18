package transport

import (
	"bytes"
	"context"
	"io"
	"testing"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/prysmaticlabs/prysm/v5/testing/require"
)

func TestTcpTransport(t *testing.T) {
	for i := 0; i < 2; i++ {
		ta, err := NewTCPTransport()
		require.NoError(t, err)
		tb, err := NewTCPTransport()
		require.NoError(t, err)

		zero := "/ip4/127.0.0.1/tcp/0"
		maddr, err := ma.NewMultiaddr(zero)
		if err != nil {
			t.Fatal(err)
		}

		ttransport.SubtestTransport(t, ta, tb, maddr, "peerA")

		envReuseportVal = false
	}
	envReuseportVal = true
}

func SubtestBasic(t *testing.T, ta, tb Transport, maddr ma.Multiaddr, peerA string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	list, err := ta.Listen(maddr)
	if err != nil {
		t.Fatal(err)
	}
	defer list.Close()

	var (
		connA, connB StreamConn
		done         = make(chan struct{})
	)
	defer func() {
		<-done
		if connA != nil {
			connA.Close()
		}
		if connB != nil {
			connB.Close()
		}
	}()

	go func() {
		defer close(done)
		var err error
		connB, err = list.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		s, err := connB.AcceptStream()
		if err != nil {
			t.Error(err)
			return
		}

		buf, err := io.ReadAll(s)
		if err != nil {
			t.Error(err)
			return
		}

		if !bytes.Equal(testData, buf) {
			t.Errorf("expected %s, got %s", testData, buf)
		}

		n, err := s.Write(testData)
		if err != nil {
			t.Error(err)
			return
		}
		if n != len(testData) {
			t.Error(err)
			return
		}

		err = s.Close()
		if err != nil {
			t.Error(err)
		}
	}()

	if !tb.CanDial(list.Multiaddr()) {
		t.Error("CanDial should have returned true")
	}

	connA, err = tb.Dial(ctx, list.Multiaddr(), peerA)
	if err != nil {
		t.Fatal(err)
	}

	s, err := connA.OpenStream(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	n, err := s.Write(testData)
	if err != nil {
		t.Fatal(err)
		return
	}

	if n != len(testData) {
		t.Fatalf("failed to write enough data (a->b)")
		return
	}

	if err = s.CloseWrite(); err != nil {
		t.Fatal(err)
		return
	}

	buf, err := io.ReadAll(s)
	if err != nil {
		t.Fatal(err)
		return
	}
	if !bytes.Equal(testData, buf) {
		t.Errorf("expected %s, got %s", testData, buf)
	}

	if err = s.Close(); err != nil {
		t.Fatal(err)
		return
	}
}
