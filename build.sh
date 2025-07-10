#!/bin/bash
# DevGen CLI Build Script

set -e

case "${1:-build}" in
    "build")
        echo "🔨 Building DevGen CLI..."
        mkdir -p build
        go build -o build/devgen main.go
        echo "✅ Build complete: build/devgen"
        ;;
    "clean")
        echo "🧹 Cleaning..."
        rm -rf build
        rm -f devgen
        echo "✅ Clean complete"
        ;;
    "install")
        echo "📦 Installing..."
        mkdir -p build
        go build -o build/devgen main.go
        sudo cp build/devgen /usr/local/bin/
        echo "✅ Installed to /usr/local/bin/devgen"
        ;;
    "test")
        echo "🚀 Testing..."
        mkdir -p build
        go build -o build/devgen main.go
        ./build/devgen --help
        ;;
    *)
        echo "DevGen CLI Build Script"
        echo "Usage: ./build.sh [build|clean|install|test]"
        ;;
esac
