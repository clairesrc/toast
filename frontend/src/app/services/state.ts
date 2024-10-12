// game state
type PlayerState = {
  x: number;
  y: number;
  zIndex: number;
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

const playerHitboxWidth: number = 24;
const playerHitboxHeight: number = 12;
const playerSpriteWidth: number = 48;
const playerSpriteHeight: number = 48;

type playerBoundingBox = {
  hitbox: {
    x: number;
    y: number;
    width: number;
    height: number;
  };
  sprite: {
    x: number;
    y: number;
    width: number;
    height: number;
  };
};

var initialState: GameState = {
  players: [],
};

const zIndexPrefix: number = 100;

/* build player bounding box based on player position */
const getPlayerBoundingBox = (player: PlayerState): playerBoundingBox => {
  return {
    hitbox: {
      x: player.x + (playerSpriteWidth - playerHitboxWidth) / 2,
      y: player.y + (playerSpriteHeight - playerHitboxHeight) / 2,
      width: playerHitboxWidth,
      height: playerHitboxHeight,
    },
    sprite: {
      x: player.x,
      y: player.y,
      width: playerSpriteWidth,
      height: playerSpriteHeight,
    },
  };
};

/* draw to the screen based on gamestate */
const renderFromState = (
  oldState: GameState,
  state: GameState,
  gameWorld: HTMLElement,
  userPlayerName: string
) => {
  // render players as needed
  state.players.forEach((player: PlayerState) => {
    // compare with oldState to figure out if player needs to be re-rendered
    const oldPlayer: false | PlayerState = findPlayer(oldState, player.name);
    if (
      oldPlayer &&
      oldPlayer.x === player.x &&
      oldPlayer.y === player.y &&
      oldPlayer.facing === player.facing &&
      oldPlayer.isAttacking === player.isAttacking &&
      oldPlayer.isWalking === player.isWalking &&
      oldPlayer.isDodging === player.isDodging &&
      oldPlayer.health === player.health &&
      player.name == userPlayerName &&
      oldPlayer.stamina === player.stamina
    ) {
      return;
    }

    player.zIndex = parseInt(`${zIndexPrefix}${player.y}`);
    const playerBoundingBox = getPlayerBoundingBox(player);

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
        playerBoundingBox.sprite.y
      }px; left: ${playerBoundingBox.sprite.x}px; z-index: ${zIndexPrefix}${
        player.zIndex
      };">
        <div class="player-sprite"></div><div class="player-data">
          <div class="player-name">${player.name}</div>
          <div class="player-health-bar">
            <div class="player-health-figure">${player.health}</div>
            <div class="player-health-bar-inner" style="width: ${
              player.health
            }%"></div>
          </div>
          ${
            player.name == userPlayerName
              ? `
          <div class="player-stamina-bar">
            <div class="player-stamina-bar-inner" style="width: ${player.stamina}%"></div>
          </div>`
              : ""
          }
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
      currentPlayer.style.zIndex = `${player.zIndex}`;
      const playerHealthBar = currentPlayer.querySelector(
        ".player-health-bar-inner"
      ) as HTMLElement;
      if (playerHealthBar) {
        playerHealthBar.style.width = `${player.health}%`;
      }

      if (player.name == userPlayerName) {
        // update player stamina
        const playerStaminaBar = currentPlayer.querySelector(
          ".player-stamina-bar-inner"
        ) as HTMLElement;
        if (playerStaminaBar) {
          playerStaminaBar.style.width = `${player.stamina}%`;
        }
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
