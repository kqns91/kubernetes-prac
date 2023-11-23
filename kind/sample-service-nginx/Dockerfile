# syntax=docker/dockerfile:1

FROM nginx:1.25.3-alpine
COPY ./nginx.conf /etc/nginx/nginx.conf
COPY ./index.html /usr/share/nginx/html/index.html

EXPOSE 80
