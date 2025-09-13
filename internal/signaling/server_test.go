package signaling

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pion/webrtc/v3"
)

func TestHelloWorldExchange(t *testing.T) {
	srv := NewServer()

	req := httptest.NewRequest(http.MethodPost, "/signal", strings.NewReader("hello"))
	req.Header.Set("Authorization", "Bearer secret")
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "hello world" {
		t.Fatalf("expected 'hello world', got %s", w.Body.String())
	}
}

func TestUnauthorized(t *testing.T) {
	srv := NewServer()

	req := httptest.NewRequest(http.MethodPost, "/signal", strings.NewReader("hello"))
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestUnauthorizedWrongToken(t *testing.T) {
	srv := NewServer()

	req := httptest.NewRequest(http.MethodPost, "/signal", strings.NewReader("hello"))
	req.Header.Set("Authorization", "Bearer wrong")
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestWebRTCOfferAnswer(t *testing.T) {
	srv := NewServer()
	ts := httptest.NewServer(srv)
	defer ts.Close()

	peer, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		t.Fatalf("failed to create peer connection: %v", err)
	}
	defer peer.Close()

	if _, err := peer.CreateDataChannel("dc", nil); err != nil {
		t.Fatalf("failed to create data channel: %v", err)
	}

	offer, err := peer.CreateOffer(nil)
	if err != nil {
		t.Fatalf("failed to create offer: %v", err)
	}
	if err := peer.SetLocalDescription(offer); err != nil {
		t.Fatalf("failed to set local description: %v", err)
	}
	<-webrtc.GatheringCompletePromise(peer)

	buf, err := json.Marshal(offer)
	if err != nil {
		t.Fatalf("failed to marshal offer: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/signal", bytes.NewReader(buf))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer secret")
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	var answer webrtc.SessionDescription
	if err := json.NewDecoder(res.Body).Decode(&answer); err != nil {
		t.Fatalf("failed to decode answer: %v", err)
	}
	if answer.Type != webrtc.SDPTypeAnswer {
		t.Fatalf("expected answer type, got %s", answer.Type.String())
	}
	if err := peer.SetRemoteDescription(answer); err != nil {
		t.Fatalf("failed to set remote description: %v", err)
	}
}

func TestSOCKS5Proxy(t *testing.T) {
	// Echo TCP server
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			go func(conn net.Conn) {
				defer conn.Close()
				io.Copy(conn, conn)
			}(c)
		}
	}()

	srv := NewServer()
	ts := httptest.NewServer(srv)
	defer ts.Close()
	defer srv.Close()

	peer, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		t.Fatalf("failed to create peer: %v", err)
	}
	defer peer.Close()

	dc, err := peer.CreateDataChannel("socks5", nil)
	if err != nil {
		t.Fatalf("create data channel: %v", err)
	}

	msgCh := make(chan []byte, 10)
	dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		msgCh <- append([]byte(nil), msg.Data...)
	})
	openCh := make(chan struct{})
	dc.OnOpen(func() { close(openCh) })

	offer, err := peer.CreateOffer(nil)
	if err != nil {
		t.Fatalf("create offer: %v", err)
	}
	if err := peer.SetLocalDescription(offer); err != nil {
		t.Fatalf("set local description: %v", err)
	}
	<-webrtc.GatheringCompletePromise(peer)

	buf, err := json.Marshal(offer)
	if err != nil {
		t.Fatalf("marshal offer: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/signal", bytes.NewReader(buf))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer secret")
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("post offer: %v", err)
	}
	defer res.Body.Close()
	var answer webrtc.SessionDescription
	if err := json.NewDecoder(res.Body).Decode(&answer); err != nil {
		t.Fatalf("decode answer: %v", err)
	}
	if err := peer.SetRemoteDescription(answer); err != nil {
		t.Fatalf("set remote: %v", err)
	}

	<-openCh

	// SOCKS5 handshake
	dc.Send([]byte{0x05, 0x01, 0x00})
	resp := <-msgCh
	if !bytes.Equal(resp, []byte{0x05, 0x00}) {
		t.Fatalf("unexpected method response %v", resp)
	}

	tcpAddr := ln.Addr().(*net.TCPAddr)
	port := uint16(tcpAddr.Port)
	addrMsg := []byte{
		0x05, 0x01, 0x00, 0x01,
		tcpAddr.IP[0], tcpAddr.IP[1], tcpAddr.IP[2], tcpAddr.IP[3],
		byte(port >> 8), byte(port),
	}
	dc.Send(addrMsg)
	resp = <-msgCh
	if len(resp) < 4 || resp[1] != 0x00 {
		t.Fatalf("connect failed: %v", resp)
	}

	dc.Send([]byte("ping"))
	resp = <-msgCh
	if string(resp) != "ping" {
		t.Fatalf("unexpected echo %s", string(resp))
	}
}

func TestHTTPProxy(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))
	defer backend.Close()

	srv := NewServer()
	ts := httptest.NewServer(srv)
	defer ts.Close()
	defer srv.Close()

	peer, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		t.Fatalf("failed to create peer: %v", err)
	}
	defer peer.Close()

	dc, err := peer.CreateDataChannel("http", nil)
	if err != nil {
		t.Fatalf("create dc: %v", err)
	}

	msgCh := make(chan []byte, 10)
	dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		msgCh <- append([]byte(nil), msg.Data...)
	})
	openCh := make(chan struct{})
	dc.OnOpen(func() { close(openCh) })

	offer, err := peer.CreateOffer(nil)
	if err != nil {
		t.Fatalf("create offer: %v", err)
	}
	if err := peer.SetLocalDescription(offer); err != nil {
		t.Fatalf("set local: %v", err)
	}
	<-webrtc.GatheringCompletePromise(peer)

	buf, err := json.Marshal(offer)
	if err != nil {
		t.Fatalf("marshal offer: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/signal", bytes.NewReader(buf))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer secret")
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("post offer: %v", err)
	}
	defer res.Body.Close()
	var answer webrtc.SessionDescription
	if err := json.NewDecoder(res.Body).Decode(&answer); err != nil {
		t.Fatalf("decode answer: %v", err)
	}
	if err := peer.SetRemoteDescription(answer); err != nil {
		t.Fatalf("set remote: %v", err)
	}

	<-openCh

	reqStr := "GET " + backend.URL + " HTTP/1.1\r\nHost: " + backend.Listener.Addr().String() + "\r\n\r\n"
	dc.Send([]byte(reqStr))
	respBytes := <-msgCh
	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(respBytes)), nil)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK || string(body) != "ok" {
		t.Fatalf("unexpected response %d %s", resp.StatusCode, string(body))
	}
}
