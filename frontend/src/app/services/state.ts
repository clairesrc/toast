// game state
type PlayerState = {
  x: number;
  y: number;
  facing: string;
  name: string;
  isAttacking: boolean;
  health: number;
  skin: string;
};
type GameState = {
  players: PlayerState[];
};

var initialState: GameState = {
  players: [],
};

const renderFromState = (
  oldState: GameState,
  state: GameState,
  gameWorld: HTMLElement
) => {
  // render players as needed
  state.players.forEach((player: PlayerState) => {
    // compare with oldState to figure out if player needs to be re-rendered
    const oldPlayer = findPlayer(oldState, player.name);
    if (
      oldPlayer &&
      oldPlayer.x === player.x &&
      oldPlayer.y === player.y &&
      oldPlayer.facing === player.facing &&
      oldPlayer.isAttacking === player.isAttacking &&
      oldPlayer.health === player.health
    ) {
      return;
    }

    // remove current player
    const currentPlayer = gameWorld.querySelector(
      `[data-playername="${player.name}"]`
    );
    if (currentPlayer) {
      currentPlayer.remove();
    }

    gameWorld.innerHTML += `<div data-playername="${
      player.name
    }" class="player ${player.isAttacking ? "attacking" : ""} ${
      player.health == 0 ? "dead" : ""
    } facing-${player.facing} skin-${player.skin}" style="top: ${
      player.y
    }px; left: ${player.x}px;">
      <div class="player-name">${player.name}</div>
      <div class="player-health">${player.health}</div>
    </div>`;
  });

  return state;
};

const updatePlayer = (state: GameState, player: PlayerState): GameState => {
  const index = state.players.findIndex((p) => p.name === player.name);
  if (index === -1) {
    return state;
  }
  const newState = JSON.parse(JSON.stringify(state));
  newState.players[index] = player;
  return newState;
};

const findPlayer = (state: GameState, name: string) => {
  return state.players.find((player) => player.name === name) || false;
};

export { initialState, renderFromState, findPlayer, updatePlayer };
