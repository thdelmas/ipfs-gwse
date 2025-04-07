# build stage
FROM node:18 AS build
WORKDIR /app
COPY . .
RUN npm install && npm run build

# serve stage (optional, if you want to serve static files)
FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html

