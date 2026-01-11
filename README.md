# [Tour de App](https://tourdeapp.cz/)
- [Materiály](https://tourdeapp.cz/vzdelavaci-materialy)
- [Zadání](https://tourdeapp.cz/zadani/introduction)
- [Tour de Cloud](https://tinyurl.com/54ekrruk)

note: this is not the only readme. There is another in apps/server.

## How to run

### Local development

- #### Server (port :3000)
    - it is now possible to run the server directly from the directory you cloned it to (TdA26-Goabuc by default), using the .bat (Windows) or the .sh (Linux / Mac) file
    - also note that the .bat file now works for both frontend, backend, and caddy simultaneously and opnens 3 new windows
      
        - #### Windows
            - in cmd, run `start_app.bat` (this starts the backend, the frontend, and Caddy at the same time &rarr; you're done)
              
        - #### Linux & MacOS
            - in bash, run `sudo chmod +x start_backend.sh` (if you haven't yet)
            - run `./start_backend.sh`
              
    - if this fails, simply run `go run ./apps/server/cmd/tourbackend` in the main directory
    - if you change something you have to rerun the command
              
-  #### Web (port :3001)
    - for better readability of logs, switch to another bash window
    - run `sudo chmod +x start_frontend.sh` (if you haven't yet)
    - in the main directory, run `./start_frontend.sh`
    - if this fails, run `cd ./apps/web` and then `npm run dev`
    - app is automatically reloaded when you change something
    - after updating the .svelte files, before commiting, run `npm run format` in apps/web

- #### Caddy (port :80)
    - for better log readability, run this in yet another bash window
    - run `sudo chmod +x start_caddy.sh` (if you haven't yet)
    - the 'start_caddy.sh' file should automatically take care of any processes running on its port and then run itself. this will require admin privileges
    - in the main directory, run `./start_caddy.sh`
    - if this fails, try running it with admin privileges: `sudo ./start_caddy.sh`
    - if this also fails, run `caddy run ./apps/caddy`
    - if this also somehow fails, run `sudo caddy run ./apps/caddy` &larr; this should always work if everything's in order. if this command is also failing, read the logs and god help you

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
