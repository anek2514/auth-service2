# Build stage
FROM node:18-alpine AS build

WORKDIR /app

# ติดตั้ง dependencies
COPY package.json package-lock.json ./
RUN npm ci

# คัดลอกโค้ดทั้งหมด
COPY . .

# สร้าง production build
RUN npm run build

# Production stage
FROM nginx:alpine

# คัดลอก production build จาก stage แรก
COPY --from=build /app/build /usr/share/nginx/html
# COPY --from=build /app/build /etc/nginx/html

# คัดลอก nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]dock