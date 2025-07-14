# Sykell Full-Stack Developer Test Task

This project is a test task for the Full-Stack Developer position at Sykell,
with a front-end focus. It should include a React TypeScript frontend and a
basic Go server.

## ⚙️ How to Run and Test

- Ensure you have **Docker, Docker Compose, Node.js, npm, and Go with Gin**
  installed on your system.
- Clone this repository.
- Run `docker-compose up -d` to set up MariaDB.
- **Wait about 10-15 seconds** for the database to initialize.
- Navigate to the server directory (`cd server`), then run the server `go run main.go` which will also create the initial database structure.
- From root navigate to the client directory (`cd client`), run `npm i`, then
  `npm run dev`.
- Open your browser `http://localhost:5173/` to view the client UI.
- Enter a few links for testing (I suggest as example adding 5 links`https://www.python.org/downloads/,https://www.python.org/downloads/,https://www.python.org/downloads/,https://www.python.org/downloads/,https://www.python.org/downloads/`
- You will see a dashboard with one pending link and others in "created" status
  (in queue).
- Server will analyze links one by one.
- You can click "Refresh" to update the data status.
- You can stop or run/rerun analysis for any of the links.
- You can also view the database structure and records via phpMyAdmin UI through
  the web interface `http://localhost:8081/`.

## ⚙️ Server Go+Gin

Despite having zero experience with Go prior to this project, I developed the
simple server with Go, providing a RESTful API for analyzing web pages. It uses
the Gin framework for handling HTTP requests and in-docker MySQL for data
storage.

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

## ⚙️ Client Vite+TS+React+Bootstrap

The client is built using modern front-end technologies: Vite, TypeScript,
React, and Bootstrap. It provides a user-friendly interface for interacting with
the web page analysis service.

### Features

- **URL Management**:

  - Input for adding multiple URLs (comma-separated)
  - URL validation
  - Preview of URLs before submission
  - Submission for analysis

- **Results Dashboard**:
  - Detailed table with all analysis results
  - Sortable columns for easy data comparison
  - Filter functionality for each column
  - Global search across all fields
  - Pagination with adjustable page size
  - Status indicators with color coding
  - Controls to start/stop analysis for each URL

### Components Structure

- **UrlManagement**: Handles URL input, validation, and submission
- **ResultsDashboard**: Displays analyzed links data with filtering and sorting
  - **ResultsTable**:
    - FilterHeader: Column-specific filtering
    - TableHeader: Sortable column headers
    - TableRows: Data display with action buttons
    - Pagination: Navigation between result pages

### State Management

The application uses React's built-in state management (useState, useEffect,
useCallback) to handle:

- Pagination state (current page, page size)
- Sorting (field and direction)
- Filtering (per-column and global search)
- API communication status

### User Experience

- Responsive design using Bootstrap
- Toast notifications for user feedback
- Visual indicators for link status
- Interactive table for easy data navigation
- Intuitive controls for managing link analysis

## ⚙️ Implementation Status

This table outlines the requirements from the task and their implementation
status in this project.

### Back-end Requirements

| Requirement                       | Status                              |
| --------------------------------- | ----------------------------------- |
| Go (Golang) with framework        | ✅ Implemented using Gin framework  |
| MySQL for data storage            | ✅ Implemented with in-docker MySQL |
| HTML version detection            | ✅ Implemented                      |
| Page title extraction             | ✅ Implemented                      |
| Heading counts (H1-H6)            | ✅ Implemented                      |
| Internal vs. external links count | ✅ Implemented                      |
| Inaccessible links detection      | ✅ Implemented                      |
| Login form detection              | ✅ Implemented                      |
| Background processing             | ✅ Implemented with worker pool     |
| API authorization mechanism       | ❌ Not implemented                  |

### Front-end Requirements

| Requirement                           | Status                                                                                                                                                  |
| ------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- |
| React with TypeScript                 | ✅ Implemented with Vite and React 19                                                                                                                   |
| Responsive design                     | ✅ Implemented using Bootstrap                                                                                                                          |
| URL Management: Add URLs              | ✅ Implemented with validation                                                                                                                          |
| URL Management: Start/stop processing | ✅ Implemented with status control buttons                                                                                                              |
| Results Dashboard: Sortable table     | ✅ Implemented with all required columns                                                                                                                |
| Results Dashboard: Pagination         | ✅ Implemented with adjustable page size                                                                                                                |
| Results Dashboard: Column filters     | ✅ Implemented for each column                                                                                                                          |
| Results Dashboard: Global search      | ✅ Implemented with fuzzy matching                                                                                                                      |
| Details View: Bar/donut charts        | ❌ Not implemented                                                                                                                                      |
| Details View: List of broken links    | ⚠️ Partially implemented - Server and DB support broken links tracking, but feature not available in client UI                                          |
| Bulk Actions: Re-run/delete URLs      | ❌ Not implemented                                                                                                                                      |
| Real-Time Progress: Status indicators | ⚠️ Partially implemented - Color-coded status indicators implemented, but real-time updates require manual refresh rather than WebSocket implementation |
| Automated front-end tests             | ❌ Not implemented                                                                                                                                      |
