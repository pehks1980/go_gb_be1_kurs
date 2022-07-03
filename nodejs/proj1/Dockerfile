FROM node:alpine
#PI
#FROM hypriot/rpi-node:latest
# Create app directory
WORKDIR /app
# Install app dependencies
ADD package.json package-lock.json /app/

RUN npm install
# If you are building your code for production
#RUN npm ci --only=production

# Bundle app source
COPY . /app

EXPOSE 8090 8080
CMD [ "node", "server.js" ]

# docker build -t nodejs .
#
# build image for raspberry pi arm64:
# docker buildx build --platform linux/arm64 -t nodejs:arm64 .
#
# docker run -p 8080-8090:8080-8090 nodejs
# docker run -it -p 8080-8090:8080-8090/tcp -network host nodejs
# image строится но при запуске все зависает...
