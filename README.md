# LinkCollector

LinkCollector is a simple web application that allows users to save, organize and share useful links. Built with Go and designed to be easy to use and deploy.

## Features

- 🔐 User registration and login
- 🔗 Save links with title, description, and tags
- 🏷️ Tag-based organization
- 🔍 Search through your links
- 👥 (Future) Share links with other users

## Screenshot

![image](https://github.com/user-attachments/assets/40ad3fff-b065-43d4-88f0-c487f910198c)
![image](https://github.com/user-attachments/assets/819de5a4-e5fc-43cd-b2a7-873a99d1b650)
![image](https://github.com/user-attachments/assets/8a483e89-27b7-460f-8fd3-4061919265bd)
![image](https://github.com/user-attachments/assets/135c8514-8892-4d30-b519-3aa17d38b090)



## Installation

### Prerequisites

- Go 1.16 or later
- MySQL database (configured for Strato hosting)

### Database Setup

1. Create a MySQL database on your Strato hosting panel
2. Note your database credentials (host, username, password, database name)
3. Configure the application to use your database by setting environment variables:
   ```
   DB_USER=your_strato_db_user
   DB_PASSWORD=your_strato_db_password
   DB_HOST=your_strato_db_host
   DB_PORT=3306
   DB_NAME=your_strato_db_name
   ```

Alternatively, you can edit the defaults in `config.go` directly.

### Setup

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/LinkCollector.git
   cd LinkCollector
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Run the application:
   ```
   go run main.go config.go
   ```

4. Access the application in your browser at:
   ```
   http://localhost:8080
   ```


## Project Structure

```
.
├── main.go             # Main application file
├── config.go           # Configuration handling
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── static/             # Static assets
│   ├── css/            # CSS files
│   └── js/             # JavaScript files
├── templates/          # HTML templates
└── linkcollector.service # Systemd service file for deployment
```

## Future Improvements

- Password reset functionality
- Link sharing between users
- Public/private link settings
- Enhanced search with filters
- Import/export functionality
- Browser extension for easy link saving


## Author

parsleyfilip - your beloved
