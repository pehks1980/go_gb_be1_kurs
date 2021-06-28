FROM node:15
# Create app directory
WORKDIR /app
# Install app dependencies
ADD package.json package-lock.json /app/

RUN npm install
# If you are building your code for production
#RUN npm ci --only=production

# Bundle app source
COPY . /app

EXPOSE 8090
CMD [ "node", "server.js" ]

# docker build -t nodejs .
# docker run -p 8090:8090 nodejs
# image строится но при запуске все зависает...
