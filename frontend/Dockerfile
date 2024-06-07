FROM node:gallium-alpine3.18 AS builder
WORKDIR /app

## Copy package.json and install dependencies first
## This later will be cached and won't executed unless
## package.json or lock are changed
COPY app/package.json app/package-lock.json* ./
RUN npm install --omit=dev

## Copy the rest and build the app
COPY app/. /app/.
RUN npm run build

FROM nginxinc/nginx-unprivileged:1.25.2-alpine
COPY app-nginx.conf /etc/nginx/conf.d/app-nginx.conf
COPY --from=builder /app/dist /usr/share/nginx/html

#RUN chown nginx.nginx /usr/share/nginx/html/ -R
EXPOSE 8080