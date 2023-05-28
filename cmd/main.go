package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(":8080", nil)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Create a new WebRTC PeerConnection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		panic(err)
	}

	// Set the handler for ICE connection state
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})

	// Wait for the WebSocket to close
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Handle incoming SDP
		err = handleIncomingSDP(peerConnection, string(message))
		if err != nil {
			fmt.Println("Error handling SDP:", err)
			break
		}
	}
}

func handleIncomingSDP(peerConnection *webrtc.PeerConnection, sdpStr string) error {
	// Unmarshal the incoming SDP
	sdp := webrtc.SessionDescription{}
	err := sdp.Unmarshal([]byte(sdpStr))
	if err != nil {
		return err
	}

	// Set the remote SDP
	err = peerConnection.SetRemoteDescription(sdp)
	if err != nil {
		return err
	}

	// Create an answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		return err
	}

	// Set the local SDP
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		return err
	}

	return nil
}