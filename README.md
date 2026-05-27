# Mock Service

<p align="center">
  <img src="web/src/assets/mock.png" alt="Mock Service Logo" width="200"/>
</p>

![technology Go](https://img.shields.io/badge/technology-go-blue.svg)

A simple and powerful mock server in Go with support for dynamic variables, assertions, and multiple persistence backends.

## Installation & Deployment

### Local with Docker

The easiest way to run the project is using Docker.

**Build the image:**

```sh
make docker-build
```

**Run with File persistence:**

```sh
docker run -v /tmp:/tmp -e MOCKS_FILE=/tmp/mocks.json -p 8080:8080 --name mock-service nicopozo/mock-service:latest
```

**Run with MySQL/PostgreSQL:**

```sh
docker run -e MOCKS_DATASOURCE=mysql -e MYSQL_URL=mysql://user:password@host:port/db_name -p 8080:8080 --name mock-service nicopozo/mock-service:latest
```

### AWS Deployment (Lambda & DynamoDB)

Mock Service is ready to be deployed as a Docker-based Lambda function behind API Gateway.

1. **Initialize Infrastructure:**

    ```sh
    make aws-create-role  # Creates IAM Role
    make aws-init-db     # Provisions DynamoDB tables
    ```

2. **Deploy:**

    ```sh
    make aws-lambda-full  # Build, Push to ECR and Deploy to Lambda
    make aws-enable-api-gateway # Configure API Gateway
    ```

### Development

**Compile locally:**

```sh
make build
./service
```

**Run tests:**

```sh
make test
```

## Configuration

| Environment Variable | Description | Default |
| --- | --- | --- |
| `MOCKS_DATASOURCE` | `file`, `mysql`, `postgres`, or `dynamo` | `file` |
| `MOCKS_FILE` | Path to JSON file (only for `file` mode) | `/tmp/mocks.json` |
| `MYSQL_URL` / `POSTGRES_URL` | Full connection string for SQL databases | |
| `DB_USER` | Database user (optional, used if URL not provided) | `root` |
| `DB_PASSWORD` | Database password | `password` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port (default 3306 for MySQL, 5432 for PostgreSQL) | `3306` |
| `DB_NAME` | Database name | `mockserver` |
| `DB_SSLMODE` | SSL mode for database connection | `disable` |
| `DYNAMO_TABLE_PREFIX` | Prefix for DynamoDB tables | `mockserver_` |
| `DYNAMO_ENDPOINT` | DynamoDB endpoint (useful for local DynamoDB like LocalStack) | |
| `AWS_REGION` | AWS Region for DynamoDB/Lambda | `us-east-1` |

## Versioning

We use a centralized versioning system. The version is stored in the `VERSION` file.

**Bump version:**

```sh
make bump PART=patch  # 3.5.1 -> 3.5.2 (default)
make bump PART=minor  # 3.5.2 -> 3.6.0
make bump PART=major  # 3.6.0 -> 4.0.0
```

## How to use it

### Administer your mocks via UI

Manage your mocks through the built-in administration panel.

**URL:** [http://localhost:8080/mock-service/admin/](http://localhost:8080/mock-service/admin/)

From the UI, you can:

- Create, edit, and delete mocks.
- Configure variables (Path, Query, Header, Body, XML, Random, Hash, Composite).
- Define assertions (Equals, Regex, Contains, JSON Schema, etc.).
- Configure **webhooks** per response to fire HTTP calls asynchronously.
- View real-time logs of mocked requests, including webhook results.

### Execute a mock

Once you have created a mock (e.g., path `/users/{id}`), you can execute it:

```sh
curl --location --request GET 'http://localhost:8080/mock-service/mock/users/123'
```

**Example Response:**

```json
{
    "user_id": "123"
}
```

### Webhooks

Each response can optionally fire an asynchronous **webhook** — an HTTP call to an external URL triggered after the mock response is returned.

**Configuration fields:**

| Field | Description |
| --- | --- |
| `url` | Target URL for the webhook call |
| `method` | HTTP method (`GET`, `POST`, `PUT`, `PATCH`, `DELETE`) |
| `headers` | Custom headers sent with the request (JSON object) |
| `body` | Request body, supports `{variable}` placeholders |
| `enabled` | Enable or disable the webhook |
| `delay` | Milliseconds to wait before firing the webhook |
| `timeout` | Maximum time (ms) to wait for the webhook response |

Webhooks are fired **asynchronously** — the mock response is returned immediately and the webhook runs in the background. Check the **Request Logs** page to see webhook results (status, duration, response body).
