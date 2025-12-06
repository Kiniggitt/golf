# Golf Course Management System

A desktop application for managing golf course memberships, fees, and user accounts built with Go and Fyne.

## Quick Start

```bash
# Copy the example configuration
cp config.json.example config.json

# Edit config.json to set your admin password
# Then run the application
go run ./cmd/golf

# Run tests
go test ./...

# Build
go build -o golf ./cmd/golf
```

See [CONFIG.md](CONFIG.md) for detailed configuration options.

## Project Structure

```
golf/
├── cmd/golf/main.go          # Application entry point, UI logic
├── internal/
│   ├── membership/           # User and fee business logic
│   ├── registrar/            # JSON persistence
│   ├── admin/                # Admin table UI
│   ├── config/               # Configuration management
│   └── ui/                   # Reusable UI components
├── data/
│   ├── users.json            # User data
│   └── standard_fees.json    # Fee templates
├── config.json.example       # Example configuration
├── CONFIG.md                 # Configuration guide
└── docs/
    ├── ARCHITECTURE.md       # Detailed architecture
    ├── API.md                # Package API reference
    ├── GO_BEST_PRACTICES.md  # Coding guidelines
    └── REFACTORING_PLAN.md   # Improvement roadmap
```

## Key Features

- User registration with automatic username generation
- User sign-in and account viewing
- Admin dashboard with configurable password
- Standard fee templates
- Fee assignment to users
- JSON data persistence
- Configurable window sizes and data paths

## Architecture Overview

### Data Flow
1. **UI Layer** (cmd/golf/main.go) - User interactions
2. **Business Logic** (internal/membership) - User/fee management
3. **Data Access** (internal/registrar) - JSON read/write

### Key Types

**User**: Golf course member
- Username (auto-generated)
- Name (first, last, suffix)
- Address
- Balance + Fees

**StandardFee**: Reusable fee template
- ID (auto-generated from name)
- Name
- Amount

### Global State
```go
// internal/membership/membership.go
var idMap = make(map[string]*User)           // Username → User
var standardFees = make(map[string]*StandardFee)  // ID → Fee
```

⚠️ **Not thread-safe** (see REFACTORING_PLAN.md)

## Testing

```bash
# All tests
go test ./...

# With coverage
go test ./... -cover

# With race detection
go test -race ./...

# Coverage report
go test ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Current coverage: **85.8%**

## Known Issues

See [REFACTORING_PLAN.md](REFACTORING_PLAN.md) for prioritized fixes:

**Critical**:
- Error handling in registrar (all errors ignored)
- No concurrency protection on global maps
- Username generation bug (line 109 in membership.go)

**High Priority**:
- Missing godoc documentation
- Hardcoded admin password

## Documentation

- [ARCHITECTURE.md](ARCHITECTURE.md) - Detailed design, data flow, patterns
- [API.md](API.md) - Complete API reference for all packages
- [GO_BEST_PRACTICES.md](GO_BEST_PRACTICES.md) - Go coding guidelines with examples
- [REFACTORING_PLAN.md](REFACTORING_PLAN.md) - Prioritized improvement plan
- [TEST_COVERAGE.md](TEST_COVERAGE.md) - Test documentation

## Admin Access

1. Click "Admin" on home screen
2. Password: `admin123` (hardcoded in main.go line ~130)
3. Features:
   - View all users in table
   - Assign fees to users
   - Edit user information
   - Manage standard fees

## Data Persistence

Data auto-saves to `data/` directory after each change:
- `users.json` - All user accounts
- `standard_fees.json` - Fee templates

**Backup**: Copy the `data/` directory
