let socket = null;

export function connectWebSocket(url, onMessageCallback) {
  socket = new WebSocket(url);

  socket.onmessage = (event) => {
    const message = JSON.parse(event.data);
    if (onMessageCallback) {
      onMessageCallback(message);
    }
  };

  socket.onopen = () => console.log("WebSocket connected");
  socket.onclose = () => console.log("WebSocket disconnected");
  socket.onerror = (error) => console.error("WebSocket error:", error);
}

export function sendMessage(message) {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify(message));
  } else {
    console.error("WebSocket is not connected");
  }
}

export function disconnectWebSocket() {
  if (socket) {
    socket.close();
    socket = null;
  }
}
