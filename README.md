# Matchmaking System

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-008ECF?style=for-the-badge&logo=gin&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white)
![gRPC](https://img.shields.io/badge/gRPC-244C5A?style=for-the-badge&logo=google&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)

A rating-based matchmaking platform built as a set of independent Go microservices. Players can be created, ratings can be tracked with Elo updates, and the matchmaking service uses Redis queues plus gRPC calls to rating-service to match queued players with similar ratings.

## Architecture

```mermaid
flowchart LR
    Client --> PlayerService[Player Service]
    Client --> RatingService[Rating Service]
    Client --> MatchmakingService[Matchmaking Service]

    PlayerService --> PlayerDB[(Player DB)]
    RatingService --> RatingDB[(Rating DB)]

    MatchmakingService --> Redis[(Redis)]
    MatchmakingService -->|gRPC GetPlayerRating| RatingService
```

## Project Structure

Each service is independently organized with its own API, business logic, database access, configuration, Dockerfile, and tests.

```text
matchmaking-system/
|-- player-service/
|-- rating-service/
|-- matchmaking-service/
|-- proto/                 # gRPC definitions
|-- docker-compose.yml
`-- LICENSE
```

## Services

### Player Service

Responsible for:

* Player creation
* Player lookup by ID
* Player lookup by username
* Username normalization and validation
* Enforcing unique usernames

### Rating Service

Responsible for:

* Creating player ratings
* Fetching player ratings
* Recording match results
* Elo rating calculation
* Rating history
* Leaderboard generation
* Exposing rating lookup over gRPC

### Matchmaking Service

Responsible for:

* Ticket creation
* Ticket lookup
* Ticket cancellation
* Redis-backed queue management
* Background match creation
* Rating-based player matching
* Match lookup
* Calling rating-service over gRPC

---

## Tech Stack

### Backend

* Go
* Gin

### Database

* PostgreSQL
* Database-per-service architecture

### Cache / Queue

* Redis, including sorted sets for queue ordering

### Service Communication

* REST APIs
* gRPC
* Protocol Buffers

### Deployment

* Docker
* Docker Compose

---

## Design Notes

* Redis sorted sets are used for the matchmaking queue, with ticket IDs scored by creation timestamp.
* `player_ticket:{player_id}` prevents a player from having more than one active queue ticket.
* rating-service owns Elo calculation and rating history.
* matchmaking-service calls rating-service over gRPC to fetch ratings during match selection.
* The matchmaker checks the 10 oldest queued tickets and selects the closest rated pair from that window.
* Rating updates use database transactions and row locks to avoid conflicting match result updates.

---

## Matchmaking Flow

```mermaid
sequenceDiagram
    participant Client
    participant Matchmaking
    participant Redis
    participant Rating

    Client->>Matchmaking: POST /matchmaking/tickets
    Matchmaking->>Redis: Store ticket JSON
    Matchmaking->>Redis: Add ticket ID to sorted set scored by timestamp

    Matchmaking->>Redis: Read 10 oldest queued tickets
    Matchmaking->>Rating: Get player ratings over gRPC
    Rating-->>Matchmaking: Return ratings
    Matchmaking->>Matchmaking: Select closest rated pair
    Matchmaking->>Redis: Store match
    Matchmaking->>Redis: Mark tickets as matched
```

---

## Match Lifecycle

```mermaid
stateDiagram-v2
    [*] --> queued
    queued --> matched
    queued --> cancelled
    matched --> [*]
    cancelled --> [*]
```

---

## Rating Flow

```mermaid
sequenceDiagram
    participant Client
    participant Rating
    participant RatingDB

    Client->>Rating: POST /ratings/matches/:match_id/result
    Rating->>RatingDB: Lock both player ratings
    Rating->>Rating: Calculate Elo changes
    Rating->>RatingDB: Store match result
    Rating->>RatingDB: Update winner and loser ratings
```

---

## Getting started

### 1. Prerequisites

- Go 1.26+
- PostgreSQL
- Redis
- Docker and Docker Compose

### 2. Run with Docker Compose

Docker Compose does not use the service `.env` files. The required environment variables are already defined in `docker-compose.yml`.

From the root directory:

```bash
docker compose up --build
```

Docker Compose starts:

```text
player-service:      http://localhost:8080
rating-service:      http://localhost:8081
matchmaking-service: http://localhost:8082
rating gRPC:         localhost:9091
player-db:           localhost:5433
rating-db:           localhost:5434
redis:               localhost:6379
```

### 3. Run locally

Start PostgreSQL and Redis first. Then configure local environment variables.

Each service loads its own `.env` file at startup and exits if a required variable is missing.

**player-service/.env**

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5433/player_service?sslmode=disable
```

**rating-service/.env**

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5434/rating_service?sslmode=disable
```

**matchmaking-service/.env**

```env
REDIS_URL=redis://localhost:6379/0
RATING_GRPC_URL=localhost:9091
```

> `RATING_GRPC_URL` must point to the rating-service gRPC server.

Apply migrations:

```bash
psql "postgres://postgres:postgres@localhost:5433/player_service?sslmode=disable" -f player-service/migrations/001_init.sql
psql "postgres://postgres:postgres@localhost:5434/rating_service?sslmode=disable" -f rating-service/migrations/001_init.sql
psql "postgres://postgres:postgres@localhost:5434/rating_service?sslmode=disable" -f rating-service/migrations/002_matches.sql
```

Then run each service in a separate terminal:

```bash
cd player-service       && go run ./cmd/server   # :8080
cd rating-service       && go run ./cmd/server   # :8081 and :9091
cd matchmaking-service  && go run ./cmd/server   # :8082
```

---

## Player Endpoints

```http
POST /players
GET  /players/:id
GET  /players/by-username/:username
```

---

## Rating Endpoints

```http
POST /ratings/players/:player_id
GET  /ratings/players/:player_id

POST /ratings/matches/:match_id/result
GET  /ratings/matches/:match_id

GET  /ratings/leaderboard
GET  /ratings/players/:player_id/matches
```

### gRPC

```text
RatingService.GetPlayerRating
```

---

## Matchmaking Endpoints

```http
POST   /matchmaking/tickets
GET    /matchmaking/tickets/:ticket_id
DELETE /matchmaking/tickets/:ticket_id

GET    /matchmaking/queue

POST   /matchmaking/matches
GET    /matchmaking/matches/:match_id
```

`POST /matchmaking/matches` manually creates one match from the current queue. The background matchmaker also creates matches automatically.

---

## Request Examples

### Create Player

```bash
curl -X POST http://localhost:8080/players \
  -H "Content-Type: application/json" \
  -d '{"username":"player_one"}'
```

### Create Player Rating

```bash
curl -X POST http://localhost:8081/ratings/players/<player_id>
```

### Create Matchmaking Ticket

```bash
curl -X POST http://localhost:8082/matchmaking/tickets \
  -H "Content-Type: application/json" \
  -d '{"player_id":"<player_id>"}'
```

Example response:

```json
{
  "ticket_id": "ticket-uuid",
  "player_id": "player-uuid",
  "status": "queued",
  "created_at": "timestamp"
}
```

When a ticket is matched, `GET /matchmaking/tickets/:ticket_id` returns the same ticket with `status: "matched"` and a `match_id` field. That `match_id` is used when reporting the final result to rating-service.

Example matched ticket:

```json
{
  "ticket_id": "ticket-uuid",
  "player_id": "player-uuid",
  "status": "matched",
  "match_id": "match-uuid",
  "created_at": "timestamp"
}
```

Both matched tickets return the same `match_id`, so either player's ticket can be polled to discover the match.

The match can then be fetched from matchmaking-service:

```bash
curl http://localhost:8082/matchmaking/matches/<match_id>
```

Example match response:

```json
{
  "match_id": "match-uuid",
  "player_a_id": "player-a-uuid",
  "player_b_id": "player-b-uuid",
  "ticket_a_id": "ticket-a-uuid",
  "ticket_b_id": "ticket-b-uuid",
  "created_at": "timestamp"
}
```

### Report Match Result

The winner is supplied by the caller. The system records the result and rating-service updates both ratings.

```bash
curl -X POST http://localhost:8081/ratings/matches/<match_id>/result \
  -H "Content-Type: application/json" \
  -d '{"winner_id":"<winner_id>","loser_id":"<loser_id>"}'
```

Example response:

```json
{
  "id": "match-uuid",
  "winner_id": "winner-player-uuid",
  "loser_id": "loser-player-uuid",
  "winner_rating_before": 1500,
  "loser_rating_before": 1500,
  "winner_rating_after": 1516,
  "loser_rating_after": 1484,
  "created_at": "timestamp"
}
```

The example rating values are illustrative.

---

## Try It

1. Create two players with `POST /players`.
2. Create ratings for both players with `POST /ratings/players/:player_id`.
3. Queue both players with `POST /matchmaking/tickets`.
4. Wait for the background matchmaker to create a match.
5. Poll `GET /matchmaking/tickets/:ticket_id` until the ticket status is `matched`.
6. Use the returned `match_id` to report the result with `POST /ratings/matches/:match_id/result`.
7. Check updated ratings with `GET /ratings/players/:player_id` or `GET /ratings/leaderboard`.

```bash
curl http://localhost:8081/ratings/leaderboard
```

---

## Tests

Run tests for each service:

```bash
cd player-service       && go test ./...
cd rating-service       && go test ./...
cd matchmaking-service  && go test ./...
```

---

## Security Notes

Input validation:

* Username validation
* UUID validation

Data integrity:

* Database-per-service isolation
* Rating updates use database transactions
* Rating rows are locked during match result updates
* Match result constraints prevent invalid winner and loser combinations

**Note: The services do not currently authenticate requests.** They should only be run in a trusted local or private development environment. Do not expose the services publicly without authentication, authorization, rate limiting, and transport security.

---

## Known Limitations

* No authentication or authorization yet
* Matchmaking only considers the 10 oldest queued tickets
* Matchmaking depends on rating-service being available
* gRPC traffic is plaintext in the local setup
* Matchmaking Redis keys do not expire
* No retry, timeout, or circuit breaker around gRPC calls yet
* Only one matchmaking worker should run at a time until Redis locking is added

---

## Future Improvements

* Authentication and authorization
* API gateway
* OpenAPI documentation
* Health check endpoints
* Structured logging
* Metrics and observability
* Redis locking for multi-worker matchmaking
* Rating range expansion based on queue wait time
* Pagination for history endpoints
* Shared proto module
* TLS for gRPC
* Kubernetes deployment
* CI/CD pipeline

---

## License

This project is licensed under the [MIT License](LICENSE).
