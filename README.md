# [Tour de App](https://tourdeapp.cz/)
- [Materiály](https://tourdeapp.cz/vzdelavaci-materialy)
- [Zadání](https://tourdeapp.cz/zadani/introduction)
- [Tour de Cloud](https://tinyurl.com/54ekrruk)

## How to run

### Local development
- #### Server:
    - in /apps/server/cmd/tourbackend run: `go run .`
    - runs on port 3000
    - if u change something you have to rerun the command
    - note: it should be now possible to run the server directly from the directory you cloned it to (TdA26-Goabuc by default), using .bat (windows) or the .sh (Linux / Mac) file
        - ##### Windows
            - in the terminal, run `start_app.bat`
        - ##### Linux & MacOS
            - in bash, run `sudo chmod +x start_backend.sh` (if you haven't already)
            - run `./start_backend.sh`
            - please note that I am unable to test on MacOS
- #### Web
    - in /apps/web run: `npm run dev`
    - app is automatically reloaded when you change something

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
git push automatically deploys to the cloud with a github action - once the fee is paid
as the fee hasn't been payed yet this is untested
