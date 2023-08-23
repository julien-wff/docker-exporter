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

ARG BUILD_DATE
ARG VCS_REF
ARG VCS_URL
ARG VERSION

LABEL org.opencontainers.image.created=$BUILD_DATE \
      org.opencontainers.image.authors="Julien W <cefadrom1@gmail.com>" \
      org.opencontainers.image.source=$VCS_URL \
      org.opencontainers.image.version=$VERSION \
      org.opencontainers.image.revision=$VCS_REF \
      org.opencontainers.image.title="docker-exporter" \
      org.opencontainers.image.description="A prometheus exporter to expose docker metrics that are normally not available through cAdvisor" \
      org.opencontainers.image.base.name="node:20-alpine" \
      org.opencontainers.image.base.version="20-alpine"

HEALTHCHECK --interval=30s --timeout=5s --retries=3 --start-period=5s CMD wget -q -O - http://localhost:${SERVER_PORT}/health || exit 1

CMD ["node", "build"]
