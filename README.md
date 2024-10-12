![image](https://github.com/user-attachments/assets/52576401-fb8e-498a-bb31-8ec0282b56ff)

# 🍞toast engine

Extremely cursed PVP multiplayer RPG engine written in Go, Typescript with Websockets. Mostly just a toy for me to use as a sandbox and figure out how common game mechanics are implemented.

Controls - movement (arrow keys), attack (spacebar), dodgeroll (ctrl), health bar, stamina. A rudimentary player collision system prevents from moving into other players.

This is a multiplayer combat game. Run the server through Docker, then players can join the session from their web browser.

## Usage
Start locally:

```bash
docker compose up --build --force-recreate -d
```

Then open: http://localhost:4000 in browser

## Attributions
<a href="https://www.freepik.com/free-photo/sand-ground-textured_1198415.htm#query=dirt%20texture%20seamless&position=0&from_view=keyword&track=ais_hybrid&uuid=7c26f5db-7716-482c-9bd3-333e77ab092a">Sand dirt ground texture by luis_molinero</a> on Freepik
