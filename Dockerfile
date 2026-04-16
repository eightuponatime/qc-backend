FROM node:20-alpine AS frontend
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY frontend/ ./frontend/
COPY vite.config.js ./
RUN npm run build

FROM golang:1.25-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/static/dist ./static/dist
RUN go build -o server ./cmd

FROM alpine:3.19
WORKDIR /app
COPY --from=backend /app/server .
COPY --from=backend /app/static/dist ./static/dist
COPY --from=backend /app/templates ./templates
COPY --from=backend /app/i18n ./i18n

EXPOSE 8080
CMD ["./server"]