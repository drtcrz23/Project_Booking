module BookingService

go 1.23.4

require (
	github.com/drtcrz23/Project_Booking/services/grpc v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/segmentio/kafka-go v0.4.47
	golang.org/x/sync v0.10.0
	google.golang.org/grpc v1.69.0
)

require (
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	golang.org/x/net v0.32.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241209162323-e6fa225c2576 // indirect
	google.golang.org/protobuf v1.35.2 // indirect
)

replace github.com/drtcrz23/Project_Booking/services/grpc => ./../hotel-service/internal/grpc
