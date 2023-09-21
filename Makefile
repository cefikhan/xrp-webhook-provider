run:
	go run .

make migrate-up:
	migrate -path db/migration/ -database "postgresql://postgres:pass123@localhost:4848/xrp_db?sslmode=disable" -verbose up

make migrate-down:
	migrate -path db/migration/ -database "postgresql://postgres:pass123@localhost:4848/xrp_db?sslmode=disable" -verbose down


make migration-create:
	migrate create -ext sql -dir db/migration/ -seq $(migration_name)


# make migration-create migration_name=create_market_maker