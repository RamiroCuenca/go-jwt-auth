# This rule creates the container with docker using the postgres image
docker-container-create:
	docker run --name go-auth-container -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres

# This rule runs the postgres container (If it is stopped)
docker-container-start:
	docker container start go-auth-container

# This rule stops the posgres container
docker-container-stop:
	docker container stop go-auth-container

# This rule creates the db on the container
docker-postgres-createdb:
	docker container exec -it go-auth-container createdb --username=postgres --owner=postgres go-auth-database

# This rule deletes the db on the container
docker-postgres-dropdb:
	docker exec -it go-auth-container dropdb bank_app

# This rule runs the migrations up
run-migrations-up:
	migrate --path database/migrations --database "postgresql://postgres:postgres@localhost:5432/go-auth-database?sslmode=disable" --verbose up

run-migrations-down:
	migrate --path database/migrations --database "postgresql://postgres:postgres@localhost:5432/go-auth-database?sslmode=disable" --verbose down

# .PHONY tell explicitly to MAKE that those rules are not associated with files
.PHONY: docker-container-create docker-container-start docker-container-stop docker-postgres-createdb docker-postgres-dropdb run-migrations-up run-migrations-down
