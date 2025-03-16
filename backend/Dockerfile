# สร้าง container จาก Go runtime โดยใช้ multi-stage build เพื่อลดขนาดของ image
FROM golang:1.23.2-alpine AS builder

# ติดตั้ง dependencies ที่จำเป็น
RUN apk add --no-cache git

# ตั้งค่า working directory
WORKDIR /app

# คัดลอกเฉพาะไฟล์ที่จำเป็นสำหรับการติดตั้ง dependencies ก่อน
# เพื่อใช้ประโยชน์จาก Docker cache
COPY go.mod go.sum ./
RUN go mod download

# คัดลอกโค้ดทั้งหมด
COPY . .

# สร้าง binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -o auth-service main.go

# สร้าง container ขนาดเล็กสำหรับการใช้งานจริง
FROM alpine:3.18

WORKDIR /app

# # คัดลอกไฟล์การตั้งค่า
COPY --from=builder /app/config/model.conf /app/config/model.conf
COPY --from=builder /app/config/policy.csv /app/config/policy.csv

# คัดลอก binary จาก stage แรก
COPY --from=builder /app/auth-service .

# เปิด port 8081
EXPOSE 8081

# รัน service
CMD ["./auth-service"]