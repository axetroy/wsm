# builder for frontend
FROM node:12.22.1-alpine AS builder

WORKDIR /app

COPY ./ .

RUN rm -rf ./.nuxt && \
  yarn config set registry https://registry.npm.taobao.org && \
  yarn && \
  yarn cache clean && \
  npm cache clean --force && \
  npm run build

# target
FROM node:12.22.1-alpine

WORKDIR /app

COPY package.json yarn.lock ./

RUN yarn config set registry https://registry.npm.taobao.org && \
  yarn --production --ignore-scripts && \
  yarn cache clean && \
  npm cache clean --force

COPY nuxt.config.js ./
COPY static ./static
COPY --from=builder /app/.nuxt ./.nuxt

ENV PORT=80
ENV HOST=0.0.0.0
ENV NODE_ENV=production

EXPOSE 80

CMD ["npm", "start"]
