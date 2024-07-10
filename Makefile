migration_up: 
	@echo "--migrate up"
	migrate -path "migrations/" -database "postgresql://postgres:postgres@localhost/recipesite?sslmode=disable" -verbose up

migration_down: 
	@echo "--migrate down"
	migrate -path "migrations/" -database "postgresql://postgres:postgres@localhost/recipesite?sslmode=disable" -verbose down

migration_fix: 
	@echo "--fix migration"
	migrate -path "migrations/" -database "postgresql://postgres:postgres@localhost/recipesite?sslmode=disable" force 1