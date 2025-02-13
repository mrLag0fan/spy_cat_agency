# Spy Cat Agency - API

This project provides an API for managing cats and missions for Spy Cat Agency. The API is developed using Go (Golang) and works with a PostgreSQL database.

## Requirements

- Docker
- Docker Compose

## Running the Project

### Step 1: Clone the Repository

First, clone the repository:

```bash
git clone https://github.com/mrLag0fan/spy_cat_agency.git
cd spy-cat-agency
```

### Step 2: Set Up Containers

The project uses `docker-compose` to run services (app and database). Make sure Docker and Docker Compose are installed on your machine.

### Step 3: Start Containers

Run the following command:

```bash
docker-compose up --build
```

This command will build and start all necessary containers (app and database). Depending on your system, this process may take a few minutes.

- The app will be available at `http://localhost:8080`
- Swagger UI will be available at `http://localhost:8080/swagger/index.html`

### Step 4: Access Swagger UI

After the containers are up and running, you can access the Swagger UI to test the API.

Open your browser and go to:

```
http://localhost:8080/swagger/index.html
```

This will allow you to view the documentation and test all available API methods.

### Step 5: Stop Containers

To stop and remove the containers, use the following command:

```bash
docker-compose down
```

## API Overview

### API Structure

1. **Cats**
   - Create, update, delete cats
   - View the list of cats

2. **Missions**
   - Create, update, delete missions
   - Assign cats to missions
   - Manage mission targets

3. **Documentation**
   - The API is accessible through Swagger UI, where you can find all endpoints, parameters, and request examples.

## Configuration

All settings, including database connection, can be modified through environment variables in the `docker-compose.yml` or `config` files.

## Logging

All HTTP requests and responses are logged for debugging and monitoring purposes.

## Additional Resources

- [TheCatAPI](https://api.thecatapi.com/) for validating cat breeds.