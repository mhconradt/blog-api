FROM golang:1.13rc-alpine3.0

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/mhconradt/blog-api

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download all the dependencies
# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

ENV PORT=8080

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["blog-api"]
