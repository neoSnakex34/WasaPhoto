# Source code for neoSnakex34's WASAPhoto  

## What am I looking at? 

This is a project made during my WASA (web and software architecture) course. 
It aims to provide an experience similar to those you can have in photo sharing social networks. 
WASAPhoto (graphically *WasaPHOTO*) is a single page WebApp implementing a **REST** API. 

* **OPENAPI** specs can be found at `doc/api.yaml`.
* API implementation and Core (BACKEND) features are implemented using **Golang** and can be found in `service/`.
* Users will use API functions through a UI (FRONTEND) made with **Vue.js**, aesthetically appealing and easy to use. Source found in `webui`. 
* Deployment is made using **Docker** containers. See the `Dockerfile.*` in root.


## How to build

If you're not using the WebUI, or if you don't want to embed the WebUI into the final executable, then:

```shell
go build ./cmd/webapi/
```

If you're using the WebUI and you want to embed it into the final executable:

```shell
./open-npm.sh
# (here you're inside the NPM container)
npm run build-embed
exit
# (outside the NPM container)
go build -tags webui ./cmd/webapi/
```

## How to run (in development mode)

You can launch the backend only using:

```shell
go run ./cmd/webapi/
```

If you want to launch the WebUI, open a new tab and launch:

```shell
./open-npm.sh
# (here you're inside the NPM container)
npm run dev
```

## How to run using docker images 

* Firstly you should build the images: 

```shell
docker build -t wasaphoto-backend:latest -f Dockerfile.backend .
```

```shell
docker build -t wasaphoto-frontend:latest -f Dockerfile.frontend .
```

* Then you can run them _(changing flags if you need, rm will delete everything after you exit the container [suitable for a demo])_
```shell
docker run -it --rm -p 3000:3000 wasaphoto-backend:latest
```

```shell
docker run -it --rm -p 8081:80 wasaphoto-frontend:latest
```

