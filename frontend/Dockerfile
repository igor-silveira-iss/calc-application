# Build stage: compiles the React app into static files.
FROM node:22-alpine AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
# The API is reached through the nginx proxy on the same origin.
ENV VITE_API_URL=/api
RUN npm run build

# Runtime stage: nginx serves the app and forwards /api to the backend.
FROM nginx:1.28-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
