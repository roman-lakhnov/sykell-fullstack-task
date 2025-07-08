# Sykell Full-Stack Developer Test Task

This project is a test task for the Full-Stack Developer position at Sykell,
with a front-end focus. It should include a React TypeScript frontend and a
basic Go server.

## Server Go+Gin

Despite having zero experience with Go prior to this project, I developed the
simple server with Go, providing a RESTful API for analyzing web pages. It
uses the Gin framework for handling HTTP requests and in-docker MySQL for data storage.

### Features

- **Link Analysis**: Analyzes web pages and extracts metadata including:
  - HTML version
  - Page title
  - Heading counts (H1-H6)
  - Internal and external links
  - Inaccessible links detection
  - Login form detection

### API Endpoints

- `POST /links`: Add URLs for analysis
  - Request body:
    `{"urls": ["https://example.com", "https://another-site.com"]}`
- `GET /links`: Retrieve analyzed links with pagination
  - Query parameters: `amount` (default: 10), `page` (default: 1)
- `PUT /links`: Update link status
  - Request body: `{"id": 1, "status": "stop"}`
  - Status can be set to either "stop" or "created"

### Background Processing

The server includes a background worker that continuously processes links with
"created" status.
