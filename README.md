# Manga Backend

## Overview

This repository contains the backend implementation for the Manga Website, developed in GoLang. The backend serves as a
central hub for managing manga data, tracking updates from various manga websites, handling user interactions such as
bookmarks and notifications, and providing API endpoints for the frontend to consume.

The primary goals of this backend are:

- To provide users with a single place to monitor updates for their favorite mangas across multiple websites.
- To store manga websites in a database and periodically check for updates (once per hour).
- To apply tags to mangas, enabling search functionality.
- To notify users when a manga receives an update.
- To bookmark the last read chapter for users.
- To link to original manga websites with a notice encouraging users to read on the source site.
- To estimate the time until the next chapter release, based on historical release patterns (with potential AI
  integration in the future).

**Note:** This backend focuses exclusively on server-side logic and API services. Frontend development is handled in a
separate repository (https://github.com/Sidler1/manga-backend/tree/master – note: this may be a placeholder; adjust as
needed).

## Tech Stack

- **Language:** GoLang (Go 1.20+ recommended)
- **Database:** PostgreSQL (or another relational database; configurable via environment variables)
- **Key Libraries:**
    - Gorilla Mux for routing (or Gin/Echo if preferred)
    - GORM for ORM and database interactions
    - Go-Cron for scheduled tasks (e.g., hourly update checks)
    - JWT for authentication
    - Other dependencies as listed in `go.mod`
- **Scraping:** Custom scrapers for manga websites (implemented with GoQuery or similar)
- **Notifications:** Integration with services like Firebase Cloud Messaging or email (e.g., via SMTP)
- **Deployment:** Docker for containerization, with support for Kubernetes or cloud platforms like AWS/Heroku

## Project Structure

Assuming a standard Go project layout (to be implemented):

```
manga-backend/
├── cmd/
│   └── main.go          # Entry point for the application
├── internal/
│   ├── api/             # API handlers and routes
│   ├── models/          # Database models (e.g., Manga, Website, Chapter, Tag, User)
│   ├── services/        # Business logic (e.g., update checker, notification sender)
│   ├── repository/      # Database interactions
│   ├── scraper/         # Website scraping logic
│   └── utils/           # Utilities (e.g., JWT, cron jobs)
├── pkg/                 # Reusable packages (if any)
├── config/              # Configuration files (e.g., env.example)
├── migrations/          # Database migration scripts
├── docker/              # Dockerfiles and compose files
├── go.mod               # Go modules
├── go.sum               # Go module checksums
└── README.md            # This file
```

## Setup and Installation

1. **Clone the Repository:**
   ```
   git clone https://github.com/Sidler1/manga-backend.git
   cd manga-backend
   ```

2. **Install Dependencies:**
   ```
   go mod tidy
   ```

3. **Environment Configuration:**
    - Copy `.env.example` to `.env` and fill in the required variables (e.g., DB credentials, JWT secret, scraping
      intervals).
    - Example `.env`:
      ```
      DB_HOST=localhost
      DB_PORT=5432
      DB_USER=postgres
      DB_PASSWORD=secret
      DB_NAME=manga_db
      JWT_SECRET=your_jwt_secret
      UPDATE_INTERVAL=1h  # For cron job
      ```

4. **Database Setup:**
    - Set up a PostgreSQL database.
    - Run migrations (using a tool like Goose or built-in GORM auto-migrate):
      ```
      go run cmd/migrate.go
      ```

5. **Run the Application:**
   ```
   go run cmd/main.go
   ```
   The server will start on `http://localhost:8080` (configurable).

6. **Docker Setup (Optional):**
   ```
   docker-compose up -d
   ```

7. **Testing:**
   ```
   go test ./...
   ```

## API Endpoints

The backend exposes RESTful APIs. Base URL: `/api/v1`

### Authentication

- All endpoints require JWT authentication (except public ones like health check).
- `/auth/login` – POST: User login to get JWT.
- `/auth/register` – POST: User registration.

### Mangas

- `GET /mangas` – Retrieve all mangas (with pagination and tag filters).
- `GET /mangas/{id}` – Get details for a specific manga.
- `POST /mangas` – Add a new manga (admin only).
- `PUT /mangas/{id}` – Update manga details.
- `DELETE /mangas/{id}` – Delete a manga.

### Websites

- `GET /websites` – List all tracked manga websites.
- `POST /websites` – Add a new website for scraping.
- `PUT /websites/{id}` – Update website details.
- `DELETE /websites/{id}` – Remove a website.

### Chapters

- `GET /mangas/{manga_id}/chapters` – Get chapters for a manga.
- `POST /mangas/{manga_id}/chapters` – Add a new chapter (typically automated via scraper).
- `GET /mangas/{manga_id}/estimate-next` – Get estimated time to next chapter.

### Tags

- `GET /tags` – List all tags.
- `POST /tags` – Add a new tag.
- `PUT /tags/{id}` – Update a tag.
- `DELETE /tags/{id}` – Delete a tag.

### Users

- `GET /users/favorites` – Get user's favorite mangas.
- `POST /users/favorites` – Add a favorite.
- `GET /users/bookmarks/{manga_id}` – Get bookmark for a manga.
- `PUT /users/bookmarks/{manga_id}` – Update bookmark.
- `GET /users/notifications` – Get user notifications.

### Other

- `GET /health` – Health check endpoint.
- Scheduled Job: Runs hourly to scrape websites for updates, update DB, and send notifications.

For detailed request/response schemas, refer to the OpenAPI/Swagger docs (to be generated at `/api/docs`).

## Scheduled Tasks

- **Update Checker:** Uses cron to run every hour, scraping configured websites for new chapters.
- **Notification Sender:** Triggers push/email notifications on updates.
- **Estimation Logic:** Calculates average release interval from chapter history; future enhancements may include
  ML-based predictions.

## Security and Best Practices

- Use HTTPS in production.
- Rate limiting on APIs.
- Input validation and sanitization to prevent injections.
- Logging with tools like Zap or Logrus.
- Error handling with standardized responses (e.g., { "error": "message" }).

## Contributing

- Fork the repo and create a pull request.
- Follow Go best practices (e.g., idiomatic code, tests).
- Report issues via GitHub Issues.

## License

MIT License (or specify as needed).

## Contact

For questions, reach out via GitHub issues or [your email/contact].