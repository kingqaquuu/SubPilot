#!/bin/sh
set -eu

if [ "${APP_ENV:-development}" = "production" ]; then
  if [ -z "${POSTGRES_PASSWORD:-}" ] || [ "${POSTGRES_PASSWORD}" = "change-me" ]; then
    echo "POSTGRES_PASSWORD must be set to a non-default value in production" >&2
    exit 1
  fi

  if [ -z "${JWT_SECRET:-}" ] || [ "${JWT_SECRET}" = "change-me-to-a-long-random-secret" ]; then
    echo "JWT_SECRET must be set to a non-default value in production" >&2
    exit 1
  fi
fi

exec "$@"
