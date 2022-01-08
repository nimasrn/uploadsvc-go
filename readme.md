# uploadsvc-go

## Description

This service is for uploading images or other files in multiple chunks in Go (Golang). We use Redis to store images/files Infos and disk storage for storing the images data.
With having to be able to upload in multiple chunks the client could have paused/resume the upload process or send smaller payloads in different networks.
this is just for using case study but it aims to be almost production-ready so the Redis and Disk storage is connected to volumes in docker so we could easily scale out the containers.

I tried to use [`go-clean-arch`](https://github.com/bxcodec/go-clean-arch.git) for the project structure over the hexagonal architecture as it was more suitable and based on Rule of Clean Architecture by Uncle Bob.

This project has 4 Domain layer :

- Models Layer
- Repository Layer
- Usecase Layer
- Delivery Layer

### Development

you can run the project while passing the config file as an argument to the run command by default the app will use config.json as config variables. you don't need to pass the file extension as ".json" here.

**run project:** `go run main.go [env]`

example:

```
 go run main.go --env=develop
```

#### Run the Applications

Here is the steps to run it with `docker-compose`

```bash
#move to directory
$ cd workspace

# Clone into YOUR $GOPATH/src
$ git clone https://github.com/nimasrn/uploadsvc-go

#move to project
$ cd uploadsvc-go

# Run the application
$ make run

# check if the containers are running
$ docker ps

# Stop
$ make stop
```

## The API

Our executable expects your HTTP API to implement the following endpoints:

- **Registering an image**:

  - **method**: `POST`
  - **URI**: `/image`
  - **Content-Type**: `application/json`
  - **Request Body**:

    ```json
    {
      "sha256": "abc123easyasdoremi...",
      "size": 123456,
      "chunk_size": 256
    }
    ```

  - **Responses**:
    | Code | Description |
    |----------------------------|------------------------------------|
    | 201 Created | Image successfully registered |
    | 409 Conflict | Image already exists |
    | 400 Bad Request | Malformed request |
    | 415 Unsupported Media Type | Unsupported payload format |

- **Uploading an image chunk**:

  - **method**: `POST`
  - **URI**: `/image/<sha256>/chunks`
  - **Content-Type**: `application/json`
  - **Request Body**:

    ```json
    {
      "id": 1,
      "size": 256,
      "data": "8   888   , 888    Y888 888 888    ,ee 888 888 888 888 ..."
    }
    ```

  - **Responses**:
    | Code | Description |
    |---------------|------------------------------------|
    | 201 Created | Chunk successfully uploaded |
    | 409 Conflict | Chunk already exists |
    | 404 Not Found | Image not found |

- **Downloading an image**:

  - **method**: `GET`
  - **URI**: `/image/<sha256>`
  - **Accept**: `text/plain`
  - **Responses**:
    | Code | Description |
    |---------------|------------------------------------|
    | 200 OK | Image successfully downloaded |
    | 404 Not Found | Image not found |

  - **Note**: This endpoint returns plain text, not JSON. It should return the whole image instead of separate chunks.

- **Errors**:

  - **Accept**: `application/json`
  - **Response body**:

    ```json
    {
      "code": "400",
      "message": "Chunk ID field is missing."
    }
    ```

### Tools Used:

In this project, I use some tools listed below. If I want to mention some I can say Gin as a core framework and Redis.

- All libraries listed in [`go.mod`](https://github.com/nimasrn/uploadsvc-go/blob/master/go.mod)
