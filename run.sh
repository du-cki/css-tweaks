#!/bin/bash

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

if ! command_exists git; then
    echo "→ error: 'git' is not installed."
    exit 1
fi

if ! command_exists docker; then
    echo "→ error: 'docker' is not installed."
    exit 1
fi

if docker compose version >/dev/null 2>&1; then
    COMPOSE_CMD="docker compose"
    echo "→ using: docker compose"
elif command_exists docker-compose; then
    COMPOSE_CMD="docker-compose"
    echo "→ using: docker-compose"
else
    echo "→ error: Neither 'docker compose' nor 'docker-compose' found."
    echo "→ please install Docker Compose."
    exit 1
fi

if git rev-parse --git-dir > /dev/null 2>&1; then
    export GIT_COMMIT=$(git rev-parse --short HEAD)
    export GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
    export GIT_ORIGIN=$(git remote get-url origin 2>/dev/null || echo "unknown")
else
    echo "→ warning: Not a git repository. Using default values."

    export GIT_COMMIT="dev"
    export GIT_BRANCH="unknown"
    export GIT_ORIGIN="unknown"
fi

$COMPOSE_CMD up -d --build