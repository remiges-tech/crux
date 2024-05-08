# Remiges Crux

This is the home of Remiges Crux, the business rules engine and workflow engine from Remiges Technologies Pvt Ltd. It is available under the Apache 2.0 licence.

# CRUX Installation Guide

## Prerequisites

Before you begin the installation process, ensure that you have the following prerequisites installed on your system:

- [Go](https://golang.org/) installed on your system
- [Tern](https://github.com/jackc/tern) for database migrations
- [Make](https://www.gnu.org/software/make/) for running Makefile commands

## Installation Steps

Follow these steps to install Crux on your system:

### Step 1: Clone the Repository

Clone the Crux  repository from GitHub:

```bash
git clone https://github.com/remiges-tech/crux.git
```

### Step 2: Navigate to the  Directory

```bash
cd crux
```

### Step 3: synchronize  Dependencies

```bash
go mod tidy
```

### Step 4: Run Database Migrations

Before running database migrations, ensure that your database-related details such as username, password, dbname, and port are correctly configured in the db/migration/tern.config file.

```bash
make db-migrate-generate
```

### Step 5: Configure and Run 

Before configuring and running the, make sure to add your configuration values in the `config.json` file located in the  directory.

```bash
go run main.go
```

For more information, head over to the [wiki](https://github.com/remiges-tech/crux/wiki)
