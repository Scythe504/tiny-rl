Here’s a complete, polished **README.md** for your TinyRL backend including Docker setup, MaxMind instructions, dev/prod workflows, and project notes:

````markdown
# Project: github.com/scythe504/tiny-rl

**TinyRL Backend** – A simple URL shortener with analytics tracking.

---

## Getting Started

These instructions will help you get a copy of the project running locally for development and testing.

---

## Prerequisites

* [Docker](https://www.docker.com/get-started)
* [Docker Compose](https://docs.docker.com/compose/install/)
* Optional: Go 1.24+ if you want to run locally without Docker
* Local Dependencies: 
  (https://github.com/pressly/goose) # Necessary for migration.
  (https://github.com/air-verse/air) # Necessary for hot reloading.

---

## Setting Up with Docker

### 1. Clone the repository

```bash
git clone https://github.com/scythe504/tiny-rl.git
cd tiny-rl
````

### 2. Create `.env` file

Copy the example:

```bash
cp .env.example .env
```

Edit `.env` with your preferred configuration:

```dotenv
PORT=8080
APP_ENV=local
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=postgres
DB_USERNAME=postgres
DB_PASSWORD=mysecretpassword
DB_SCHEMA=public
MAXMIND_ACCOUNT_ID=your-account-id
MAXMIND_LICENSE_KEY=your-license-key
POSTGRES_CONN_URL=postgresql://postgres:mysecretpassword@localhost:5432/postgres
```

### 2b. Configure MaxMind GeoIP

TinyRL uses the MaxMind GeoLite2 database for country lookups. You need an **account ID** and **license key** from MaxMind.

1. Sign up at [MaxMind](https://www.maxmind.com/en/geolite2/signup) to get your credentials.
2. Add them to your `.env` file (see above).
3. The backend automatically downloads the GeoLite2-Country.mmdb file on startup.

   * It will be stored in `/app/data/` inside the container.
   * It will be stored in `${PWD}/data/`, if running locally without Docker.

---

### 3. Start services

```bash
docker compose up -d
```

* This starts:

  * **Postgres** for persistent storage
  * **Backend** API server
* Air (hot reload) is enabled in dev if you’re using the `Dockerfile.dev` setup.

---

### 4. Run migrations & seed data

**Option A — migrations automatically run on startup**
The backend runs `goose` migrations automatically when the container starts.

**Option B — manually seed**

```bash
docker compose run --rm seed
```

---

### 5. Access API

* Base URL: `http://localhost:8080`
* Example endpoints:

  * `GET /` – Hello world
  * `GET /{shortCode}` – Redirect to full URL
  * `POST /api/shorten` – Shorten a URL
  * `POST /api/update-link` – Update destination URL
  * Analytics endpoints under `/api/analytics/{shortCode}/...`

---

## Running in Development

* Mount source code and use Air for hot reload:

```bash
docker compose -f docker-compose.yml up backend
```

* Changes to `.go` files will trigger automatic rebuilds.

---

## Running in Production

* Use the `prod` stage of `Dockerfile.dev` (or your production Dockerfile)
* The binary `tiny-rl` is prebuilt and runs without Go installed

```bash
docker compose -f docker-compose.yml up --build backend
```

---

## Project Structure

```
cmd/
    api/           # main backend server entrypoint
    seed/          # seed script for initial data
data/
    COPYRIGHT.txt
    GeoLite2-Country.mmdb
    LICENSE.txt
internal/          # core packages
migrations/        # database migrations
scripts/           # scripts to download GeoIP DB and run migrations
Dockerfile.dev     # dev + prod multi-stage build
docker-compose.yml # Docker compose setup
```

---

## Notes

* Ensure the Postgres DB is healthy before running backend or seed containers.
* `{shortCode}` route should be registered **after `/api/...` routes** to avoid accidental route collisions.
* Air hot reload stores temporary files in `/tmp` (or `air_tmp` volume) for faster rebuilds.
* MaxMind GeoIP requires a valid Account ID and License Key. The backend downloads the database automatically if credentials are provided.