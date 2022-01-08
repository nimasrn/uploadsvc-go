# uploadsvc-go

## Description

### Development

you can run the project while passing the config file as an argument to the run command by default the app will use config.json as config variables.

**run project:** `go run main.go [env]`

example:

```
 go run main.go --env=env.develop
```

#### Run the Applications

Here is the steps to run it with `docker-compose`

```bash
#move to directory
$ cd workspace

# Clone into YOUR $GOPATH/src
$ git clone https://github.com/nimasrn/ecgsvc-go

#move to project
$ cd ecgsvc-go

# Run the application
$ make run

# check if the containers are running
$ docker ps

# Stop
$ make stop
```

## The API

### Tools Used:

In this project, I use some tools listed below. If I want to mention some I can say Gin as a core framework and Postgres.

- All libraries listed in [`go.mod`](https://github.com/nimasrn/ecgsvc-go/blob/main/go.mod)
