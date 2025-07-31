# 🤖 Discord AI Tech News Bot

A Discord bot built with Go that provides AI-powered tech news updates and includes a REST API built with Gin framework.

## 🚀 Features

- **Discord Bot Integration**: Responds to messages in specific channels
- **Channel Restriction**: Only operates in the "🔥┃ai-tech-news" channel
- **REST API**: Built with Gin framework for external integrations
- **Health Monitoring**: Health check endpoints for monitoring
- **Webhook Support**: Ready for external webhook integrations
- **Graceful Shutdown**: Proper cleanup of resources on exit

## 📁 Project Structure

```
discord-ai-tech-news/
├── cmd/
│   └── main.go              # Application entry point
├── config/
│   └── config.go           # Configuration and environment loading
├── internal/
│   ├── bot/
│   │   └── discord_bot.go  # Discord bot initialization
│   ├── handler/
│   │   ├── discord/
│   │   │   └── message_handler.go  # Discord message handling
│   │   └── http/
│   │       └── routes.go   # HTTP routes and handlers
│   ├── repository/        # Data access layer
│   ├── service/           # Business logic services
│   └── usecase/
│       └── message_usecase.go  # Message processing logic
├── .air.toml              # Air hot reload configuration
├── .github/
│   └── workflows/         # GitHub Actions CI/CD
│       ├── ci.yml         # Main CI pipeline
│       └── lint.yml       # Linting workflow
├── .golangci.yml          # golangci-lint configuration
├── .env.example           # Environment variables template
├── Makefile               # Build automation
├── go.mod                # Go module dependencies
└── README.md             # This file
```

## 🛠️ Prerequisites

- Go 1.23.5 or higher
- Discord Bot Token
- Discord Server with appropriate permissions
- [Air](https://github.com/cosmtrek/air) for hot reload development (optional)
- [golangci-lint](https://golangci-lint.run/) for code linting (optional)
- [Make](https://www.gnu.org/software/make/) for build automation (optional)

## ⚙️ Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/meterai07/discord-ai-tech-news.git
   cd discord-ai-tech-news
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Install Air for hot reload (optional)**
   ```bash
   go install github.com/cosmtrek/air@latest
   ```

4. **Install development tools (optional)**
   ```bash
   # Install golangci-lint for linting
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   
   # Or use the Makefile
   make install-tools
   ```

5. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` and add your Discord bot token:
   ```env
   TOKEN=your_discord_bot_token_here
   APP_PORT=8080
   ```

6. **Build the application**
   ```bash
   go build ./cmd
   ```

## 🏃‍♂️ Running the Application

### Development Mode with Hot Reload (Recommended)
```bash
# Using Air directly
air

# Or using Makefile
make dev
```

### Development Mode (Standard)
```bash
# Using Go directly
go run ./cmd

# Or using Makefile
make run
```

### Production Mode
```bash
# Build first
go build -o discord-bot ./cmd
# Or: make build

# Run the binary
./discord-bot
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `TOKEN` | Discord Bot Token | - | ✅ |
| `APP_PORT` | HTTP server port | `8080` | ❌ |

### Discord Bot Setup

1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Create a new application
3. Go to the "Bot" section
4. Create a bot and copy the token
5. Enable necessary intents (Message Content Intent if needed)
6. Invite the bot to your server with appropriate permissions

### Required Permissions
- Send Messages
- Read Message History
- View Channels

## 📡 API Endpoints

The bot includes a REST API server with the following endpoints:

### Health Check
```
GET /health
```
Response:
```json
{
  "status": "healthy",
  "bot": "online"
}
```

### Bot Status
```
GET /
```
Response:
```json
{
  "message": "Discord AI Tech News Bot API",
  "status": "running"
}
```

### Webhook
```
POST /webhook
```
Response:
```json
{
  "message": "webhook received"
}
```

## 🎯 Bot Commands

The bot supports the following commands in the "🔥┃ai-tech-news" channel:

### News Commands
- `news`, `berita`, `tech`, `teknologi` - Get latest tech news
- `search <keyword>`, `cari <keyword>` - Search news (coming soon)

### General Commands  
- `hello`, `hi`, `halo` - Greet the bot
- `help`, `bantuan` - Show available commands
- `ping` - Check bot connection
- `status` - View bot status

*Note: The bot only responds in channels named "🔥┃ai-tech-news"*

## 🔨 Development

### Hot Reload Development

This project supports hot reload using Air for faster development:

1. **Install Air** (if not already installed):
   ```bash
   go install github.com/cosmtrek/air@latest
   ```

2. **Start development with hot reload**:
   ```bash
   air
   ```

3. **Air Configuration**: The project includes `.air.toml` for custom configuration

### Code Quality & Linting

This project uses [golangci-lint](https://golangci-lint.run/) for comprehensive linting:

1. **Run linter**:
   ```bash
   # Using golangci-lint directly
   golangci-lint run
   
   # Or using Makefile
   make lint
   ```

2. **Auto-fix issues**:
   ```bash
   # Using golangci-lint directly
   golangci-lint run --fix
   
   # Or using Makefile
   make lint-fix
   ```

3. **Configuration**: See `.golangci.yml` for linting rules

### Makefile Commands

The project includes a Makefile for common development tasks:

```bash
make help          # Show all available commands
make build         # Build the application
make run           # Run the application
make dev           # Start development with hot reload
make test          # Run tests
make test-coverage # Run tests with coverage report
make lint          # Run linter
make lint-fix      # Run linter with auto-fix
make fmt           # Format code
make tidy          # Tidy dependencies
make check         # Run all checks (lint, test, build)
make pre-commit    # Prepare code for commit
make clean         # Clean build artifacts
make install-tools # Install development tools
```

### Adding New Commands

1. Edit `internal/usecase/message_usecase.go` to add new command logic
2. The message handler already filters for the correct channel

### Adding New API Endpoints

1. Edit `internal/handler/http/routes.go` to add new routes
2. Follow the existing pattern for consistency

### Project Architecture

This project follows Clean Architecture principles:

- **cmd/**: Application entry points
- **config/**: Configuration management
- **internal/bot/**: Discord bot initialization
- **internal/handler/**: Input handlers (HTTP and Discord)
- **internal/usecase/**: Business logic
- **internal/service/**: External service integrations
- **internal/repository/**: Data access layer

### Development Tools

- **Air**: Hot reload for Go applications during development
- **Gin**: HTTP web framework with built-in middleware
- **DiscordGo**: Official Discord API wrapper for Go

## 🧪 Testing

```bash
# Run tests
go test ./...
# Or: make test

# Run tests with coverage
go test -cover ./...
# Or: make test-coverage

# Run tests with race detection
go test -race ./...
```

### Continuous Integration

The project includes GitHub Actions workflows for:

- **Linting**: Runs golangci-lint on every pull request
- **Testing**: Runs tests across different Go versions
- **Security**: Runs security scans with gosec
- **Coverage**: Uploads coverage reports to Codecov

Workflows are triggered on:
- Pull requests to `master` branch
- Direct pushes to `master` branch

## 📦 Dependencies

- [discordgo](https://github.com/bwmarrin/discordgo) - Discord API wrapper
- [gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [godotenv](https://github.com/joho/godotenv) - Environment variable loading

### Development Dependencies

- [air](https://github.com/cosmtrek/air) - Hot reload for Go applications
- [golangci-lint](https://golangci-lint.run/) - Fast Go linters runner

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run quality checks (`make pre-commit`)
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Code Quality Standards

- All code must pass `golangci-lint` checks
- Tests are required for new features
- Maintain test coverage above 80%
- Follow Go best practices and conventions
- Use meaningful commit messages

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Troubleshooting

### Common Issues

1. **Bot not responding**
   - Check if the bot is in the correct channel ("🔥┃ai-tech-news")
   - Verify the bot has proper permissions
   - Check the logs for any error messages

2. **Build errors**
   - Run `go mod tidy` to resolve dependencies
   - Ensure Go version is 1.23.5 or higher

3. **Environment variables not loading**
   - Make sure `.env` file exists in the project root
   - Check that `TOKEN` is properly set

4. **Hot reload not working**
   - Ensure Air is properly installed: `go install github.com/cosmtrek/air@latest`
   - Check if `.air.toml` configuration file exists
   - Try running `air -v` to verify Air installation

5. **Linting errors**
   - Install golangci-lint: `make install-tools`
   - Run linter to see issues: `make lint`
   - Auto-fix some issues: `make lint-fix`
   - Check `.golangci.yml` for configuration

6. **CI/CD issues**
   - Ensure all linting checks pass locally before pushing
   - Check GitHub Actions logs for detailed error messages
   - Make sure `.github/workflows/` files are properly configured

### Logs

The application provides detailed logging for debugging:
- Discord message processing
- HTTP server status
- Error messages

## 🚀 Deployment

### Docker (Optional)

Create a `Dockerfile`:
```dockerfile
FROM golang:1.23.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o discord-bot ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/discord-bot .
COPY --from=builder /app/.env .
CMD ["./discord-bot"]
```

### Environment-specific Configuration

For different environments, create separate `.env` files:
- `.env.development`
- `.env.staging`
- `.env.production`

## 📞 Support

If you have any questions or issues, please:
1. Check the troubleshooting section
2. Search existing issues on GitHub
3. Create a new issue with detailed information

---

**Made with ❤️ and Go**
