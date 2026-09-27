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

buildforlinux: generate
    mkdir -p bin
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/eriksoftware-linux

# Build and deploy the Linux binary, then restart the service.
deploy: buildforlinux
    scp bin/eriksoftware-linux erik@erikjermanis.me:/home/erik/sites/eriksoftware.hr/eriksoftware-tmp
    ssh erik@erikjermanis.me 'mv /home/erik/sites/eriksoftware.hr/eriksoftware-tmp /home/erik/sites/eriksoftware.hr/eriksoftware'
    ssh erik@erikjermanis.me 'sudo /usr/bin/systemctl restart eriksoftware.service'
