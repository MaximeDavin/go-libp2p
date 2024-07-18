package transport

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/hashicorp/yamux"
	"github.com/libp2p/go-libp2p/p2p/security/noise"

	mplex "github.com/libp2p/go-mplex"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	mss "github.com/multiformats/go-multistream"
)

type TcpTransport struct{}

var _ Transport = &TcpTransport{}

func NewTCPTransport() (*TcpTransport, error) {
	return &TcpTransport{}, nil
}

type MuxedConn struct {
	yamux.Session
	mplex.Multiplex
}

func upgrade(ctx context.Context, conn manet.Conn, addr ma.Multiaddr, pid string, initiator bool) (StreamConn, error) {
	// Upgrade security
	_, err := mss.SelectOneOf([]string{NOISE_ID}, conn)
	if err != nil {
		// Not compatible with noise
		return nil, err
	}
	sconn, err := noise.Secure(conn, initiator, pid)
	if err != nil {
		return nil, err
	}

	// Upgrade stream muxer
	var streamMuxerIDs = []string{YAMUX_ID, MPLEX_ID} // TODO: build this from config
	streamproto, err := mss.SelectOneOf(streamMuxerIDs, sconn)
	if err != nil {
		// No stream muxer compatible
		return nil, err
	}
	switch streamproto {
	case YAMUX_ID:
		// TODO: build and pass yamux config
		yconn, err := yamux.Client(sconn, nil)
		if err != nil {
			return nil, err
		}
		yadapter := YamuxAdapter{Session: yconn}
		return &yadapter, nil
	case MPLEX_ID:
		mconn, err := mplex.NewMultiplex(sconn, initiator, nil)
		if err != nil {
			return nil, err
		}
		madapter := &MplexAdapter{Session: mconn}
		return madapter, nil
	default:
		return nil, errors.New("Programming error: this muxer is not supported")
	}
}

func (t *TcpTransport) Dial(ctx context.Context, addr ma.Multiaddr, pid string) (StreamConn, error) {
	// Get raw connection
	var d manet.Dialer
	conn, err := d.DialContext(ctx, addr)
	if err != nil {
		return nil, err
	}
	return upgrade(ctx, conn, addr, pid, true)
}

func (t *TcpTransport) Listen(addr ma.Multiaddr) (Listener, error) {
	// l, err := manet.Listen(addr)
	// if err != nil {
	// 	return nil, err
	// }
	// var wg sync.WaitGroup
	// defer func() {
	// 	l.Close()
	// 	wg.Wait()
	// }()
	// for {
	// 	conn, err := l.Accept()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	log.Printf("listener %s got connection: %s <---> %s",
	// 		l,
	// 		conn.LocalMultiaddr(),
	// 		conn.RemoteMultiaddr(),
	// 	)
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		// Let's upgrade the connection
	// 		sconn, err := upgrade(context.Background(), conn, conn.RemoteMultiaddr(), "", true)
	// 		if err != nil {
	// 			// Don't bother bubbling this up. We just failed
	// 			// to completely negotiate the connection.
	// 			log.Printf("accept upgrade error: %s (%s <--> %s)",
	// 				err,
	// 				conn.LocalMultiaddr(),
	// 				conn.RemoteMultiaddr())
	// 			return
	// 		}
	// 		select {
	// 		case l.incoming <- sconn:
	// 		}
	// 	}()
	// }
}

type TcpListener struct {
	manet.Listener
	incoming chan StreamConn
}

var _ Listener = &TcpListener{}

func NewTCPListener() (*TcpTransport, error) {
	return &TcpTransport{}, nil
}

func (l *TcpListener) Accept() (StreamConn, error) {
	var wg sync.WaitGroup
	defer func() {
		l.Listener.Close()
		wg.Wait()
		close(l.incoming)
	}()
	for {
		conn, err := l.Listener.Accept()
		if err != nil {
			return nil, err
		}
		log.Printf("listener %s got connection: %s <---> %s",
			l,
			conn.LocalMultiaddr(),
			conn.RemoteMultiaddr(),
		)
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Let's upgrade the connection
			sconn, err := upgrade(context.Background(), conn, conn.RemoteMultiaddr(), "", true)
			if err != nil {
				// Don't bother bubbling this up. We just failed
				// to completely negotiate the connection.
				log.Printf("accept upgrade error: %s (%s <--> %s)",
					err,
					conn.LocalMultiaddr(),
					conn.RemoteMultiaddr())
				return
			}
			select {
			case l.incoming <- sconn:
			}
		}()

	}
}

func (l *TcpListener) Close() error {
	return l.Listener.Close()
}
