#!/bin/bash

# Define the echo server hostname or IP address and port
echo_server="server"
echo_port=12345

# Send the string "hello" to the echo server using netcat and capture the response
response=$(echo "hello" | nc $echo_server $echo_port)

# Check if the response matches the expected output "hello"
if [ "$response" = "hello" ]; then
  # If the response matches, print a success message
  echo "Echo test successful. Response: $response. Server is working"
else
  # If the response does not match, print an error message
  echo "Echo test failed. Expected 'hello', got '$response'"
fi
