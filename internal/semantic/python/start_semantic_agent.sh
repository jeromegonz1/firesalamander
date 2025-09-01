#!/bin/bash
#
# Start the Semantic Analysis Agent server
#

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Default configuration
PORT=${1:-8003}
PYTHON_ENV="venv"

echo "=== Fire Salamander - Semantic Analysis Agent ==="
echo "Port: $PORT"
echo "Working directory: $SCRIPT_DIR"

# Check if virtual environment exists
if [ ! -d "$PYTHON_ENV" ]; then
    echo "ERROR: Virtual environment not found at $PYTHON_ENV"
    echo "Please run: python3 -m venv $PYTHON_ENV && source $PYTHON_ENV/bin/activate && pip install -r requirements.txt"
    exit 1
fi

# Activate virtual environment
source "$PYTHON_ENV/bin/activate"

# Check if dependencies are installed
if ! python -c "import flask, yaml, numpy, sklearn" 2>/dev/null; then
    echo "Installing dependencies..."
    pip install -r requirements.txt
fi

# Start the server
echo "Starting semantic analysis server on port $PORT..."
python semantic_server.py $PORT