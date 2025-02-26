FROM node:12
#PI
#FROM hypriot/rpi-node:latest
# Create app directory
WORKDIR /usr/src/app
# Install app dependencies
COPY package*.json ./

RUN npm install
# If you are building your code for production
#RUN npm ci --only=production

# копируем исходный код
COPY . .

EXPOSE 8090
CMD [ "node", "server.js" ]

# docker build . -t pehks1980/node-web-app_1
# build image for raspberry pi arm64:
# docker buildx build --platform linux/arm64 . -t node-web-app:arm64
# docker save node-web-app:arm64  > node-web-arm64.tar
# scp node-web-arm64.tar user@192.168.1.204:/home/user
# docker load -i node-web-arm64.tar
# docker run -p 8090:8090 node-web-app
