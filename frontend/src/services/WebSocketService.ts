let socket: WebSocket | null = null;
let onMessageCallback: Function | null = null;

export function connectWebSocket(url: string, messageCallback: Function): void {
  const token = "Mjo5V8xQW3xes6I2DnItWObxlpnh4uQIU0xWGiDNCb92nFJoo7q6NDya2x/ZRT1U1";
  const safeToken = encodeURIComponent(token);

  socket = new WebSocket(url+'?token='+safeToken);

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
