#!/bin/sh

echo "running migrations"
# create if not there already
touch /opt/canary/data/measurements.db
# run idempotent migrations
sqlite3 /opt/canary/data/measurements.db '.read /opt/canary/migrations.sql'
echo "done running migrations."

echo "starting server"
/opt/canary/canary