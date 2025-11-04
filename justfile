clean:
    rm -f pssh

build: clean
    go build -o pssh .