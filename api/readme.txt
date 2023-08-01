# Go Bank Management API

This is a Go API project that provides a bank management system. It uses the Gin framework for building the API endpoints and interacts with a PostgreSQL database for data storage.

## Prerequisites

Before running this API, make sure you have the following requirements installed:

- Go (1.16 or higher)
- PostgreSQL database

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/Jai0987/ginbankapi.git
   cd ginbankapi
   ```

2. Install the project dependencies:

   ```bash
   go mod download
   ```

3. Set up the PostgreSQL database:
   - Create a new PostgreSQL database for the project.
   - Update the database connection string in `main.go` to match your database details.

4. Set up Google OAuth credentials:
   - Create a new project in the [Google Cloud Console](https://console.cloud.google.com/).
   - Enable the "Google+ API" for the project.
   - Configure the OAuth consent screen with the necessary details.
   - Create credentials (OAuth client ID) for the project and obtain the client ID and client secret.
   - Set the `CLIENT_ID` and `CLIENT_SECRET` environment variables with your credentials.

5. Run the API:

   ```bash
   go run main.go
   ```

   The API will be accessible at `http://localhost:8080`.

## API Endpoints

The API provides the following endpoints:

### User Details

- `POST /user` - Create a new user.
- `GET /user/:id` - Get user details by ID.
- `PUT /user/:id` - Update user details by ID.
- `DELETE /user/:id` - Delete a user by ID.

### Account Operations

- `GET /account/:id` - Get account details by ID.
- `POST /account/:id/pay` - Mark a bill as paid for the account.
- `GET /account/:id/due` - Get payment due date for the account.
- `GET /account/:id/score` - Get credit score for the account.

## Contribution

Contributions to this project are welcome. Feel free to open issues and submit pull requests to improve the code, add new features, or fix bugs.

## Contact

For any questions or inquiries, please contact [jainamkashyap1204@gmail.com](mailto:jainamkashyap1204@gmail.com).
