# Mock Service

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
| `DYNAMO_TABLE_PREFIX` | Prefix for DynamoDB tables | `mockserver_` |
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
- Configure variables (Path, Query, Header, Body, Random, Hash).
- Define assertions (Equals, Regex, Contains, JSON Schema, etc.).
- View real-time logs of mocked requests.

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
