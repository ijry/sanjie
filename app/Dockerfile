FROM node:20-alpine AS build

WORKDIR /src

COPY package*.json ./
RUN npm install --legacy-peer-deps

COPY . .

ARG VITE_SANJIE_API_BASE_URL=""
ENV VITE_SANJIE_API_BASE_URL=${VITE_SANJIE_API_BASE_URL}

RUN npm run build:h5

FROM nginx:1.27-alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /src/dist/build/h5 /usr/share/nginx/html

EXPOSE 80
