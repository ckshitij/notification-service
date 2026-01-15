# Notification Service (`notify-srv`)

A robust, scalable, and observable notification service built with Go. It provides a centralized system for managing and sending notifications across multiple channels like Email, Slack, and In-App.

## Features

- **Multi-Channel Support**: Easily send notifications via Email, Slack, and In-App channels.
- **Dynamic Templates**: Utilizes Go's templating engine (`text/template`) for dynamic and version-controlled notification content.
- **RESTful API**: A clean and simple API for managing templates and sending notifications. (See `api/openapi.yaml` for the full specification).
- **Asynchronous Processing**: Leverages Kafka for queuing and processing notification requests asynchronously, ensuring high throughput and resilience.
- **Database Migrations**: Manages database schema changes cleanly using a dedicated migrator tool.
- **Observability**: Exposes application metrics in Prometheus format for easy monitoring and alerting.
- **Containerized**: Comes with a complete `docker-compose` setup for all dependencies, enabling a one-command local environment startup.

## Tech Stack

- **Language**: Go (v1.24)
- **Framework**: Chi (for routing) & Viper (for configuration)
- **Database**: MySQL 8.4
- **Message Broker**: Kafka
- **Cache**: Redis
- **Monitoring**: Prometheus & Grafana
- **Local Email Testing**: MailHog
- **Containerization**: Docker & Docker Compose

## Prerequisites

Before you begin, ensure you have the following installed:

- [Go](https://go.dev/doc/install) (version 1.24 or later)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Getting Started

Follow these steps to get the notification service running on your local machine.

### 1. Clone the Repository

```bash
git clone https://github.com/ckshitij/notify-srv.git
cd notify-srv
```

### 2. Start Dependencies

All required backing services (MySQL, Kafka, Redis, etc.) are defined in the `docker-compose.yml` file. Start them all in detached mode:

```bash
docker-compose -f deployments/docker-compose.yml up -d
```

This will start all services and expose their default ports. You can check the status with `docker-compose -f deployments/docker-compose.yml ps`.

### 3. Run Database Migrations

Once the database container is healthy, run the database migrations to set up the required tables and seed initial data.

```bash
go run cmd/migrator/main.go
```

### 4. Run the Application

Finally, start the main notification service application:

```bash
go run cmd/app/main.go
```

The service should now be running and connected to all its dependencies. By default, it runs on port `8000`.

### 5. Accessing Services

Once everything is running, you can access the various components:

- **Notification Service Swagger UI**: `http://localhost:8098/swagger/index.html`
- **Kafka UI**: `http://localhost:8081`
- **Grafana Dashboard**: `http://localhost:3000` (Login: `notif_admin` / `Grafana@123`)
- **MailHog (Email Viewer)**: `http://localhost:8025`
- **Prometheus Targets**: `http://localhost:9090`

## API Documentation

The API is documented using the OpenAPI specification. The definition can be found in `api/openapi.yaml`.

## Running Tests

To run the test suite, execute the following command from the project root:

```bash
go test ./...
```

## Project Structure

- `api/`: Contains the OpenAPI specification file.
- `cmd/`: Main application entry points.
  - `app/`: The main notification service server.
  - `migrator/`: The database migration tool.
- `config/`: Contains the application configuration file (`config.yml`).
- `deployments/`: Docker and `docker-compose` files for local development.
- `internal/`: All private application logic.
  - `config/`: Configuration loading logic.
  - `logger/`: Application logger setup.
  - `metrics/`: Prometheus metrics setup.
  - `pkg/`: Core business logic for notifications, templates, and senders.
  - `server/`: HTTP server setup, routing, and middleware.
- `migrations/`: SQL migration files (`.up.sql` and `.down.sql`).
- `static/`: Static assets (if any).


## Core Logic (`internal/pkg`)

The `internal/pkg` directory contains the core business logic for the notification service. It is organized into several sub-packages, each with a distinct responsibility.

### `template`
- **Purpose:** Manages notification templates.
- **Functionality:** Provides CRUD (Create, Read, Update, Delete) operations for templates, which are stored in a MySQL database. It exposes HTTP handlers for managing these templates via an API.

### `notification`
- **Purpose:** Orchestrates the creation, scheduling, and sending of notifications.
- **Functionality:** This is the central service. It receives requests to send notifications, retrieves the appropriate template using the `template` service, renders the content, and then dispatches it through various channels. It also manages the state of notifications (e.g., pending, sent, failed) in the database.

### `renderer`
- **Purpose:** Renders notification content from templates.
- **Functionality:** Takes a template and a set of data (variables) and uses Go's `text/template` package to produce the final content for a notification, such as an email body or a Slack message.

### `senders`
- **Purpose:** Handles the actual delivery of notifications to external services.
- **Functionality:** Implements a strategy pattern with a common `Sender` interface. Concrete implementations for different channels are provided:
    - `email`: Sends notifications via an email service.
    - `slack`: Sends notifications to a Slack channel.
    - `inapp`: Stores notifications to be displayed within a UI.

### `schedular`
- **Purpose:** Processes scheduled and stuck notifications.
- **Functionality:** A background worker that periodically queries the database for notifications that are due to be sent or have been stuck in a "sending" state for too long. It then enqueues them for processing by the `notification` service.