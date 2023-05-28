import React, { useRef } from 'react';
import { w3cwebsocket as W3CWebSocket } from 'websocket';

const App = () => {
  const wsClient = useRef(null);

  const startCall = async () => {
    wsClient.current = new W3CWebSocket('ws://localhost:8080/ws');

    const peerConnection = new RTCPeerConnection();

    wsClient.current.onmessage = async (message) => {
      const sdp = JSON.parse(message.data);
      await peerConnection.setRemoteDescription(new RTCSessionDescription(sdp));
      const answer = await peerConnection.createAnswer();
      await peerConnection.setLocalDescription(answer);
      wsClient.current.send(JSON.stringify(answer));
    };

    const offer = await peerConnection.createOffer();
    await peerConnection.setLocalDescription(offer);
    wsClient.current.onopen = () => {
      wsClient.current.send(JSON.stringify(offer));
    };
  };

  return (
    <div>
      <button onClick={startCall}>Call Support</button>
    </div>
  );
};

export default App;