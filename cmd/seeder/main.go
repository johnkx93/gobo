package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/user/coc/internal/db"
)

const (
	defaultUsers     = 50
	defaultAdmins    = 5
	defaultAddresses = 150 // ~3 addresses per user on average
)

func main() {
	// Parse command line flags
	numUsers := flag.Int("users", defaultUsers, "Number of users to generate")
	numAdmins := flag.Int("admins", defaultAdmins, "Number of admins to generate")
	numAddresses := flag.Int("addresses", defaultAddresses, "Number of addresses to generate")
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

	// Seed admins
	fmt.Printf("👨‍💼 Generating %d admins...\n", *numAdmins)
	adminCount, err := seedAdmins(ctx, queries, *numAdmins)
	if err != nil {
		log.Fatalf("Failed to seed admins: %v", err)
	}
	fmt.Printf("✅ Created %d admins\n", adminCount)

	// Seed users
	fmt.Printf("👥 Generating %d users...\n", *numUsers)
	userIDs, err := seedUsers(ctx, queries, *numUsers)
	if err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}
	fmt.Printf("✅ Created %d users\n", len(userIDs))

	// Seed addresses
	fmt.Printf("📍 Generating %d addresses...\n", *numAddresses)
	addressCount, err := seedAddresses(ctx, queries, userIDs, *numAddresses)
	if err != nil {
		log.Fatalf("Failed to seed addresses: %v", err)
	}
	fmt.Printf("✅ Created %d addresses\n", addressCount)

	// Orders seeding removed

	fmt.Println("\n🎉 Seeding completed successfully!")
}

func clearDatabase(ctx context.Context, pool *pgxpool.Pool) error {
	// Delete addresses (must be before users due to foreign key)
	if _, err := pool.Exec(ctx, "DELETE FROM addresses"); err != nil {
		return fmt.Errorf("failed to delete addresses: %w", err)
	}

	// Delete users
	if _, err := pool.Exec(ctx, "DELETE FROM users"); err != nil {
		return fmt.Errorf("failed to delete users: %w", err)
	}

	// Delete admins (but keep the default super admin from migration)
	if _, err := pool.Exec(ctx, "DELETE FROM admins WHERE email != 'admin@example.com'"); err != nil {
		return fmt.Errorf("failed to delete admins: %w", err)
	}

	return nil
}

func seedAdmins(ctx context.Context, queries *db.Queries, count int) (int, error) {
	// Hash a default password for all test admins
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	adminCount := 0
	roles := []string{"admin", "super_admin", "moderator"}

	for i := 0; i < count; i++ {
		person := gofakeit.Person()

		// Generate unique email and username
		email := fmt.Sprintf("admin%d@%s", i, gofakeit.DomainName())
		username := fmt.Sprintf("admin_%s%d", gofakeit.Username(), i)

		firstName := pgtype.Text{String: person.FirstName, Valid: true}
		lastName := pgtype.Text{String: person.LastName, Valid: true}

		// Assign role (first one is super_admin, rest are random)
		role := roles[gofakeit.Number(0, len(roles)-1)]
		if i == 0 {
			role = "super_admin"
		}

		_, err := queries.CreateAdmin(ctx, db.CreateAdminParams{
			Email:        email,
			Username:     username,
			PasswordHash: string(passwordHash),
			FirstName:    firstName,
			LastName:     lastName,
			Role:         role,
			IsActive:     true,
		})
		if err != nil {
			// Skip duplicates and continue
			if i < count-1 {
				continue
			}
			return adminCount, fmt.Errorf("failed to create admin %d: %w", i, err)
		}

		adminCount++

		// Progress indicator
		if (adminCount)%5 == 0 {
			fmt.Printf("  ... %d/%d admins created\n", adminCount, count)
		}
	}

	return adminCount, nil
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

func seedAddresses(ctx context.Context, queries *db.Queries, userIDs []pgtype.UUID, count int) (int, error) {
	if len(userIDs) == 0 {
		return 0, fmt.Errorf("no users available to create addresses for")
	}

	addressCount := 0
	userAddressMap := make(map[int][]pgtype.UUID) // Track addresses per user index

	for i := 0; i < count; i++ {
		// Random user
		userIndex := gofakeit.Number(0, len(userIDs)-1)
		userID := userIDs[userIndex]

		// Generate address data
		address := gofakeit.Street()
		if len(address) > 50 {
			address = address[:50]
		}

		floor := fmt.Sprintf("%d", gofakeit.Number(1, 30))
		if len(floor) > 10 {
			floor = floor[:10]
		}

		unitNo := fmt.Sprintf("%02d", gofakeit.Number(1, 99))
		if len(unitNo) > 10 {
			unitNo = unitNo[:10]
		}

		// Block/Tower (optional - 60% chance)
		blockTower := pgtype.Text{Valid: false}
		if gofakeit.Bool() && gofakeit.Number(1, 100) <= 60 {
			block := fmt.Sprintf("Block %s", gofakeit.Letter())
			if len(block) > 25 {
				block = block[:25]
			}
			blockTower = pgtype.Text{String: block, Valid: true}
		}

		// Company name (optional - 30% chance for business address)
		companyName := pgtype.Text{Valid: false}
		if gofakeit.Number(1, 100) <= 30 {
			company := gofakeit.Company()
			if len(company) > 25 {
				company = company[:25]
			}
			companyName = pgtype.Text{String: company, Valid: true}
		}

		newAddress, err := queries.CreateAddress(ctx, db.CreateAddressParams{
			UserID:      userID,
			Address:     address,
			Floor:       floor,
			UnitNo:      unitNo,
			BlockTower:  blockTower,
			CompanyName: companyName,
		})
		if err != nil {
			return addressCount, fmt.Errorf("failed to create address %d: %w", i, err)
		}

		addressCount++

		// Track addresses per user for setting default later
		userAddressMap[userIndex] = append(userAddressMap[userIndex], newAddress.ID)

		// Progress indicator
		if (addressCount)%50 == 0 {
			fmt.Printf("  ... %d/%d addresses created\n", addressCount, count)
		}
	}

	// Set default address for users (randomly pick one of their addresses)
	fmt.Println("  🎯 Setting default addresses for users...")
	defaultCount := 0
	for userIndex, addressIDs := range userAddressMap {
		if len(addressIDs) > 0 {
			userID := userIDs[userIndex]

			// Pick random address as default
			defaultAddressID := addressIDs[gofakeit.Number(0, len(addressIDs)-1)]

			_, err := queries.SetDefaultAddress(ctx, db.SetDefaultAddressParams{
				ID:               userID,
				DefaultAddressID: defaultAddressID,
			})
			if err == nil {
				defaultCount++
			}
		}
	}
	fmt.Printf("  ✅ Set %d default addresses\n", defaultCount)

	return addressCount, nil
}
