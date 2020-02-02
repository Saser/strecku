module github.com/Saser/strecku/testing

go 1.13

require (
	github.com/Saser/strecku/backend v0.0.0
	github.com/golang/protobuf v1.3.3
	github.com/stretchr/testify v1.4.0
	google.golang.org/grpc v1.27.0
)

replace github.com/Saser/strecku/backend => ../backend
