# SETA Training Project

This repository contains the source code for the SETA Golang/Node.js training project.

## Contents

- Source code (Go)
- `docker-compose.yml` for running the system
- This `README.md` explaining setup and API usage

---

## Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- (Optional) PostgreSQL and Redis if you want to run the backend without Docker

---

## Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/ppvan/seta_training.git
cd seta_training
```

### 2. Start the system with Docker Compose

```bash
docker-compose up --build
```

This will build and launch all services defined in `docker-compose.yml`.

By default, the API service runs on port `8000`.

---

## API Endpoints

Base URL: `http://localhost:8000`

### 1. Health Check

**GET** `/v1/healthcheck`

```bash
curl http://localhost:8000/v1/healthcheck
```

### 2. Create a Post

**POST** `/v1/posts`

```bash
curl -X POST http://localhost:8000/v1/posts \
  -H "Content-Type: application/json" \
  -d '{"title":"Hello World","content":"My first post","tags":["training","golang"]}'
```

### 3. Search Posts by Tag

**GET** `/v1/posts/search-by-tag?tag=training`

```bash
curl "http://localhost:8000/v1/posts/search-by-tag?tag=training"
```

### 4. Get a Post by ID

**GET** `/v1/posts/{id}`

Example (replace `{id}` with an actual post ID, e.g., `1`):

```bash
curl http://localhost:8000/v1/posts/1
```

---

## Notes

- If you need to change ports or service settings, edit `docker-compose.yml`.
- For troubleshooting, check logs in your terminal or use `docker-compose logs`.
- The backend requires a PostgreSQL and Redis instance as defined. Adjust credentials in your `.env` or `docker-compose.yml` if needed.

---

## Contributing

Contributions are welcome! Please fork the repo and submit a pull request.

---

## License

See `LICENSE` file for details.
