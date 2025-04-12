package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

// Version information
var (
	buildTime = "2025-04-11 07:34:25"
	buildUser = "Anton-Bondarchuk"
	version   = "1.0.0"
)

func main() {
	// Command line arguments
	var (
		dbConnection   string
		migrationsPath string
		migrationsTable string
		verbose        bool
		showVersion    bool
	)

	flag.StringVar(&dbConnection, "db", os.Getenv("DATABASE_URL"), "PostgreSQL connection string")
	flag.StringVar(&migrationsPath, "migrations-path", "./db/migrations", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "schema_migrations", "name of migrations table")
	flag.BoolVar(&verbose, "verbose", false, "show verbose output")
	flag.BoolVar(&showVersion, "version", false, "show version information")
	flag.Parse()

	// Show version if requested
	if showVersion {
		fmt.Printf("Migrator v%s\n", version)
		fmt.Printf("Build: %s by %s\n", buildTime, buildUser)
		return
	}

	// Validate required params
	if dbConnection == "" {
		log.Fatal("Database connection string is required. Provide it with -db flag or DATABASE_URL environment variable")
	}
	if migrationsPath == "" {
		log.Fatal("Migrations path is required")
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to the database
	conn, err := pgx.Connect(ctx, dbConnection)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	// Create the migrations table if it doesn't exist
	if err := ensureMigrationsTable(ctx, conn, migrationsTable); err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	// Get all migration files
	migrationFiles, err := getMigrationFiles(migrationsPath)
	if err != nil {
		log.Fatalf("Failed to read migration files: %v", err)
	}

	if len(migrationFiles) == 0 {
		fmt.Println("No migration files found")
		return
	}

	// Get applied migrations
	appliedMigrations, err := getAppliedMigrations(ctx, conn, migrationsTable)
	if err != nil {
		log.Fatalf("Failed to get applied migrations: %v", err)
	}

	// Apply pending migrations
	count := 0
	for _, file := range migrationFiles {
		// Skip if already applied
		if containsMigration(appliedMigrations, file.Name) {
			if verbose {
				fmt.Printf("Skipping already applied migration: %s\n", file.Name)
			}
			continue
		}

		// Apply migration
		if verbose {
			fmt.Printf("Applying migration: %s\n", file.Name)
		}

		// Read migration file
		content, err := os.ReadFile(filepath.Join(migrationsPath, file.Name))
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", file.Name, err)
		}

		// Begin transaction
		tx, err := conn.Begin(ctx)
		if err != nil {
			log.Fatalf("Failed to begin transaction: %v", err)
		}

		// Execute migration
		if _, err := tx.Exec(ctx, string(content)); err != nil {
			tx.Rollback(ctx)
			log.Fatalf("Failed to execute migration %s: %v", file.Name, err)
		}

		// Record migration
		if _, err := tx.Exec(ctx, 
			fmt.Sprintf("INSERT INTO %s (version, applied_at) VALUES ($1, NOW())", migrationsTable),
			file.Name); err != nil {
			tx.Rollback(ctx)
			log.Fatalf("Failed to record migration %s: %v", file.Name, err)
		}

		// Commit transaction
		if err := tx.Commit(ctx); err != nil {
			log.Fatalf("Failed to commit transaction: %v", err)
		}

		count++
		fmt.Printf("Applied migration: %s\n", file.Name)
	}

	if count == 0 {
		fmt.Println("No migrations to apply")
	} else {
		fmt.Printf("Successfully applied %d migrations\n", count)
	}
}

// Migration represents a migration file
type Migration struct {
	Name string
	Path string
}

// ensureMigrationsTable creates the migrations table if it doesn't exist
func ensureMigrationsTable(ctx context.Context, conn *pgx.Conn, tableName string) error {
	_, err := conn.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE NOT NULL
		)
	`, tableName))
	return err
}

// getMigrationFiles returns a sorted list of migration files
func getMigrationFiles(path string) ([]Migration, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []Migration
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		name := entry.Name()
		// Only consider .sql files
		if !strings.HasSuffix(strings.ToLower(name), ".sql") {
			continue
		}
		
		files = append(files, Migration{
			Name: name,
			Path: filepath.Join(path, name),
		})
	}

	// Sort files by name to ensure they're applied in the right order
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	return files, nil
}

// getAppliedMigrations returns a list of already applied migrations
func getAppliedMigrations(ctx context.Context, conn *pgx.Conn, tableName string) ([]string, error) {
	rows, err := conn.Query(ctx, fmt.Sprintf("SELECT version FROM %s ORDER BY version", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []string
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	return versions, rows.Err()
}

// containsMigration checks if a migration has already been applied
func containsMigration(migrations []string, name string) bool {
	for _, m := range migrations {
		if m == name {
			return true
		}
	}
	return false
}