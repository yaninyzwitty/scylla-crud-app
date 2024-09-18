# Scylla-Go-App

## Overview

`scylla-go-app` is a Go application demonstrating integration with ScyllaDB, a high-performance NoSQL database designed for high availability and scalability. This project showcases how to connect to ScyllaDB, perform CRUD operations, and manage database schemas using Go.

## Features

- **Database Connection**: Establishes a connection to ScyllaDB using the `gocql` and `gocqlx` libraries.
- **Schema Management**: Automatically creates keyspaces and tables if they do not exist.
- **Configuration Management**: Loads configuration settings from environment variables.

## Project Structure

- **`/configuration`**: Contains code for managing application configuration.
  - `config.go`: Defines the `Config` struct and functions to load configuration from environment variables.
  
- **`/database`**: Handles database interactions with ScyllaDB.
  - `database.go`: Contains functions for establishing a connection to ScyllaDB and managing schemas.
  - 
- **`/controller`**: Handles  incoming HTTP requests and returning appropriate HTTP responses..
  - `songs_controller.go`: Contains functions for handling those requests.

- **`/service`**: Has business logic of the application.
  - `songs_service.go`: Contains functions for handling the business logic of the application.
## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/yaninyzwitty/scylla-crud-app.git
   cd scylla-crud-app
