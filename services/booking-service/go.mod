module BookingService

go 1.23.4

require (
	github.com/drtcrz23/Project_Booking/services/grpc v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/segmentio/kafka-go v0.4.47
	golang.org/x/sync v0.9.0
	google.golang.org/grpc v1.69.0
)

require (
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
	google.golang.org/protobuf v1.35.2 // indirect
)

replace github.com/drtcrz23/Project_Booking/services/grpc => ../grpc
