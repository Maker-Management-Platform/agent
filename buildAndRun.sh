#!/bin/bash
set -e

# Build mode: v1 | v2 | both
MODE=${1:-v2}

echo "================================"
echo "Building MMP Agent in mode: $MODE"
echo "================================"

case $MODE in
  v1)
    echo ""
    echo "Building V1 (legacy)..."
    echo "------------------------"
    go build -o ./bin/mmp-v1 .
    echo "✓ V1 build complete: ./bin/mmp-v1"
    ;;

  v2)
    echo ""
    echo "Building V2 (modern)..."
    echo "------------------------"

    # Check if node_modules exists, if not run npm install
    if [ ! -d "frontend/node_modules" ]; then
      echo "Installing frontend dependencies..."
      cd frontend && npm install && cd ..
    fi

    # Frontend build
    echo "Building React frontend..."
    cd frontend && npm run build && cd ..
    echo "✓ Frontend build complete"

    # Backend build
    echo "Building Go backend..."
    go build -o ./bin/mmp-v2 ./cmd/command/
    echo "✓ V2 build complete: ./bin/mmp-v2"
    ;;

  both)
    echo ""
    echo "Building both versions..."
    echo "------------------------"

    # Build V1
    echo "Building V1..."
    go build -o ./bin/mmp-v1 .
    echo "✓ V1 build complete"

    # Build V2
    echo ""
    echo "Building V2..."
    if [ ! -d "frontend/node_modules" ]; then
      echo "Installing frontend dependencies..."
      cd frontend && npm install && cd ..
    fi
    echo "Building React frontend..."
    cd frontend && npm run build && cd ..
    echo "✓ Frontend build complete"

    echo "Building Go backend..."
    go build -o ./bin/mmp-v2 ./cmd/command/
    echo "✓ V2 build complete"
    ;;

  desktop)
    echo ""
    echo "Building V2 Desktop App..."
    echo "------------------------"

    # Check if node_modules exists
    if [ ! -d "frontend/node_modules" ]; then
      echo "Installing frontend dependencies..."
      cd frontend && npm install && cd ..
    fi

    # Frontend build
    echo "Building React frontend..."
    cd frontend && npm run build && cd ..
    echo "✓ Frontend build complete"

    # Desktop build
    echo "Building desktop application..."
    go build -o ./bin/mmp-desktop ./cmd/desktop/
    echo "✓ Desktop build complete: ./bin/mmp-desktop"
    ;;

  *)
    echo "Unknown mode: $MODE"
    echo ""
    echo "Usage: $0 [v1|v2|both|desktop]"
    echo ""
    echo "  v1       - Build legacy V1 version (Go only)"
    echo "  v2       - Build modern V2 version (React + Go) [default]"
    echo "  both     - Build both V1 and V2"
    echo "  desktop  - Build V2 desktop application (Wails)"
    echo ""
    exit 1
    ;;
esac

echo ""
echo "================================"
echo "Build complete!"
echo "================================"
echo ""
echo "To run:"
case $MODE in
  v1)
    echo "  ./bin/mmp-v1"
    ;;
  v2)
    echo "  ./bin/mmp-v2"
    ;;
  both)
    echo "  V1: ./bin/mmp-v1"
    echo "  V2: ./bin/mmp-v2"
    ;;
  desktop)
    echo "  ./bin/mmp-desktop"
    ;;
esac
echo ""
