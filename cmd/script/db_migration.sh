#!/bin/bash

migrate -path ./database/migrations -database "postgres://postgres:admin123@localhost:5432/ecomm?sslmode=disable" down -all

go run cmd/script/migration.go