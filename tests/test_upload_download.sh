#!/bin/bash
set -e

SERVER_URL="http://localhost:8080"
USERNAME="testuser"
PASSWORD="testpass"
FILE_ORIG="one.txt"
FILE_DOWNLOADED="downloads/one.txt"

echo "This is a test file for upload and download." > $FILE_ORIG

./pkg/client/zerodupe-client signup --server $SERVER_URL --username $USERNAME --password $PASSWORD --confirm-password $PASSWORD || true

TOKENS=$(./pkg/client/zerodupe-client login --server $SERVER_URL --username $USERNAME --password $PASSWORD 2>&1)
ACCESS_TOKEN=$(echo "$TOKENS" | grep "Access token:" | awk '{print $3}')
REFRESH_TOKEN=$(echo "$TOKENS" | grep "Refresh token:" | awk '{print $3}')

UPLOAD_OUTPUT=$(./pkg/client/zerodupe-client upload --server $SERVER_URL --token $ACCESS_TOKEN $FILE_ORIG 2>&1)
FILE_HASH=$(echo "$UPLOAD_OUTPUT" | grep "File hash:" | awk '{print $3}')

mkdir -p downloads
./pkg/client/zerodupe-client download --server $SERVER_URL --token $ACCESS_TOKEN -o downloads -n one.txt $FILE_HASH

if cmp -s "$FILE_ORIG" "$FILE_DOWNLOADED"; then
    echo "SUCCESS: Downloaded file matches original."
else
    echo "ERROR: Downloaded file does NOT match original."
    exit 1
fi