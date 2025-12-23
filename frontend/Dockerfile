# ============================================
# Dockerfile для Frontend
# Путь: frontend/Dockerfile
# ============================================

FROM node:18-alpine AS builder

WORKDIR /app

# Копируем package.json и package-lock.json
COPY package*.json ./

# Устанавливаем зависимости
RUN npm ci

# Копируем весь frontend
COPY . .

# Собираем production build
RUN npm run build

# ============================================
# Production образ с nginx
# ============================================
FROM nginx:alpine

# Копируем собранное приложение
COPY --from=builder /app/build /usr/share/nginx/html

# Копируем конфиг nginx (если есть)
# COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]