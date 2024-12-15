#!/bin/bash

# Set the application directory
APP_DIR="/Users/ngocduongminh/Projects/genomic-system/genomic-service"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Navigate to the application directory
cd "$APP_DIR" || { echo -e "${RED}Error: Directory $APP_DIR does not exist.${NC}"; exit 1; }

# Start the application in the background
echo -e "${BOLD}Starting the application...${NC}"
go run main.go &
APP_PID=$!

# Ensure the application starts correctly
sleep 5
if ! ps -p $APP_PID > /dev/null; then
  echo -e "${RED}Error: Application failed to start.${NC}"
  exit 1
fi

echo -e "${GREEN}Application started with PID $APP_PID.${NC}"

USER_ADDRESS=${1:-"0x1234567890abcdef1234567890abcdef12345678"}
# Define the user data to be sent in the POST request
USER_DATA=$(jq -n --arg address "$USER_ADDRESS" '{ "address": $address }')

# Step 1: Register user
echo -e "${BOLD}Step 1 - Register user${NC}"
echo "Sending POST request to /v1/users..."
HTTP_RESPONSE=$(curl -s -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d "$USER_DATA")

# Extract the public key
PUBLIC_KEY=$(echo "$HTTP_RESPONSE" | jq -r '.data.public_key')
echo -e "HTTP Response Body:\n$(echo "$HTTP_RESPONSE" | jq '.')"

# Step 2: Upload Data
RAW_DATA=$(jq -n --arg pubKey "$PUBLIC_KEY" --arg data "low risk" \
  '{ "public_key": $pubKey, "data": $data }')

echo -e "${BOLD}Step 2 - Upload Data${NC}"
echo "Sending POST request to /v1/tee/encrypt to encrypt raw data"
HTTP_RESPONSE=$(curl -s -X POST http://localhost:8080/v1/tee/encrypt \
  -H "Content-Type: application/json" \
  -d "$RAW_DATA")

ENCRYPTED_DATA=$(echo "$HTTP_RESPONSE" | jq -r '.data.encrypted_data')
echo -e "HTTP Response Body:\n$(echo "$HTTP_RESPONSE" | jq '.')"

UPLOAD_DATA=$(jq -n --arg pubKey "$USER_ADDRESS" --arg encrypted_data "$ENCRYPTED_DATA" \
  '{ "address": $pubKey, "encrypted_data": $encrypted_data, "signature": "signature" }')

echo "Sending POST request to /v1/users/upload to upload data"
HTTP_RESPONSE=$(curl -s -X POST http://localhost:8080/v1/users/upload \
  -H "Content-Type: application/json" \
  -d "$UPLOAD_DATA")

FILE_ID=$(echo "$HTTP_RESPONSE" | jq -r '.data.file_id')
echo -e "HTTP Response Body:\n$(echo "$HTTP_RESPONSE" | jq '.')"

# Pause to simulate waiting for confirmation
sleep 10

# Step 3: Confirm and reward
echo -e "${BOLD}Step 3 - Confirm and reward for user${NC}"
echo "Fetching uploaded data to confirm the data"
HTTP_RESPONSE=$(curl -s -X GET "http://localhost:8080/v1/genetic-data/${FILE_ID}" \
  -H "Content-Type: application/json")

echo -e "HTTP Response Body:\n$(echo "$HTTP_RESPONSE" | jq '.')"

# Stop the application
echo -e "${BOLD}Stopping the application...${NC}"
if lsof -ti :8080 > /dev/null; then
    echo "Killing process on port 8080..."
    kill -9 $(lsof -ti :8080)
fi

echo -e "${GREEN}Application stopped.${NC}"
exit 0