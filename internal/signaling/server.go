package signaling

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/pion/webrtc/v3"
)

// Server represents a basic signaling server.
type Server struct {
	mu    sync.Mutex
	conns []*webrtc.PeerConnection
}

// NewServer creates a new Server.
func NewServer() *Server {
	return &Server{}
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || r.URL.Path != "/signal" {
		http.NotFound(w, r)
		return
	}

	if r.Header.Get("Authorization") != "Bearer secret" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if string(body) == "hello" {
		_, _ = w.Write([]byte("hello world"))
		return
	}

	var offer webrtc.SessionDescription
	if err := json.Unmarshal(body, &offer); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Handle incoming data channels for proxying.
	peerConnection.OnDataChannel(func(dc *webrtc.DataChannel) {
		switch dc.Label() {
		case "socks5":
			handleSOCKS5(dc)
		case "http":
			handleHTTP(dc)
		}
	})

	s.mu.Lock()
	s.conns = append(s.conns, peerConnection)
	s.mu.Unlock()

	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
	if err := peerConnection.SetLocalDescription(answer); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	<-gatherComplete

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(peerConnection.LocalDescription()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Close shuts down all active peer connections.
func (s *Server) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, pc := range s.conns {
		_ = pc.Close()
	}
	s.conns = nil
}

// handleSOCKS5 proxies a SOCKS5 connection over the given data channel.
func handleSOCKS5(dc *webrtc.DataChannel) {
	type state struct {
		stage int
		conn  net.Conn
	}
	st := &state{}

	dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		data := msg.Data
		switch st.stage {
		case 0:
			// methods negotiation
			if len(data) < 3 || data[0] != 0x05 {
				_ = dc.Close()
				return
			}
			_ = dc.Send([]byte{0x05, 0x00})
			st.stage = 1
		case 1:
			// connect request (IPv4 only)
			if len(data) < 10 || data[0] != 0x05 || data[1] != 0x01 || data[3] != 0x01 {
				_ = dc.Close()
				return
			}
			addr := net.IP(data[4:8]).String()
			port := binary.BigEndian.Uint16(data[8:10])
			conn, err := net.Dial("tcp", net.JoinHostPort(addr, strconv.Itoa(int(port))))
			if err != nil {
				_ = dc.Send([]byte{0x05, 0x01, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
				_ = dc.Close()
				return
			}
			st.conn = conn
			_ = dc.Send([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
			st.stage = 2
			go func() {
				io.Copy(&dataChannelWriter{dc}, conn)
				_ = dc.Close()
			}()
		case 2:
			if st.conn != nil {
				st.conn.Write(data)
			}
		}
	})
}

// handleHTTP proxies HTTP/HTTPS over the given data channel.
func handleHTTP(dc *webrtc.DataChannel) {
	var (
		stage int
		conn  net.Conn
		buf   bytes.Buffer
	)

	dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		if stage == 0 {
			buf.Write(msg.Data)
			if !bytes.Contains(buf.Bytes(), []byte("\r\n\r\n")) {
				return
			}
			req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(buf.Bytes())))
			if err != nil {
				_ = dc.Close()
				return
			}
			if req.Method == http.MethodConnect {
				host := req.Host
				conn, err = net.Dial("tcp", host)
				if err != nil {
					dc.Send([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
					_ = dc.Close()
					return
				}
				dc.Send([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
				stage = 1
				go func() {
					io.Copy(&dataChannelWriter{dc}, conn)
					_ = dc.Close()
				}()
			} else {
				if req.URL.Scheme == "" {
					req.URL.Scheme = "http"
				}
				if req.URL.Host == "" {
					req.URL.Host = req.Host
				}
				req.RequestURI = ""
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					dc.Send([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
					_ = dc.Close()
					return
				}
				var b bytes.Buffer
				resp.Write(&b)
				dc.Send(b.Bytes())
				_ = dc.Close()
			}
			buf.Reset()
		} else {
			if conn != nil {
				conn.Write(msg.Data)
			}
		}
	})
}

// dataChannelWriter adapts a DataChannel to an io.Writer.
type dataChannelWriter struct {
	dc *webrtc.DataChannel
}

func (w *dataChannelWriter) Write(p []byte) (int, error) {
	if err := w.dc.Send(p); err != nil {
		return 0, err
	}
	return len(p), nil
}
