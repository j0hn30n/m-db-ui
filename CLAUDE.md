# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a MongoDB management tool built with Go language, providing a web-based UI for MongoDB database administration. The application uses the Gin web framework and MongoDB Go driver to provide database management capabilities.

## Development Commands

### Building and Running
```bash
# Install dependencies
go mod download

# Run development server
go run main.go

# Build production binary
go build -o m-db-ui main.go

# Run tests
go test ./...
```

### Environment Setup
```bash
# Copy environment configuration
cp .env.example .env

# Set MongoDB connection URI
export MONGO_URI="mongodb://localhost:27017"

# Set server port (default: 8082)
export PORT="8082"
```

## Architecture

### Project Structure
```
m-db-ui/
├── main.go                    # Application entry point
├── config.yaml               # YAML configuration file
├── internal/                 # Internal packages
│   ├── config/              # Configuration management
│   ├── database/            # MongoDB operations and models
│   ├── handlers/            # HTTP request handlers
│   └── services/            # Business logic services
├── web/                     # Web assets
│   ├── static/             # CSS, JS, images
│   │   ├── css/style.css   # Custom styles
│   │   └── js/app.js       # Frontend JavaScript
│   └── templates/          # HTML templates
│       ├── base.html       # Base template
│       ├── index.html      # Database list page
│       ├── database.html   # Database view page
│       ├── collection.html # Collection view page
│       └── error.html      # Error page
└── README.md               # Project documentation
```

### Key Components

#### Database Layer (`internal/database/`)
- `database.go`: Core MongoDB operations (CRUD for databases, collections, documents)
- `models.go`: Data structures for database info, collections, documents, statistics
- `utils.go`: Helper functions for ObjectID handling and query parsing

#### HTTP Layer (`internal/handlers/`)
- `handlers.go`: HTTP request handlers for all API endpoints and web pages
- Supports both JSON API and HTML rendering

#### Configuration (`internal/config/`)
- Environment variable support
- YAML configuration file support
- Default values for development

#### Web UI
- Bootstrap 5 for responsive design
- JSONEditor for document editing
- AJAX-based interactions
- Real-time statistics and pagination

### API Endpoints

#### Database Management
- `GET /api/v1/databases` - List all databases
- `GET /api/v1/databases/{name}` - Get database details
- `DELETE /api/v1/databases/{name}` - Delete database

#### Collection Management
- `GET /api/v1/databases/{db}/collections` - List collections
- `POST /api/v1/databases/{db}/collections` - Create collection
- `DELETE /api/v1/databases/{db}/collections/{collection}` - Delete collection

#### Document Management
- `GET /api/v1/databases/{db}/collections/{collection}/documents` - List documents with pagination
- `POST /api/v1/databases/{db}/collections/{collection}/documents` - Create document
- `PUT /api/v1/databases/{db}/collections/{collection}/documents/{id}` - Update document
- `DELETE /api/v1/databases/{db}/collections/{collection}/documents/{id}` - Delete document
- `POST /api/v1/databases/{db}/collections/{collection}/query` - Query documents

#### Statistics
- `GET /api/v1/stats` - Server statistics and status

## Technology Stack

- **Backend**: Go 1.19+, Gin web framework
- **Database**: MongoDB 3.6+
- **Frontend**: Bootstrap 5, JavaScript ES6+, JSONEditor
- **Dependencies**:
  - `github.com/gin-gonic/gin` - Web framework
  - `go.mongodb.org/mongo-driver` - MongoDB driver
  - `github.com/gin-contrib/cors` - CORS support

## Development Notes

- Default port is 8082 to avoid conflicts with other services
- The application uses environment variables for configuration
- Web templates use Go's template syntax with custom helper functions
- JSON documents are displayed and edited using JSONEditor for better UX
- Pagination is implemented for document collections to handle large datasets
- CORS is configured for development (all origins allowed)

## Configuration Priority

1. Environment variables (highest priority)
2. YAML configuration file (`config.yaml`)
3. Default values (lowest priority)

## Testing

The application uses Go's built-in testing framework. Test files should be placed alongside the source files with `_test.go` suffix.

## Production Deployment

For production deployment:
1. Set appropriate environment variables
2. Configure CORS origins to specific domains
3. Use HTTPS/TLS termination
4. Set up proper MongoDB authentication
5. Configure connection pooling appropriately for your load