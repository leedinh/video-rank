# Video Ranking Microservice

## Architecture Diagram
![alt text](https://github.com/leedinh/video-rank/blob/main/video-ranking.drawio.png)

## Setup

1. Install Docker and Docker Compose.
2. Clone the repository.
3. Run `docker-compose up --build`.

## API Endpoints

- **POST /api/interactions**: Submit an interaction (view, like, etc.).
- **GET /api/rankings**: Retrieve top videos (global or per-user).

## Usage

Visit `http://localhost:8080/swagger/index.html` for the documentaion of the APIs. 

## Testing

Run `go test ./...` to execute unit tests.