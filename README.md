# Mock Service

![technology Go](https://img.shields.io/badge/technology-go-blue.svg)

A simple mock server in Go.

## Installation

### With Dockers

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
docker run -e MOCKS_DATASOURCE=mysql -e DB_USER={{user}} -e DB_PASSWORD={{password}} -e DB_HOST={{host}} -e DB_PORT={{port}} -e DB_NAME={{db_name}} -p 8080:8080 --name mock-service mock-service
```

If you are deploying to a PaaS like Railway, you can alternatively use `MYSQL_URL` directly:

```sh
docker run -e MOCKS_DATASOURCE=mysql -e MYSQL_URL=mysql://user:password@host:port/db_name -p 8080:8080 --name mock-service mock-service
```

Database must contain a schema `mockserver` (or the one specified in `DB_NAME`/`MYSQL_URL`) with required tables. Follow[`this link`](https://github.com/nicopozo/mockserver/blob/master/scripts/init.sql "Init sql script") to get the creation script.

### By compiling with Go

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
export DB_NAME={{db_name}}
```

Alternatively, you can provide a full connection URL (e.g. if deploying to Railway):

```sh
export MOCKS_DATASOURCE=mysql
export MYSQL_URL="mysql://user:password@host:port/db_name"
```

and then run the app with `./mocks` command (run the [`init database script`](https://github.com/nicopozo/mockserver/blob/master/scripts/init.sql "Init sql script") before running the app).

## How to use it

### Administer your mocks via UI

The easiest way to manage your mocks is through the built-in administration panel.

**URL:** [http://localhost:8080/mock-service/admin/](http://localhost:8080/mock-service/admin/)

From the UI, you can:

- Create, edit, and delete mocks.
- Configure variables (Path, Query, Header, Body).
- Define assertions (Equals, Regex, Contains, JSON Schema, etc.).
- View real-time logs of mocked requests.

![UI](assets/ui.png)

### Execute a mock

Once you have created a mock (for example, with path `/users/{id}`), you can execute it by calling the following endpoint:

```sh
curl --location --request GET 'http://localhost:8080/mock-service/mock/users/123'
```

This will trigger the mock engine, which will:

1. Extract variables (like `id=123`).
2. Run assertions.
3. Replace variables in the response body.
4. Return the configured response.

**Example Response:**

```json
{
    "user_id": "123"
}
```
