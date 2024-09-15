# To-Do List API

This is a simple To-Do List API built with Go, PostgreSQL, and Docker. The API allows you to create to-do lists, add items to them, mark items as completed, and delete items.

## Features

- Create to-do lists.
- Add items to lists.
- Remove items from lists.
- Mark items as completed.
- Get the list along with its items.

## Project Structure

```
todo-api/
├── docker-compose.yml   # Docker Compose configuration
├── Dockerfile           # Dockerfile for Go app
├── go.mod               # Go module definition
├── go.sum               # Go dependencies
├── main.go              # Go source code for the API
└── sql/
    └── init.sql         # PostgreSQL initialization script
```

## Requirements

Before you begin, make sure you have the following installed:

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Go (if you want to develop locally)](https://golang.org/doc/install)

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/Luke1491/to-do-list-api.git
cd to-do-list-api
```

### 2. Running the Application

The easiest way to run the application is by using Docker and Docker Compose.

#### Using Docker Compose

```bash
docker-compose up --build
```

This will:

- Start a PostgreSQL container and run the `init.sql` script to create the necessary database and tables.
- Start the Go application.

The API will be accessible at `http://localhost:8080`.

### 3. API Endpoints

Here are the available API endpoints:

#### 3.1 Create a To-Do List

**Endpoint**: `POST /lists`

**Request Body**:

```json
{
  "name": "Groceries"
}
```

**Response**:

```json
{
  "id": "uuid-of-list",
  "name": "Groceries"
}
```

#### 3.2 Add an Item to a To-Do List

**Endpoint**: `POST /items`

**Request Body**:

```json
{
  "list_id": "uuid-of-list",
  "description": "Buy Milk"
}
```

**Response**:

```json
{
  "id": "uuid-of-item",
  "list_id": "uuid-of-list",
  "description": "Buy Milk",
  "is_checked": false
}
```

#### 3.3 Get a To-Do List (with its items)

**Endpoint**: `GET /lists/{id}`

**Response**:

```json
{
  "list": {
    "id": "uuid-of-list",
    "name": "Groceries"
  },
  "items": [
    {
      "id": "uuid-of-item",
      "list_id": "uuid-of-list",
      "description": "Buy Milk",
      "is_checked": false
    }
  ]
}
```

#### 3.4 Update an Item (Mark as Completed)

**Endpoint**: `PUT /items/{id}`

**Request Body**:

```json
{
  "description": "Buy Milk",
  "is_checked": true
}
```

**Response**: `200 OK`

#### 3.5 Delete an Item

**Endpoint**: `DELETE /items/{id}`

**Response**: `200 OK`

### 4. Environment Variables

The following environment variables are used to configure the application:

- `DB_HOST`: The host of the PostgreSQL database (default: `db` in Docker Compose).
- `DB_PORT`: The port of the PostgreSQL database (default: `5432`).
- `DB_USER`: The username for the PostgreSQL database.
- `DB_PASSWORD`: The password for the PostgreSQL database.
- `DB_NAME`: The name of the PostgreSQL database.

You can modify these in the `docker-compose.yml` file if necessary.

### 5. Database Schema

The `init.sql` file initializes the PostgreSQL database with the following schema:

```sql
CREATE TABLE IF NOT EXISTS todo_lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS todo_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    list_id UUID REFERENCES todo_lists(id) ON DELETE CASCADE,
    description VARCHAR(255) NOT NULL,
    is_checked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### 6. Stopping the Application

To stop and remove all containers:

```bash
docker-compose down
```

### 7. Development

If you want to develop the application without Docker, follow these steps:

1. Make sure you have PostgreSQL running on your local machine.
2. Update the environment variables in the `.env` file or directly in your Go code.
3. Run the application using Go:
   ```bash
   go run main.go
   ```

### 8. License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
