import $ from "jquery";
import "styles/app.css";
import "./vendor";
import websocketClient from "./services/websocket";

$(document).ready(function () {
  // instantiate the websocket client
  const ws = websocketClient.connect("ws://localhost:8181/echo", (message) => {
    console.log(message);
  });

  $("#send").click(function () {
    const message = "hello !!!!!";
    websocketClient.send(ws, message);
  });
});
