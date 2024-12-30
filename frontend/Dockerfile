FROM node:16-alpine AS build

RUN apk add --no-cache \
    libvips-dev \
    build-base \
    python3 \
    make \
    g++ \
    && npm install -g npm@latest \
    && npm install -g yarn

WORKDIR /app

COPY package.json yarn.lock ./

RUN yarn install --frozen-lockfile \
    && npx update-browserslist-db@latest

COPY . .

RUN yarn build

FROM nginx:alpine

COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
