
# Trading App GUI

A web-based trading application that allows users to execute trades through an intuitive dashboard interface.

## Features

- **User Registration and Authentication**: Secure user sign-up and login functionalities.
- **Dashboard Interface**: Interactive web-based dashboard for managing trades.
- **Order Placement**: Ability to place buy and sell orders directly from the dashboard.
- **Order Matching Engine**: Automated matching of buy and sell orders based on predefined criteria.
- **Trade History**: View detailed history of executed trades.

## Requirements

- Go 1.20 or later
- CockroachDB installed and running
- Environment variables configured in a `.env` file

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   ```
2. Navigate to the project directory:
   ```bash
   cd trading_app_gui
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Start CockroachDB:
   ```bash
   cockroach start-single-node --insecure --listen-addr=localhost:26257
   ```
5. Create the database:
   ```bash
   cockroach sql --insecure --host=localhost -e "CREATE DATABASE trading;"
   ```
6. Run the application:
   ```bash
   go run main.go
   ```

<!-- For a detailed explanation in Japanese, please refer to the following article:
[CLI通貨取引アプリ作ってみた](https://zenn.dev/genn_tmm/articles/70a2d9025ebdb4) -->
