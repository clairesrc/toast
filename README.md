# 🍞toast engine

![image](https://github.com/user-attachments/assets/cea0ccda-5d4c-43dc-8c71-c357e72d5bb2)

Extremely cursed multiplayer RPG engine written in Go, Typescript with Websockets. Not quite what I would call a "game" or an "engine," mostly just a toy for me to use as a sandbox and figure out how common game mechanics are implemented.

Combat mechanics include movement (arrow keys), attack (spacebar), dodgeroll (ctrl), health bar, stamina. A rudimentary player collision system prevents from moving into other players.

The cursed part mainly refers to how I draw things in the browser: instead of something sensible like canvas or webgl, everything you see is drawn using direct DOM manipulation and CSS -- every object onscreen is an HTML element, animations are done using classes and keyframes, etc.

Every action, from movement to attack animation duration to collision detection, is validated server-side -- the client just reacts to updates in the state object.

You may notice the Docker compose file specifies a some persistence layer stuff i.e. a MongoDB and Postgres images -- these aren't in use yet but will eventually allow for multiple server instances, player accounts, and real RPG mechanics like experience points and equipment.

Some other interesting things planned for the serverside include Kafka for game event delegation, which combined with some strategically placed load balancers could enable some pretty ridiculous scaling since I would be able to decompose the server monolith into a series of individual event handlers. At that point all the server entrypoint does is handle incoming Websocket messages by emitting them to Kafka, and also sending out heartbeat game state updates, which at that point will be fetched from MongoDB.

of course, to really leverage the server scaling gains i will need to optimize the clientside as well, mainly by ensuring the game state it receives only contains data relevant to what's onscreen or at least within the player's current world instance.

basically, i have a lot of neat ideas that will keep me hacking on this project for a long time. not sure it will ever be of use to anyone, but i am having fun regardless!


## Usage
Start locally:

```bash
docker compose up --build --force-recreate -d
```

Then open: http://localhost:4000 in browser

## Attributions
<a href="https://www.freepik.com/free-photo/sand-ground-textured_1198415.htm#query=dirt%20texture%20seamless&position=0&from_view=keyword&track=ais_hybrid&uuid=7c26f5db-7716-482c-9bd3-333e77ab092a">Sand dirt ground texture by luis_molinero</a> on Freepik
