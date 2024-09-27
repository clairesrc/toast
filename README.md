# breadgine

![screenshot](https://github.com/user-attachments/assets/4748bbcb-c7d3-4f0d-b616-c57a02cb8bba)

Extremely cursed multiplayer RPG engine written in Go, Typescript with Websockets. Not quite what I would call a "game" or an "engine," mostly just a toy for me to use as a sandbox and figure out how common game mechanics are implemented.

The cursed part mainly refers to how I draw things in the browser: instead of something sensible like canvas or webgl, everything you see is drawn using direct DOM manipulation and CSS -- every object onscreen is an HTML element, animations are done using classes and keyframes, etc.

Every action, from movement to attack animation duration to collision detection, is validated server-side -- the client just reacts to updates in the state object.

You may notice the Docker compose file specifies a some persistence layer stuff i.e. a MongoDB and Postgres images -- these aren't in use yet but will eventually allow for multiple server instances, player accounts, and real RPG mechanics like experience points and equipment.


## Usage
Start locally:

```bash
docker compose up --build --force-recreate -d
```

Then open: http://localhost:4000 in browser

## Attributions
<a href="https://www.freepik.com/free-photo/sand-ground-textured_1198415.htm#query=dirt%20texture%20seamless&position=0&from_view=keyword&track=ais_hybrid&uuid=7c26f5db-7716-482c-9bd3-333e77ab092a">Sand dirt ground texture by luis_molinero</a> on Freepik
