default: dev

# Generate Go code from pages/*.templ (templ is pinned in go.mod).
generate:
    go run -mod=mod github.com/a-h/templ/cmd/templ generate

# Start locally at http://localhost:8080 (or set PORT).
dev: generate
    go run .

# Create a self-contained binary with the embedded public/ assets.
build: generate
    mkdir -p bin
    go build -o bin/eriksoftware .

test: generate
    go test ./...

push:
    git add .
    git commit -m "push"
    git push
