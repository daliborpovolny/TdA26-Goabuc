# [Tour de App](https://tourdeapp.cz/)
- [Materiály](https://tourdeapp.cz/vzdelavaci-materialy)
- [Zadání](https://tourdeapp.cz/zadani/introduction)
- [Tour de Cloud](https://tinyurl.com/54ekrruk)

note: this is not the only readme. There is another in apps/server.

## How to run

### Local development
- #### Server:
    - in /apps/server/cmd/tourbackend run: `go run .`
    - runs on port 3000
    - if u change something you have to rerun the command
    - note: it should be now possible to run the server directly from the directory you cloned it to (TdA26-Goabuc by default), using the .bat (Windows) or the .sh (Linux / Mac) file
    - also note that the .bat file now works for both frontend and backend simultaneously and opnens 2 new windows
        - #### Windows
            - in cmd, run `start_app.bat` (this starts the backend, the frontend, and Caddy at the same time => you're done)
        - #### Linux & MacOS
            - in bash, run `sudo chmod +x start_backend.sh` (if you haven't already)
            - run `./start_backend.sh`
            - if this fails, simply run `go run ./apps/server/cmd/tourbackend` in the main directory, or change directories first using `cd` and then run `go run .`
- #### Web
    - for better readability of logs, switch to another bash window
    - run `sudo chmod +x start_frontend.sh`
    - in the main directory, run `./start_frontend.sh`
    - if this fails, run `cd ./apps/web` and then `npm run dev`
    - app is automatically reloaded when you change something

- #### Caddy
    - for better log readability, run this in yet another bash window
    - run `sudo chmod +x start_caddy.sh`
    - in the main directory, run `./start_caddy.sh` &larr; this does not currently work for some reason
    - if this fails, run `caddy run ./apps/caddy`

### Local deployment with docker
This is the apps as they will be run in the cloud
After each change you'll have to rebuild the images and rerun the containers

- #### Server:
    - in /apps/server run: `docker build -t tourbackend`.
    - this will create a docker image
    - to run this docker image ei create a container:
        - `docker run --rm -p 3000:3000 --name tourbackend tourbackend`
- #### Web:
    - in /apps/web run: `docker build -t tourfrontend`.
    - this will create a docker image
    - to run this docker image ei create a container:
        - `docker run --rm -p 3001:3001 --name tourfrontend tourfrontend`

## Deploy to Tour de Cloud
git push automatically deploys to the cloud with a github action
