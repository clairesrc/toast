// websockets client

const websocketClient = {
  connect: function (url: string, onMessage: (message: string) => void) {
    var ws = new WebSocket(url);
    ws.onmessage = function (evt: any) {
      onMessage(evt);
    };
    ws.onopen = function (evt) {
      console.log("OPEN");
    };
    ws.onclose = function (evt) {
      console.log("CLOSE");
      ws = null;
    };
    ws.onerror = function (evt: any) {
      console.log("ERROR ", evt);
    };
    return ws;
  },
  send: function (ws: WebSocket, message: string) {
    ws.send(message);
  },
  close: function (ws: WebSocket) {
    ws.close();
  },
};

export default websocketClient;
