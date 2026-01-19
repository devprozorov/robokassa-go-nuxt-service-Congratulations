# syntax=docker/dockerfile:1
FROM node:20-alpine AS build
RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates
WORKDIR /app
COPY package.json ./
RUN npm install
COPY . .
RUN npm run build

FROM node:20-alpine
RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates
WORKDIR /app
ENV NODE_ENV=production
COPY --from=build /app/.output ./.output
EXPOSE 3000
CMD ["node", ".output/server/index.mjs"]
