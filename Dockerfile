# Sử dụng bản Go official làm base image
FROM golang:1.20-alpine as builder

# Thiết lập working directory bên trong container
WORKDIR /app

# Sao chép go.mod và go.sum vào thư mục hiện tại của container
# và tải về các dependency
COPY go.mod go.sum ./
RUN go mod download

# Sao chép mã nguồn của ứng dụng vào container
COPY . .

# Build ứng dụng Go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Sử dụng scratch làm base image cho image cuối cùng
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Sao chép file thực thi từ builder stage
COPY --from=builder /app/main .

# Cổng mà ứng dụng sẽ chạy trên đó
EXPOSE 8080

# Chạy ứng dụng
CMD ["./main"]
