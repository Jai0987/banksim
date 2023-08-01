SHELL := /bin/bash

.PHONY: help
help:
	@echo "Usage make <TARGET>"
	@echo ""
	@echo "  Targets:"
	@echo "	   build				Build the main executable"
	@echo "	   run					Run the main program from code"
	@echo "	   run-bin				Run the main executable"
	@echo "	   build-user			Build the user executable"
	@echo "	   run-user			Run the user program from code"
	@echo "	   run-user-bin			Run the user executable"
	@echo "	   build-account		Build the account executable"
	@echo "	   run-account			Run the account program from code"
	@echo "	   run-account-bin		Run the account executable"
	@echo "	   load-db				Add sample data via curl"
	@echo "	   get-by-id			Get a user by id, pass id=<id> on command line"
	@echo "	   get-all				Get all users"
	@echo "	   update-2				Update user with id=2, pass a new title in using title=<title> on command line"
	@echo "	   delete-all			Delete all users"
	@echo "	   delete-by-id			Delete a user by id, pass id=<id> on command line"
	@echo "	   get-v2				Get all users by done status, pass done=<true|false> on command line"
	@echo "	   get-v2-all			Get all users using version 2"
	@echo "	   get-account			Get an account by id, pass id=<id> on command line"
	@echo "	   pay-bill				Mark a bill as paid, pass id=<id> on command line"
	@echo "	   get-due-date			Get the payment due date for an account, pass id=<id> on command line"
	@echo "	   get-credit-score		Get the credit score for a user, pass id=<id> on command line"
	@echo "	   build-amd64-linux	Build amd64/Linux executables"
	@echo "	   build-arm64-linux	Build arm64/Linux executables"

.PHONY: build
build:
	go build -o main .

.PHONY: run
run:
	go run .

.PHONY: run-bin
run-bin:
	./main

#To interact with users:
.PHONY: get-by-id
get-by-id: run-user-bin
	./user getUser $(id)

.PHONY: get-all
get-all: run-user-bin
	./user getUser

.PHONY: update-2
update-2: run-user-bin
	./user updateUser $(title)

.PHONY: delete-by-id
delete-by-id: run-user-bin
	./user deleteUser $(id)

# For example, to interact with accounts:
.PHONY: get-account
get-account: run-account-bin
	./account getAccount $(id)

.PHONY: pay-bill
pay-bill: run-account-bin
	./account payBill $(id)

.PHONY: get-due-date
get-due-date: run-account-bin
	./account getDueDate $(id)

.PHONY: get-credit-score
get-credit-score: run-account-bin
	./account getCreditScore $(id)

# Add other targets as needed, specific to user.go and account.go files
