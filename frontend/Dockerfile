# syntax=docker/dockerfile:1.7

ARG NODE_VERSION=20.18.1

FROM node:${NODE_VERSION}-alpine AS builder
WORKDIR /app

RUN apk add --no-cache bash

COPY frontend/package.json frontend/package-lock.json* ./
RUN npm install

COPY frontend .
COPY docs /docs

ARG VITE_API_BASE_URL=http://localhost:8080/v1
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}
ENV SPEC_PATH=/docs/api.yaml

RUN npm run build

FROM nginx:1.29.3-alpine
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
