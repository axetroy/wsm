FROM node:12.14.0-alpine

WORKDIR /root/app

RUN mkdir -p /root/app

COPY package.json yarn.lock ./

RUN yarn config set registry https://registry.npm.taobao.org && \
  yarn --production --ignore-scripts && \
  yarn cache clean && \
  npm cache clean --force

COPY nuxt.config.js ./
COPY .nuxt ./.nuxt
COPY static ./static

ENV PORT=80

EXPOSE 80

CMD ["npm", "start"]