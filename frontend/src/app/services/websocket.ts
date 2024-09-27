// websockets client
class wsClient {
  socket: WebSocket | null = null;
  connect(url: string, onMessage: (evt: any) => void): Promise<wsClient> {
    this.socket = new WebSocket(url);
    this.socket.onmessage = function (evt: any) {
      onMessage(evt.data);
    };
    this.socket.onopen = function (evt) {
      console.log("OPEN");
    };
    this.socket.onclose = function (evt) {
      console.log("CLOSE");
    };
    this.socket.onerror = function (evt: any) {
      console.log("ERROR ", evt);
    };
    return new Promise((resolve, reject) => {
      this.socket?.addEventListener("open", () => {
        resolve(this);
      });
    });
  }
  send(message: string) {
    this.socket.send(message);
  }
  close() {
    this.socket.close();
  }
}

export default wsClient;
