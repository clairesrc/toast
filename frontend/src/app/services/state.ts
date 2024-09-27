// game state
type PlayerState = {
  x: number;
  y: number;
  facing: string;
  name: string;
  isWalking: boolean;
  isAttacking: boolean;
  isDodging: boolean;
  health: number;
  stamina: number;
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
      oldPlayer.isWalking === player.isWalking &&
      oldPlayer.isDodging === player.isDodging &&
      oldPlayer.stamina === player.stamina &&
      oldPlayer.health === player.health
    ) {
      return;
    }

    // if player has been removed from state
    if (!player) {
      // remove current player
      const currentPlayer = gameWorld.querySelector(
        `[data-playername="${player.name}"]`
      );
      if (currentPlayer) {
        currentPlayer.remove();
      }
      return;
    }

    // add new player element if just added
    if (!oldPlayer) {
      gameWorld.innerHTML += `<div data-playername="${
        player.name
      }" class="player ${player.isAttacking ? "attacking" : ""} ${
        player.isDodging ? "dodging" : ""
      } ${player.isWalking ? "walking" : ""} ${
        player.health == 0 ? "dead" : ""
      } facing-${player.facing} skin-${player.skin}" style="top: ${
        player.y
      }px; left: ${player.x}px;">
        <div class="player-sprite"></div><div class="player-data">
          <div class="player-name">${player.name}</div>
          <div class="player-health-bar">
            <div class="player-health-figure">${player.health}</div>
            <div class="player-health-bar-inner" style="width: ${
              player.health
            }%"></div>
          </div>
        </div>
      </div>`;
      return;
    }

    // update player classes and position
    const currentPlayer = gameWorld.querySelector(
      `[data-playername="${player.name}"]`
    ) as HTMLElement;
    if (currentPlayer) {
      currentPlayer.className = `player ${
        player.isAttacking ? "attacking" : ""
      } ${player.isDodging ? "dodging" : ""} ${
        player.isWalking ? "walking" : ""
      } ${player.health == 0 ? "dead" : ""} facing-${player.facing} skin-${
        player.skin
      }`;
      const playerHealthFigure = currentPlayer.querySelector(
        ".player-health-figure"
      );
      if (playerHealthFigure) {
        playerHealthFigure.innerHTML = `${player.health}`;
      }

      // update player position using inline style tag
      currentPlayer.style.top = `${player.y}px`;
      currentPlayer.style.left = `${player.x}px`;
      const playerHealthBar = currentPlayer.querySelector(
        ".player-health-bar-inner"
      ) as HTMLElement;
      if (playerHealthBar) {
        playerHealthBar.style.width = `${player.health}%`;
      }
    }
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

export { initialState, renderFromState, findPlayer, updatePlayer, PlayerState };
