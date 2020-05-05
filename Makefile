gen-calculatorpb:
	protoc calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:.
run-calculator-server:
	go run calculator/server/server.go
run-calculator-client:
	go run calculator/client/client.go
	
gen-platformpb:
	protoc platform/platformpb/platform.proto --go_out=plugins=grpc:.
run-platform-server:
	go run platform/server/server.go
run-platform-client:
	go run platform/client/client.go
