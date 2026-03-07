#!/bin/bash


echo "Starting axiomMP Microservices..."

echo "Booting MathCompiler (Python)..."
cd MathCompiler
source venv/bin/activate
uvicorn main:app --port 8081 &
PID_MATH=$!
cd ..

sleep 1

echo "Booting PhysicsEngine (Julia)..."
julia --project=PhysicsEngine PhysicsEngine/server.jl &
PID_PHYSICS=$!

sleep 5

echo "Booting GameServer (Go)..."
cd GameServer
go run main.go &
PID_SERVER=$!
cd ..

sleep 2

echo "Booting GameClient (Rust Conduit)..."
cd GameClient
cargo run --release


echo "Shutting down the universe..."
kill $PID_MATH
kill $PID_PHYSICS
kill $PID_SERVER

echo "axiomMP Offline."
