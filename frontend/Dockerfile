# Стадия сборки
FROM node:20-alpine AS builder

WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install
COPY . .
RUN npm run build

# Стадия раздачи через nginx
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

# Конфиг nginx для React Router (все пути → index.html)
RUN echo 'server { \
    listen 80; \
    location / { \
        root /usr/share/nginx/html; \
        try_files $uri $uri/ /index.html; \
    } \
}' > /etc/nginx/conf.d/default.conf

EXPOSE 80
