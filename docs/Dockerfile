FROM node:18-alpine

WORKDIR /app

RUN npm install -g swagger-markdown && \
    apk add --no-cache darkhttpd curl

COPY start.sh /app/start.sh

RUN chmod +x /app/start.sh

CMD ["/app/start.sh"]
