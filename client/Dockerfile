FROM node:14.15.1
COPY package-lock.json .
COPY package.json .
COPY craco.config.js .
RUN npm install
COPY . .
CMD npm start