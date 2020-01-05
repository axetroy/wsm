# builder for frontend
FROM node:12.14.0-alpine AS builder

WORKDIR /app

ADD ./ .

RUN rm -rf ./.nuxt && \
  yarn config set registry https://registry.npm.taobao.org && \
  yarn && \
  yarn cache clean && \
  npm cache clean --force && \
  npm run build

# target
FROM node:12.14.0-alpine

WORKDIR /app

RUN mkdir -p /app

COPY package.json yarn.lock ./

RUN yarn config set registry https://registry.npm.taobao.org && \
  yarn --production --ignore-scripts && \
  yarn cache clean && \
  npm cache clean --force

COPY nuxt.config.js ./
COPY .nuxt ./.nuxt
COPY static ./static
COPY --from=builder /app/.nuxt ./.nuxt

ENV PORT=80

EXPOSE 80

CMD ["npm", "start"]
