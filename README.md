# Mock Service

![technology Go](https://img.shields.io/badge/technology-go-blue.svg)

A simple mock server in Go.

## Installation

#### With Dockers

Running this project with Dockers is the best and easiest option.

Once we have checked out our project, and we are in the root folder we need to build our image:

```sh
docker build -t mock-service:latest . 
```

By default, the application will look for mocks in the file set in `MOCKS_FILE` environment variable inside the container. So we can simply run this project by running our image with the following command:

```sh
docker run -v /tmp:/tmp -e MOCKS_FILE=/tmp/mocks.json -p 8080:8080 --name mock-service mock-service
```

Alternatively, Mock Service can be run with MySQL database:
```sh
docker run -e MOCKS_DATASOURCE=mysql -e DB_USER={{user}} -e DB_PASSWORD={{password}} - DB_HOST={{host}} -e DB_PORT={{port}} -p 8080:8080 --name mock-service mock-service
```
Database must contain a schema `mockserver` with required tables. Follow[`this link`](https://github.com/nicopozo/mockserver/blob/master/scripts/init.sql "Init sql script") to get the creation script.

#### By compiling with Go
Mock Service can be compiled and run without the need of a Dockers installation. In order to compile this application, we need Go 1.18 installed.

```sh
cd cmd/mocks 
go build .
```

Before running we configure the JSON file with the mocks.

```sh
export MOCKS_FILE=/tmp/mocks.json
```

Now, simply run the application:

```sh
./mocks
```

In order to use the application with MySQL database instead of a file, set these environment variables before running the app.

```sh
export MOCKS_DATASOURCE=mysql
export DB_USER={{user}}
export DB_PASSWORD={{password}} 
export DB_HOST={{host}}
export DB_PORT={{port}}
```

and then run the app witj `./mocks` command (run the [`init database script`](https://github.com/nicopozo/mockserver/blob/master/scripts/init.sql "Init sql script") before running the app).

## How to use it

#### Create a new mock

```sh
curl --location --request POST 'localhost:8080/mock-service/rules' \
--header 'Content-Type: application/json' \
--data-raw '{
    "group": "Users",
    "name": "Get User",
    "path": "/users/{id}",
    "strategy": "normal",
    "method": "GET",
    "status": "enabled",
    "responses": [
        {
            "body": "{\"user_id\":\"{the_id}\"}",
            "content_type": "application/json",
            "http_status": 200,
            "delay": 0,
            "scene": ""
        }
    ],
    "variables": [
        {
            "type": "path",
            "name": "the_id",
            "key": "id"
        }
    ]
}'
```

#### Execute the mock for the path set in the previously created mock
```sh
curl --location --request GET 'http://localhost:8080/mock-service/mock/users/123'
```

This example will return the following response:
```json
{
    "user_id": "123"
}
```

Alternatively, mocks can be created and edited [via UI](http://localhost:8080/mock-service/admin/#/)

![UI](https://raw.githubusercontent.com/nicopozo/mockserver/master/assets/ui.png)
Thanks [hecekiel](https://github.com/hecekiel) for your artwork!
