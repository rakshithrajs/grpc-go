module github.com/rakshithrajs/cloud/MMS

go 1.26.4

require (
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.12.3
	google.golang.org/grpc v1.82.1
	google.golang.org/protobuf v1.36.11
)

require (
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260414002931-afd174a4e478 // indirect
)

replace github.com/rakshithrajs/cloud/UMS => ../UMS
