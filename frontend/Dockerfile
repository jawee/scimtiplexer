FROM node:22 AS frontend_builder
WORKDIR /frontend

COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:alpine AS frontend
COPY --from=frontend_builder /frontend/dist /usr/share/nginx/html
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
