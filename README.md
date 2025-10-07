# Newsletter Service - Stay Connected, Stay Informed!

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://example.com/build)
[![License](https://img.shields.io/badge/license-MIT-blue)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/go-1.23-blue)](https://go.dev/)

A simple newsletter service for managing subscribers and sending email updates.

## Overview

This project provides a basic newsletter management system allowing users to subscribe and receive updates. It leverages a PostgreSQL database for storing subscriber information and uses Go's standard libraries for core functionality. This service can be used as a foundation for building more complex communication platforms.

**Problem:** Manually managing newsletter subscribers and sending updates is inefficient and prone to errors.

**Solution:** This service automates the subscription and distribution process, ensuring reliable and timely delivery of newsletters.

**Target Audience:** Small businesses, content creators, and organizations looking to manage their email communications efficiently.

**Use Cases:**
*   Sending weekly updates to subscribers.
*   Distributing announcements about new products or services.
*   Sharing curated content with a targeted audience.

**Key Technologies:**
*   Go 1.23.6
*   PostgreSQL
*   `github.com/lib/pq` for PostgreSQL driver
*   `github.com/joho/godotenv` for environment variable management
*   `github.com/gomarkdown/markdown` for rendering markdown
*   Docker Compose for easy setup

**Unique Value:** Provides a streamlined and easily deployable solution for newsletter management, focusing on simplicity and core functionality.

## Features

*   **Subscriber Management:**
    *   Add new subscribers to the database.
    *   Remove subscribers from the database.
    *   Retrieve a list of all subscribers.
*   **Newsletter Creation:**
    *   Write newsletters in Markdown format.
    *   Convert Markdown to HTML for email delivery.
*   **Email Sending (Placeholder):**
    *   *Note: This project currently focuses on subscriber management and Markdown rendering. Email sending functionality would need to be implemented using an email service provider (e.g., SendGrid, Mailgun).*
*   **Database Integration:**
    *   Uses PostgreSQL for persistent storage of subscriber data.
*   **Configuration:**
    *   Environment variable-based configuration for easy deployment.
*   **Markdown Rendering:**
    *   Converts Markdown content to HTML using `github.com/gomarkdown/markdown` and `github.com/yuin/goldmark` for rich text formatting.
    *   Supports GitHub Flavored Markdown (GFM).
    *   Syntax highlighting via `github.com/yuin/goldmark-highlighting`.

## Quick Start

1.  Clone the repository:
    ```bash
    git clone https://github.com/pixperk/newsletter.git
    cd newsletter
    ```
2.  Start the PostgreSQL database using Docker Compose:
    ```bash
    docker-compose up -d
    ```
3.  *Note: The core functionality involves setting up the database.  Further development is needed to interact with the database and send newsletters.*

## Installation

### Prerequisites

*   Go 1.23.6 or higher
*   Docker and Docker Compose (for PostgreSQL setup)

### Steps

1.  **Clone the Repository:**
    ```bash
    git clone https://github.com/pixperk/newsletter.git
    cd newsletter
    ```

2.  **Set up PostgreSQL using Docker Compose:**
    ```bash
    docker-compose up -d
    ```
    This will start a PostgreSQL instance named `newsletter-db` with the database `newsletter`.

3.  **Configure Environment Variables:**
    Create a `.env` file in the project root and set the following variables:

    ```
    POSTGRES_HOST=localhost
    POSTGRES_PORT=5432
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
    POSTGRES_DB=newsletter
    ```

4.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```

5.  **Build the Application:**
    ```bash
    go build -o newsletter .
    ```

### Verification

1.  Ensure the PostgreSQL container is running:
    ```bash
    docker ps
    ```
    You should see a container named `newsletter-db`.

2.  Verify that you can connect to the PostgreSQL database using a client like `psql`:
    ```bash
    psql -h localhost -p 5432 -U postgres -d newsletter
    ```

## Usage

This project currently provides a Markdown to HTML conversion utility.  Further development is required to implement the core newsletter functionality.

### Markdown to HTML Conversion

The `utils/markdown.go` file provides a function to convert Markdown to HTML.

```go
package main

import (
	"fmt"
	"github.com/pixperk/newsletter/utils"
)

func main() {
	markdown := `# Hello, Newsletter!
This is a **test** newsletter.
`
	html := utils.MarkdownToHTML(markdown)
	fmt.Println(html)
}
```

**Output:**

```html
<h1>Hello, Newsletter!</h1>

<p>This is a <strong>test</strong> newsletter.</p>
```

### Environment Variables

The following environment variables are used to configure the application:

*   `POSTGRES_HOST`: Hostname of the PostgreSQL server (default: `localhost`).
*   `POSTGRES_PORT`: Port number of the PostgreSQL server (default: `5432`).
*   `POSTGRES_USER`: Username for connecting to the PostgreSQL database (default: `postgres`).
*   `POSTGRES_PASSWORD`: Password for connecting to the PostgreSQL database (default: `postgres`).
*   `POSTGRES_DB`: Name of the PostgreSQL database (default: `newsletter`).

## Architecture

The project follows a modular design:

*   **`main.go`:** Entry point of the application. (Currently a placeholder).
*   **`utils/markdown.go`:** Contains utility functions, specifically for Markdown to HTML conversion.
*   **`docker-compose.yml`:** Defines the PostgreSQL database service.

The intended data flow would involve:

1.  Subscribers are added to the PostgreSQL database.
2.  Newsletters are written in Markdown.
3.  The Markdown content is converted to HTML.
4.  (Future) The HTML content is sent to subscribers via an email service.

## Configuration

Configuration is primarily handled through environment variables. These variables are read at runtime to configure the database connection and other settings.

### Environment Variables

| Variable            | Description                                   | Default Value |
| ------------------- | --------------------------------------------- | ------------- |
| `POSTGRES_HOST`     | Hostname of the PostgreSQL server             | `localhost`   |
| `POSTGRES_PORT`     | Port number of the PostgreSQL server          | `5432`        |
| `POSTGRES_USER`     | Username for connecting to the PostgreSQL database | `postgres`    |
| `POSTGRES_PASSWORD` | Password for connecting to the PostgreSQL database | `postgres`    |
| `POSTGRES_DB`       | Name of the PostgreSQL database               | `newsletter`  |

### Configuration File

*This project does not currently use a configuration file.*

## Development

### Setup

1.  Install Go dependencies:
    ```bash
    go mod tidy
    ```

### Build

```bash
go build -o newsletter .
```

### Testing

*This project currently lacks automated tests.*

### Code Organization

*   `main.go`: Main application entry point.
*   `utils/`: Utility functions (e.g., Markdown conversion).

## Performance & Scaling

*This section is speculative as the core functionality is not yet implemented.*

*   Database performance will depend on the size of the subscriber list and the efficiency of database queries.
*   Consider using connection pooling to improve database performance.
*   Email sending performance will depend on the email service provider used.

## Troubleshooting

*   **Cannot connect to PostgreSQL:**
    *   Ensure the PostgreSQL container is running (`docker ps`).
    *   Verify the environment variables are correctly configured.
    *   Check the PostgreSQL logs for errors.
*   **Markdown conversion issues:**
    *   Ensure the Markdown content is valid.
    *   Check for any errors during Markdown parsing.

## Contributing

Contributions are welcome!

1.  Fork the repository.
2.  Create a new branch for your feature or bug fix.
3.  Write tests for your code.
4.  Submit a pull request.

## License & Legal

This project is licensed under the MIT License. See the `LICENSE` file for details.

### Third-Party Dependencies

This project uses the following third-party dependencies:

*   `github.com/lib/pq` (MIT License)
*   `github.com/joho/godotenv` (MIT License)
*   `github.com/gomarkdown/markdown` (BSD-3-Clause License)
*   `github.com/yuin/goldmark` (MIT License)
*   `github.com/yuin/goldmark-highlighting` (MIT License)

---

<div style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; font-size: 14px; line-height: 1.6; color: #333; margin-top: 32px; padding: 24px; border-top: 1px solid #e0e0e0; background-color: #fafafa;">
  
  <img src="https://www.pixperk.tech/assets/avatar.jpg" alt="PixPerk" width="48" height="48" style="border-radius: 50%; vertical-align: middle; margin-right: 14px; border: 2px solid #ffffff; box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);" />

  <strong style="font-size: 16px; color: #222;">PixPerk</strong><br/>
  <span style="font-weight: 500; color: #444;">Yashaswi</span><br/>
  <em style="color: #666;">Backend Developer • <a href="https://www.pixperk.tech" target="_blank" style="color: #0066cc; text-decoration: none; font-weight: 500;">pixperk.tech</a></em><br/>

  <a href="https://twitter.com/pixperk_" target="_blank" style="color: #1DA1F2; text-decoration: none; font-size: 14px;">🐦 @pixperk_</a>

  <hr style="margin: 20px 0; border: none; border-top: 1px solid #e0e0e0;" />

  <p style="margin: 0 0 12px 0; color: #555;">If you found this useful, share it with someone who builds.</p>

  <a href="https://www.buymeacoffee.com/pixperk" target="_blank">
    <img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;" >
  </a>

  <p style="margin-top: 12px;">
    <a href="https://pixperk.tech/?unsubscribe=true" target="_blank" style="color: #999; text-decoration: none; font-weight: 400;">Unsubscribe</a>
  </p>

  <p style="font-size: 12px; color: #999; margin-top: 24px; font-weight: 400;">
    © 2025 Yashaswi – All bytes reserved.
  </p>

</div>