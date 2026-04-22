FROM node:20-alpine AS builder

WORKDIR /app
COPY package.json ./
RUN npm install

COPY . .
# Point API at the Go service name inside docker network
ARG VITE_API_URL=http://server:8080
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
# SPA fallback — redirect all routes to index.html for react-router
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
