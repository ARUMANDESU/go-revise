#!/bin/bash

# Check if port is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <port>"
    exit 1
fi

PORT=$1

# Get the script's directory and project root directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Check if ngrok is installed
if ! command -v ngrok &> /dev/null; then
    echo "ngrok is not installed. Please install it first."
    exit 1
fi

# Start ngrok in the background and capture the URL
echo "Starting ngrok on port $PORT..."
ngrok http $PORT > /dev/null 2>&1 &

# Wait for ngrok to start
sleep 3

# Get the public URL from ngrok API
NGROK_URL=$(curl -s localhost:4040/api/tunnels | grep -o '"public_url":"https:\/\/[^"]*' | cut -d'"' -f4)

if [ -z "$NGROK_URL" ]; then
    echo "Failed to get ngrok URL"
    exit 1
fi

echo "ngrok URL: $NGROK_URL"

# Update .env file in project root directory
ENV_FILE="$PROJECT_ROOT/.env"

if [ -f "$ENV_FILE" ]; then
    # Check if TELEGRAM_WEBHOOK_URL exists in .env
    if grep -q "TELEGRAM_WEBHOOK_URL=" "$ENV_FILE"; then
        # Replace existing TELEGRAM_WEBHOOK_URL
        sed -i.bak "s|TELEGRAM_WEBHOOK_URL=.*|TELEGRAM_WEBHOOK_URL=$NGROK_URL|" "$ENV_FILE"
    else
        # Add new TELEGRAM_WEBHOOK_URL
        echo "TELEGRAM_WEBHOOK_URL=$NGROK_URL" >> "$ENV_FILE"
    fi
    echo "Updated TELEGRAM_WEBHOOK_URL in .env file"
else
    echo "TELEGRAM_WEBHOOK_URL=$NGROK_URL" > "$ENV_FILE"
    echo "Created new .env file with TELEGRAM_WEBHOOK_URL"
fi

echo "Setup complete! Your webhook URL is: $NGROK_URL"

# Print process info for cleanup
NGROK_PID=$(pgrep ngrok)
echo "ngrok process ID: $NGROK_PID"
echo "To stop ngrok, run: kill $NGROK_PID"
