# Database Migrations

This directory contains SQL migration files for all GORM models using [golang-migrate](https://github.com/golang-migrate/migrate) format for PostgreSQL.

## Migration Files

Each migration has two files:
- `{version}_{name}.up.sql` - Migration to apply
- `{version}_{name}.down.sql` - Migration to rollback

### Migration Order

1. **000001_create_role_table** - Creates the `role` table
2. **000002_create_menu_table** - Creates the `menu` table
3. **000003_create_menu_mapping_table** - Creates the `menu_mapping` table (requires role and menu tables)
4. **000004_create_users_table** - Creates the `users` table (requires role table)
5. **000005_create_face_datasets_table** - Creates the `face_datasets` table (requires users table)
6. **000006_create_model_training_table** - Creates the `model_training` table
7. **000007_create_parameter_table** - Creates the `parameter` table

## Installation

### Install golang-migrate CLI

**macOS:**
```bash
brew install golang-migrate
```

**Linux:**
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/migrate
```

**Windows:**
```powershell
choco install migrate
```

Or download from: https://github.com/golang-migrate/migrate/releases

## Usage

### Set Database URL

Set your PostgreSQL connection string as an environment variable:

```bash
export DATABASE_URL="postgres://username:password@localhost:5432/database_name?sslmode=disable"
```

Or use the full connection format:

```bash
export DATABASE_URL="postgres://username:password@host:port/database?sslmode=disable"
```

### Run Migrations

**Apply all pending migrations:**
```bash
migrate -path migrations -database "$DATABASE_URL" up
```

**Apply specific number of migrations:**
```bash
migrate -path migrations -database "$DATABASE_URL" up 2
```

**Rollback last migration:**
```bash
migrate -path migrations -database "$DATABASE_URL" down 1
```

**Rollback all migrations:**
```bash
migrate -path migrations -database "$DATABASE_URL" down
```

**Check migration version:**
```bash
migrate -path migrations -database "$DATABASE_URL" version
```

**Force migration version (use with caution):**
```bash
migrate -path migrations -database "$DATABASE_URL" force {version}
```

### Create New Migration

```bash
migrate create -ext sql -dir migrations -seq create_new_table
```

This will create:
- `{next_number}_create_new_table.up.sql`
- `{next_number}_create_new_table.down.sql`

## Using in Go Code

### Install golang-migrate library

```bash
go get -u github.com/golang-migrate/migrate/v4
go get -u github.com/golang-migrate/migrate/v4/database/postgres
go get -u github.com/golang-migrate/migrate/v4/source/file
```

### Example Usage

```go
package main

import (
    "database/sql"
    "log"
    
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/lib/pq"
)

func main() {
    db, err := sql.Open("postgres", "postgres://user:pass@localhost/dbname?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        log.Fatal(err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations",
        "postgres", driver)
    if err != nil {
        log.Fatal(err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal(err)
    }
}
```

## Important Notes

- Migrations are run in numerical order (000001, 000002, etc.)
- Foreign key constraints are defined, so tables must be created in the correct order
- The migrations use `CREATE TABLE IF NOT EXISTS` to prevent errors if tables already exist
- PostgreSQL doesn't support `ON UPDATE CURRENT_TIMESTAMP` like MySQL. A trigger function (migration 000008) automatically updates the `updated_at` field on UPDATE for `role`, `menu`, and `menu_mapping` tables
- All tables use PostgreSQL's default character encoding (UTF-8)

## Database Connection Update Required

⚠️ **Important**: Your current `app/connection/connection.go` file is configured for MySQL. Since you're using PostgreSQL, you'll need to update it to use the PostgreSQL driver:

1. Update `go.mod` to include PostgreSQL driver:
   ```bash
   go get gorm.io/driver/postgres
   ```

2. Update `app/connection/connection.go`:
   ```go
   import "gorm.io/driver/postgres"
   
   func NewDatabaseConnection(c *config.Database) *gorm.DB {
       dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
           c.Host,
           c.Username,
           c.Password,
           c.Database,
           c.Port)
       
       db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
       // ... rest of the code
   }
   ```

## Table Relationships

- `users` → references `role` (via `role_id`)
- `face_datasets` → references `users` (via `username`)
- `menu_mapping` → references `menu` (via `menu_id`) and `role` (via `role_id`)
