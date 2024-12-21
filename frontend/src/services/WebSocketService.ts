let socket: WebSocket | null = null;
let onMessageCallback: Function | null = null;

export function connectWebSocket(url: string, messageCallback: Function): void {
  socket = new WebSocket(url);
  onMessageCallback = messageCallback;

  socket.onmessage = (event) => {
    const message = JSON.parse(event.data);
    if (onMessageCallback) {
      onMessageCallback(message); // Appelle le callback pour mettre à jour l'interface
    }
  };

  socket.onopen = () => console.log("WebSocket connected");
  socket.onclose = () => console.log("WebSocket disconnected");
  socket.onerror = (error) => console.error("WebSocket error:", error);
}

export function sendMessage(message: any): void {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify(message));
  } else {
    console.error("WebSocket is not connected");
  }
}

export function disconnectWebSocket(): void {
  if (socket) {
    socket.close();
    socket = null;
  }
}
