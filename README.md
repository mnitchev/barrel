# barrel
Repository for the garden hell week 🔥

# Bild 
```command
GOOS=linux GOARCH=amd64 go build -o barrel
```

# Run
```command
./barrel roll <command> [ -- <args>* ] 
```
Example:
```command
./barrel roll /bin/sh -- "-c" "hostname foo; hostname"
```

# Test on MacOS
```
docker run --privileged -v $PWD:/app --rm -it mnitchev/barrel-test  ginkgo -mod=vendor -r
```
