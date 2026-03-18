# Apply all pending migrations
go run ./cmd/migrate up

# Roll back the last migration
go run ./cmd/migrate down

# Roll back 3 migrations
go run ./cmd/migrate down -steps 3

# See what's applied vs pending
go run ./cmd/migrate status