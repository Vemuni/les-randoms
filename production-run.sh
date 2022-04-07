#!/bin/sh

. ./setup-secret-keys.sh

export DEBUG_MODE="FALSE"
export WEBSITE_URL="http://vemuni.ga"
export PORT="8080"

sudo iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-port $PORT

echo "Ctrl+C to stop running the local server"
echo "See on $WEBSITE_URL:$PORT"
mkdir logs
go build -o bin/les-randoms -v . && ./bin/les-randoms 2>logs/_lastcrash.log &

# To start postgresl console :
# sudo -u postgres psql

# Deprecated :
# use 'winpty sqlite3 sqlite-database.db' to run sqlite3 command on git bash
