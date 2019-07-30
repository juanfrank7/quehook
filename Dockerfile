FROM golang
ADD . app/
WORKDIR app/
EXPOSE 8080
CMD go run main.go
