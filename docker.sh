# Simple shell script to create and run container
docker build -f Dockerfile.dev  -t onemdp-dev-1 . && docker run -it -p 8080:8080 --name onemdp-dev-1 onemdp-dev-1