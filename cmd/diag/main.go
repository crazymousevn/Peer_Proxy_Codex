package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pion/webrtc/v3"
)

func main() {
	signalURL := flag.String("signal", "http://localhost:8080/signal", "signaling URL")
	forceRelay := flag.Bool("force-relay", false, "force TURN relay mode")
	flag.Parse()

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
			{URLs: []string{"turn:turn.example.com:3478"}, Username: "user", Credential: "pass"},
		},
	}
	if *forceRelay {
		config.ICETransportPolicy = webrtc.ICETransportPolicyRelay
	}

	pc, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	pc.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Println("connection state:", s)
	})

	dc, err := pc.CreateDataChannel("diag", nil)
	if err != nil {
		panic(err)
	}
	dc.OnOpen(func() { fmt.Println("data channel open") })

	offer, err := pc.CreateOffer(nil)
	if err != nil {
		panic(err)
	}
	gatherComplete := webrtc.GatheringCompletePromise(pc)
	if err := pc.SetLocalDescription(offer); err != nil {
		panic(err)
	}
	<-gatherComplete

	reqBody, err := json.Marshal(struct {
		webrtc.SessionDescription
		ForceRelay bool `json:"forceRelay"`
	}{SessionDescription: *pc.LocalDescription(), ForceRelay: *forceRelay})
	if err != nil {
		panic(err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, *signalURL, bytes.NewReader(reqBody))
	if err != nil {
		panic(err)
	}
	httpReq.Header.Set("Authorization", "Bearer secret")
	httpReq.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("bad status: %s", res.Status))
	}

	var answer webrtc.SessionDescription
	if err := json.NewDecoder(res.Body).Decode(&answer); err != nil {
		panic(err)
	}
	if err := pc.SetRemoteDescription(answer); err != nil {
		panic(err)
	}

	fmt.Println("Waiting for connection, press Ctrl+C to exit")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
