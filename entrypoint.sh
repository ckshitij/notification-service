#!/bin/sh
set -e

echo "Waiting for MySQL..."
until nc -z mysql 3306; do
  sleep 2
done

echo "Running migrations..."
./migrator

echo "Starting notification service..."
exec ./notify-srv
