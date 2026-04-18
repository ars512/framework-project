DB_URL="postgresql://postgres:1234@localhost:5432/shop?sslmode=disable"

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	migrate -path migrations -database $(DB_URL) -verbose up

migrate-down:
	migrate -path migrations -database $(DB_URL) -verbose down

migrate-down-1:
	migrate -path migrations -database $(DB_URL) -verbose down 1

migrate-force:
	migrate -path migrations -database $(DB_URL) force $(v)