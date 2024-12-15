module github.com/drtcrz23/Project_Booking/HotelService

go 1.23

require (
	github.com/joho/godotenv v1.5.1
	github.com/mattn/go-sqlite3 v1.14.24
	golang.org/x/sync v0.9.0
	google.golang.org/grpc v1.68.0
)

require (
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace github.com/drtcrz23/Project_Booking/services/grpc => ./internal/grpc
