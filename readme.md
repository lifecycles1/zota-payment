# Zota Payment Integration Project

## Overview

This project integrates with the Zota payment gateway, providing a seamless payment experience through client-side flow with redirects. The architecture consists of a React/Vite frontend, a Golang backend, and ngrok for secure tunneling. The Go API serves as a proxy to the Zota API, handling deposit requests and order status queries. The frontend interacts with this Go API and provides a user interface for initiating and managing payments.

## API Endpoints

### Flowless Endpoints

1. **Deposit Request**

   - **URL**: `/api/v1/deposit/request/{endpointID}/`
   - **Method**: POST
   - **Description**: Creates a deposit request and returns the response from Zota.
   - **Body Parameters**: JSON payload required by Zota.

2. **Order Status Request**
   - **URL**: `/api/v1/query/order-status/`
   - **Method**: GET
   - **Description**: Checks the order status. The user needs to retry until a final status is obtained.
   - **Query Parameters**:
     - `merchantID`
     - `orderID`
     - `merchantOrderID`
     - `timestamp`
     - `signature`

### Flow Endpoints

3. **Client Flow with Redirects**

   - **URL**: `/api/v1/deposit/client-flow/{endpointID}/`
   - **Method**: POST
   - **Description**: Initiates the deposit process, redirects the user to Zota's deposit page, and handles the final status via client-side polling.

## Setup Instructions

### Prerequisites

1. **ngrok**: Download and install ngrok from the [official website](https://ngrok.com/).
2. **Node.js**: Ensure Node.js is installed for running the frontend.
3. **Go**: Ensure Go is installed for running the backend.

### Installation and Startup Steps

1. **Install ngrok and Set Up Authentication**

   - Install ngrok and add the folder path to your system's PATH.
   - Register at the ngrok website to obtain an auth token.
   - Open `ngrok/ngrok.yml` and update the `authtoken` with your ngrok auth token.

   ```yaml
   authtoken: your-ngrok-auth-token
   ```

2. **Setup and Run**

   - Navigate to the `client` directory and install dependencies:

   ```bash
   cd client
   npm install
   ```

   - Start ngrok from the `client` directory:

   ```bash
   npm run start:ngrok
   ```

   - Run the frontend: open another terminal, navigate to the `client` directory, and start the development server:

   ```bash
   cd client
   npm run dev
   ```

   - Run the backend: open a third terminal, navigate to the `server` directory, and run the Go application:

   ```bash
   cd server
   go run main.go
   ```

   - Launch the application:
     - Open your browser and go to http://localhost:5173 to interact with the frontend.
     - Test the deposit flow or the single endpoints through the frontend interface.
     - All endpoints except 'Client Flow' can also be tested with tools like Postman, Insomnia, etc.

## Testing Endpoints

### Test through frontend client, Postman or similar tools

1. **Deposit Request**

   - **URL**: `http://localhost:8080/api/v1/deposit/request/{endpointID}/`
   - **Method**: POST
   - **Body**: JSON payload required by Zota.

2. **Order Status Request**
   - **URL**: `http://localhost:8080/api/v1/query/order-status/`
   - **Method**: GET
   - **Query Params**:
     - `merchantID`
     - `merchantOrderID`
     - `orderID`
     - `timestamp`
     - `signature`

## Testing Flow

### Test `Client Flow` through frontend client only.

1. **Testing the Client Flow**

   - Trigger the deposit request from the frontend.
   - Observe the frontend handling the redirects to Zota's deposit page.
   - Complete the payment process.
   - Zota redirects back to frontend's PaymentReturn redirect url which displays the order status results and keeps polling for a final status.

## Additional Tests

### Run Tests on Frontend with Vitest/React Testing Library

1. **Unit/Integration/End-to-end**

   - Navigate to `client` directory, run the test command, then observe the results:

   ```bash
   cd client
   npm run test
   ```

### Run Tests on Backend with Built-in Testing library and Testify

2. **Integration**

   - Navigate to `server` directory, run the test command to recursively find and run all test files in the backend project. Include verbose flag:

   ```bash
   cd server
   go test -v ./...
   ```
