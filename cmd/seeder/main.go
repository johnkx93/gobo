package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/user/coc/internal/db"
)

const (
	defaultUsers  = 50
	defaultOrders = 200
)

func main() {
	// Parse command line flags
	numUsers := flag.Int("users", defaultUsers, "Number of users to generate")
	numOrders := flag.Int("orders", defaultOrders, "Number of orders to generate")
	clearData := flag.Bool("clear", false, "Clear existing data before seeding")
	flag.Parse()

	// Load environment variables
	_ = godotenv.Load()

	// Get database URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/appdb?sslmode=disable"
	}

	// Connect to database
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	queries := db.New(pool)

	// Clear data if requested
	if *clearData {
		fmt.Println("🗑️  Clearing existing data...")
		if err := clearDatabase(ctx, pool); err != nil {
			log.Fatalf("Failed to clear database: %v", err)
		}
		fmt.Println("✅ Data cleared")
	}

	// Seed users
	fmt.Printf("👥 Generating %d users...\n", *numUsers)
	userIDs, err := seedUsers(ctx, queries, *numUsers)
	if err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}
	fmt.Printf("✅ Created %d users\n", len(userIDs))

	// Seed orders
	fmt.Printf("📦 Generating %d orders...\n", *numOrders)
	orderCount, err := seedOrders(ctx, queries, userIDs, *numOrders)
	if err != nil {
		log.Fatalf("Failed to seed orders: %v", err)
	}
	fmt.Printf("✅ Created %d orders\n", orderCount)

	fmt.Println("\n🎉 Seeding completed successfully!")
}

func clearDatabase(ctx context.Context, pool *pgxpool.Pool) error {
	// Delete orders first (due to foreign key constraint)
	if _, err := pool.Exec(ctx, "DELETE FROM orders"); err != nil {
		return fmt.Errorf("failed to delete orders: %w", err)
	}

	// Delete users
	if _, err := pool.Exec(ctx, "DELETE FROM users"); err != nil {
		return fmt.Errorf("failed to delete users: %w", err)
	}

	return nil
}

func seedUsers(ctx context.Context, queries *db.Queries, count int) ([]pgtype.UUID, error) {
	userIDs := make([]pgtype.UUID, 0, count)

	// Hash a default password for all test users
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	for i := 0; i < count; i++ {
		person := gofakeit.Person()

		// Generate unique email and username
		email := gofakeit.Email()
		username := fmt.Sprintf("%s%d", gofakeit.Username(), i)

		firstName := pgtype.Text{String: person.FirstName, Valid: true}
		lastName := pgtype.Text{String: person.LastName, Valid: true}

		user, err := queries.CreateUser(ctx, db.CreateUserParams{
			Email:        email,
			Username:     username,
			PasswordHash: string(passwordHash),
			FirstName:    firstName,
			LastName:     lastName,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create user %d: %w", i, err)
		}

		userIDs = append(userIDs, user.ID)

		// Progress indicator
		if (i+1)%10 == 0 {
			fmt.Printf("  ... %d/%d users created\n", i+1, count)
		}
	}

	return userIDs, nil
}

func seedOrders(ctx context.Context, queries *db.Queries, userIDs []pgtype.UUID, count int) (int, error) {
	if len(userIDs) == 0 {
		return 0, fmt.Errorf("no users available to create orders for")
	}

	statuses := []string{"pending", "completed", "shipped", "cancelled"}
	orderCount := 0

	for i := 0; i < count; i++ {
		// Random user
		userID := userIDs[gofakeit.Number(0, len(userIDs)-1)]

		// Generate order number
		orderNumber := fmt.Sprintf("ORD-%s-%06d",
			time.Now().Format("20060102"), gofakeit.Number(100000, 999999))

		// Random status
		status := statuses[gofakeit.Number(0, len(statuses)-1)]

		// Random amount between 10 and 1000
		amount := fmt.Sprintf("%.2f", gofakeit.Price(10, 1000))
		totalAmount := pgtype.Numeric{}
		if err := totalAmount.Scan(amount); err != nil {
			return orderCount, fmt.Errorf("failed to scan amount: %w", err)
		}

		// Random notes (50% chance of having notes)
		notes := pgtype.Text{Valid: false}
		if gofakeit.Bool() {
			notes = pgtype.Text{
				String: gofakeit.Sentence(gofakeit.Number(5, 15)),
				Valid:  true,
			}
		}

		_, err := queries.CreateOrder(ctx, db.CreateOrderParams{
			UserID:      userID,
			OrderNumber: orderNumber,
			Status:      status,
			TotalAmount: totalAmount,
			Notes:       notes,
		})
		if err != nil {
			// Skip duplicates and continue
			if i < count-1 {
				continue
			}
			return orderCount, fmt.Errorf("failed to create order %d: %w", i, err)
		}

		orderCount++

		// Progress indicator
		if (orderCount)%50 == 0 {
			fmt.Printf("  ... %d/%d orders created\n", orderCount, count)
		}
	}

	return orderCount, nil
}
