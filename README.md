# grpc-learn
## Calculator api
- Tìm hiểu về unary API:
  - Tìm hiểu cách định nghĩa API bên trong file .proto
  - Thực hiện phần xử lý ở server side
  - Tìm hiểu cách client gọi server 
  
- Tìm hiểu về server streaming API:
  - Là 1 loại API mới dùng HTTP/2, client gửi 1 request và nhận được nhiều response từ server, thích hợp cho:
    - Khi server cần gửi một lượng lớn data
    - Khi server cần push thông tin cho client mà không cần client request thêm
  - Tìm hiểu cách mà client chỉ gửi 1 request mà lại nhận về được nhiều response
  - Tìm hiểu cách client cần làm để có thể nhận nhiều response từ server
  - Tìm hiểu cách để biết server streaming kết thúc

- Tìm hiểu về client streaming API:
  - Là 1 loại API mới dùng HTTP/2, client gửi nhiều request và nhận được một response từ server ở 1 lúc nào đó, thích hợp cho:
    - Khi client cần gửi một lượng lớn data
    - Khi việc xử lý server là tốn kém, chỉ thực hiện khi client gửi data
    - Khi client cần push data tới server mà không quan tâm response
  - Tìm hiểu cách mà client có thể gửi nhiều request và nhận được response từ server
  - Tìm hiểu cách server cần làm để có thể nhận nhiều request từ client
  - Tìm hiểu cách để biết client streaming kết thúc

- Bao gồm các api sau:

```
rpc Sum(SumRequest) returns (SumResponse) {}
```
  - Unary:
    - Khi client gửi một yêu cầu gồm 2 số nguyên đến server và nhận về một phản hồi giống như một cuộc gọi phương thức bình thường trả về tổng của 2 số đã gửi.

```
rpc PrimeNumberDecomposition(PNDRequest) returns (stream PNDResponse) {}
```
  - Server streaming:
    - Khi client gửi một yêu cầu gồm 1 số nguyên đến server và nhận về một stream trả về tích các thừa số nguyên tố của số đã gửi.

```
rpc Average(stream AverageRequest) returns (AverageResponse) {}
```
  - Client streaming:
    - Khi client gửi n số đến server, stream được sử dụng, client sẽ hoàn thành việc gửi thông điệp của nó, sau đó chờ server phản hồi về giá trị trung bình cộng của n số đã gửi.
    
- Run:
```
make run-calculator-server
make run-calculator-client
```
