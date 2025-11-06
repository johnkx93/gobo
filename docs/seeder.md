# Database Seeder

A Go-based database seeder using [gofakeit](https://github.com/brianvoe/gofakeit) to generate realistic fake data for development and testing.

## Features

- 🎲 Generates realistic fake data using gofakeit
- 👥 Creates users with fake names, emails, and usernames
- 📦 Creates orders with random statuses, amounts, and notes
- 🔄 Configurable number of records
- 🗑️ Option to clear existing data before seeding
- 📊 Progress indicators for long-running seeds

## Prerequisites

The seeder uses your existing database connection from the `DATABASE_URL` environment variable, or falls back to:
```
postgres://postgres:postgres@localhost:5432/appdb?sslmode=disable
```

## Usage

### Using Make (Recommended)

```bash
# Seed with default values (50 users, 200 orders)
make seed

# Seed with custom values
make seed users=100 orders=500

# Clear existing data and reseed
make seed-clear users=20 orders=50
```

### Using Go Directly

```bash
# Default: 50 users, 200 orders
go run cmd/seeder/main.go

# Custom values
go run cmd/seeder/main.go -users=100 -orders=500

# Clear and reseed
go run cmd/seeder/main.go -clear -users=20 -orders=50
```

## Command Line Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-users` | 50 | Number of users to generate |
| `-orders` | 200 | Number of orders to generate |
| `-clear` | false | Clear existing data before seeding |

## Generated Data

### Users
- **Email**: Realistic fake emails (e.g., `johndoe@example.com`)
- **Username**: Unique usernames with numeric suffix
- **Password**: All users have password `password123` (bcrypt hashed)
- **First/Last Name**: Realistic fake names using gofakeit

### Orders
- **User ID**: Randomly assigned to created users
- **Order Number**: Format `ORD-YYYYMMDD-NNNNNN`
- **Status**: Random from `pending`, `completed`, `shipped`, `cancelled`
- **Total Amount**: Random between $10.00 - $1000.00
- **Notes**: 50% chance of having a random sentence

## Example Output

```bash
$ make seed users=10 orders=20
👥 Generating 10 users...
  ... 10/10 users created
✅ Created 10 users
📦 Generating 20 orders...
✅ Created 20 orders

🎉 Seeding completed successfully!
```

## Tips

- Run migrations before seeding: `make migrate-up`
- Use `-clear` flag to reset data between test runs
- Adjust numbers based on your testing needs
- All test users use password `password123` for easy login testing

## Troubleshooting

**Connection Error**: Make sure your database is running and migrations are up to date:
```bash
make docker-up
make migrate-up
```

**Duplicate Key Errors**: Use `-clear` flag to remove existing data first:
```bash
make seed-clear
```
