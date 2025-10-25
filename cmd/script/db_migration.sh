#!/bin/bash

migrate -path=database/migrations -database "postgres://postgres:admin123@localhost:5432/ecomm?sslmode=disable" -verbose up
