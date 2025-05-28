module gomike

go 1.21.0

require (
	google.golang.org/protobuf v1.35.2
	lib v0.0.0
)

require github.com/google/uuid v1.6.0

require github.com/golang/protobuf v1.5.0 // indirect

replace lib => ../../lib
