FROM node:18-slim

RUN apt-get update && apt-get install -y chromium \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

EXPOSE 3001
CMD ["npm", "start"]
