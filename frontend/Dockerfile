FROM node:18-alpine AS build
WORKDIR /app
COPY package.json ./
COPY package-lock.json ./
RUN npm install
RUN npm install axios
COPY . ./
CMD [ "npm", "run", "dev", "--", "--host"]