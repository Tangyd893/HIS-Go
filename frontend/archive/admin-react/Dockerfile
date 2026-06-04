# React 前端构建器（产物供 Nginx 挂载）
FROM node:22-alpine

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm ci --legacy-peer-deps

COPY . .

RUN npm run build
