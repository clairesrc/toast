import $ from "jquery";
import "styles/app.css";
import "./vendor";
import wsClient from "./services/websocket";
import {
  initialState,
  renderFromState,
  findPlayer,
  updatePlayer,
} from "./services/state";

$(document).ready(function () {
  var state = JSON.parse(JSON.stringify(initialState));
  var oldState = JSON.parse(JSON.stringify(state));

  // random player name
  const playerName = `player${Math.floor(Math.random() * 1000)}`;

  console.log(playerName);

  const player1 = {
    x: 0,
    y: 0,
    name: playerName,
    isAttacking: false,
    health: 100,
    facing: "right",
    skin: "default",
  };

  // instantiate the websocket client
  const webSocketClient = new wsClient();
  webSocketClient
    .connect("ws://localhost:8181/state", (gameStateFromServer) => {
      // hold on to last version of state
      oldState = JSON.parse(JSON.stringify(state));

      // update state from websocket
      state = JSON.parse(gameStateFromServer);
    })
    .then(initGame);

  function initGame(webSocketClient: wsClient) {
    // re-render the game world
    const gameWorld = document.getElementById("game-world");
    const render = () => {
      renderFromState(oldState, state, gameWorld);
      oldState = JSON.parse(JSON.stringify(state));
      requestAnimationFrame(render);
    };

    render();

    // send refresh event every 100ms
    setInterval(() => {
      webSocketClient.send(
        JSON.stringify({
          data: {},
          type: "refresh",
        })
      );
    }, 100);

    // send a message to the server to add the player
    webSocketClient.send(
      JSON.stringify({
        data: player1,
        type: "join",
      })
    );

    // react to keypress events using jquery
    $(document).keydown(function (e) {
      const player = findPlayer(state, playerName);
      if (!player) {
        return;
      }
      if (e.key === "ArrowUp") {
        player.facing = "up";
        webSocketClient.send(
          JSON.stringify({
            data: player,
            type: "walk",
          })
        );
      } else if (e.key === "ArrowDown") {
        player.facing = "down";
        webSocketClient.send(
          JSON.stringify({
            data: player,
            type: "walk",
          })
        );
      } else if (e.key === "ArrowLeft") {
        player.facing = "left";
        webSocketClient.send(
          JSON.stringify({
            data: player,
            type: "walk",
          })
        );
      } else if (e.key === "ArrowRight") {
        player.facing = "right";
        webSocketClient.send(
          JSON.stringify({
            data: player,
            type: "walk",
          })
        );
      } else if (e.key === " ") {
        player.isAttacking = true;
        webSocketClient.send(
          JSON.stringify({
            data: player,
            type: "attack",
          })
        );
      }
      state = updatePlayer(state, player);
    });
  }

  // close the websocket connection when the page is closed
  window.addEventListener("beforeunload", function () {
    webSocketClient.send(
      JSON.stringify({
        data: player1,
        type: "leave",
      })
    );
    webSocketClient.close();
  });
});
