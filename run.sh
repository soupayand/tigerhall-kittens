#!/bin/bash


# Function to handle cleanup (kill the background processes) on SIGINT
function cleanup {
    echo "Stopping the server..."
    kill -15 $main_pid
    echo "Server stopped."
    exit 0
}

# Trap the SIGINT signal and call the cleanup function
trap cleanup SIGINT
> error.log
> info.log
# Run the main.go and worker_main.go in the background and capture their PIDs
go run main.go & main_pid=$!

# Echo the message after 5 seconds
echo "Server is up and running"

# Echo the PIDs for later use (e.g., stopping the processes)
echo "Main PID: $main_pid"
echo "Press control + c  to stop the server"

# Wait for both background processes to finish
wait $main_pid

# The cleanup function will be automatically called on SIGINT (Ctrl+C)

