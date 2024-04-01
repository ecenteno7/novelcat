#!/bin/bash
# Load .env file
export $(egrep -v '^#' .env | xargs)

# Variables
REMOTE_USER=root
REMOTE_DIR=/pb/pb_data
LOCAL_DIR=./pb_data

# Sync remote directory to local directory
echo "Syncing remote directory to local..."
rsync -avzc --delete --exclude='.autocert_cache' -e ssh $REMOTE_USER@$REMOTE_IP:$REMOTE_DIR/ $LOCAL_DIR/
if [ $? -eq 0 ]; then
    echo "Sync successful."
else
    echo "Sync failed."
fi