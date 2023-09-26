FROM node:19.8.1-alpine as build
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY dist /app/dist
#RUN npm install && npm run build

FROM nginx:1.16.0-alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY deploy/nginx.conf /etc/nginx/conf.d
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
