# Contributing to Timelith

Thank you for your interest in contributing to Timelith! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for all contributors.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with:
- Clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, Docker version, etc.)
- Relevant logs or screenshots

### Suggesting Features

Feature requests are welcome! Please:
- Check existing issues first
- Provide clear use case
- Explain expected behavior
- Consider implementation complexity

### Pull Requests

1. **Fork the repository**
   ```bash
   git clone https://github.com/yourusername/timelith.git
   cd timelith
   ```

2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes**
   - Write clean, readable code
   - Follow existing code style
   - Add comments where necessary
   - Update documentation if needed

4. **Test your changes**
   ```bash
   # Rails tests
   cd rails-app
   bundle exec rspec

   # Go tests
   cd go-backend
   go test ./...
   ```

5. **Commit your changes**
   ```bash
   git add .
   git commit -m "Add feature: description"
   ```

6. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Create Pull Request**
   - Provide clear description
   - Reference related issues
   - Explain testing performed

## Development Setup

### Prerequisites
- Docker and Docker Compose
- Git
- Text editor or IDE

### Local Development

1. **Clone and setup**
   ```bash
   git clone https://github.com/yourusername/timelith.git
   cd timelith
   cp .env.example .env
   # Edit .env with your credentials
   ```

2. **Start services**
   ```bash
   docker-compose up -d
   ```

3. **Run migrations**
   ```bash
   docker-compose exec rails-app bundle exec rails db:migrate
   ```

4. **Access application**
   - Rails UI: http://localhost:3000
   - Go API: http://localhost:8080

### Code Style

**Ruby/Rails:**
- Follow Ruby Style Guide
- Use 2 spaces for indentation
- Keep methods short and focused
- Write descriptive variable names

**Go:**
- Run `go fmt` before committing
- Follow Go Code Review Comments
- Use meaningful variable names
- Add godoc comments for public functions

**General:**
- Write self-documenting code
- Add comments for complex logic
- Keep functions small and focused
- Use consistent naming conventions

## Project Structure

```
timelith/
├── rails-app/          # Rails Web UI
│   ├── app/           # Controllers, models, views
│   ├── config/        # Configuration
│   └── db/            # Migrations and schema
├── go-backend/        # Go Backend
│   ├── cmd/           # Main entry point
│   ├── internal/      # Internal packages
│   └── pkg/           # Public packages
├── docker/            # Docker configurations
└── docs/              # Additional documentation
```

## Testing

### Rails Tests
```bash
cd rails-app
bundle exec rspec
```

### Go Tests
```bash
cd go-backend
go test ./... -v
```

### Integration Tests
```bash
# Run full integration test suite
./scripts/integration-tests.sh
```

## Documentation

When adding features:
- Update README.md if needed
- Add inline code comments
- Update API documentation
- Add examples if applicable

## Release Process

1. Update CHANGELOG.md
2. Bump version in relevant files
3. Create git tag
4. Build and test Docker images
5. Create GitHub release

## Questions?

Feel free to:
- Open an issue for questions
- Join discussions
- Reach out to maintainers

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

Thank you for contributing to Timelith! 🎉
