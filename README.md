# Instructions — Run the project

## Prerequisites
- Docker and docker-compose installed and running.
- Open a terminal in the project's source directory (the directory containing docker-compose.yml).

## Start
1. In the source directory run:
```bash
docker-compose up --build
```
2. Wait — the first startup can be slow because databases and Kafka must be created.

3. Open the frontend in your browser:
```
http://localhost:8081
```

## Stop and remove containers
From the same source directory run:
```bash
docker-compose down
```

Notes:
- To run in background, add `-d`: `docker-compose up --build -d`
- To follow logs: `docker-compose logs -f`
