import { PlayerState } from "./services/state";
import wsClient from "./services/websocket";

export default class Player {
  websocketClient: wsClient;
  player: PlayerState;
  constructor(player: PlayerState, websocketClient: wsClient) {
    this.player = player;
    this.websocketClient = websocketClient;
  }

  send(eventType: string, data) {
    this.websocketClient.send(
      JSON.stringify({
        data: this.player,
        type: eventType,
      })
    );
  }
  moveUp() {
    this.player.facing = "up";
    this.send("walk", this.player);
  }
  moveDown() {
    this.player.facing = "down";
    this.send("walk", this.player);
  }
  moveLeft() {
    this.player.facing = "left";
    this.send("walk", this.player);
  }
  moveRight() {
    this.player.facing = "right";
    this.send("walk", this.player);
  }
  attack() {
    this.player.isAttacking = true;
    this.send("attack", this.player);
  }
  dodge() {
    this.send("dodge", this.player);
  }
}
