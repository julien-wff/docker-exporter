FROM node:20-alpine AS build

WORKDIR /app

COPY package.json /app
RUN npm install && npm cache clean --force
COPY . /app
RUN npm run build

FROM node:20-alpine AS production

ENV NODE_ENV=production
ENV SERVER_PORT=9100

WORKDIR /app

COPY --from=build /app/package.json /app
RUN npm install --production && npm cache clean --force

COPY --from=build /app/build /app/build

HEALTHCHECK --interval=30s --timeout=5s --retries=3 --start-period=5s CMD wget -q -O - http://localhost:${SERVER_PORT}/health || exit 1

CMD ["node", "build"]
