#!/bin/bash

pkg install golang

# Clone the Koppel repository
echo "Cloning Koppel repository..."
pkg install git
git clone https://github.com/TiyoNotFound/Koppel.git

# Change directory to Koppel
cd Koppel

# Initialize a Go module
echo "Initializing Go module..."
go mod init main.go

# Fetch dependencies
echo "Fetching dependencies..."
go get .

# Build the koppel binary
echo "Building koppel binary..."
go build -o koppel

# Check if the build was successful
if [ $? -ne 0 ]; then
  echo "Failed to build koppel binary."
  exit 1
fi

# Move the koppel binary to usr/bin
echo "Moving koppel binary to usr/bin..."
mv koppel /data/data/com.termux/files/usr/bin

# Check if moving the binary was successful
if [ $? -ne 0 ]; then
  echo "Failed to move koppel binary to usr/bin."
  exit 1
fi

echo "Koppel installed successfully."
