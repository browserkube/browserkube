FROM node:gallium-alpine3.18 AS builder
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

## Copy package.json and install dependencies first
## This later will be cached and won't executed unless
## package.json or lock are changed
COPY app/package.json app/package-lock.json* ./
RUN npm install

## Copy the rest and build the app
COPY app/. /usr/src/app/.
RUN npm run start