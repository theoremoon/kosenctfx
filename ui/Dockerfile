FROM node:15.3 AS builder

COPY . .
RUN yarn && yarn build

FROM nginx:1.19
COPY --from=builder /dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
