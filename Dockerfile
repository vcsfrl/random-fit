FROM golang:1.26-bookworm AS base
ARG username
ARG exec_user_id
RUN groupadd -g $exec_user_id -o $username \
    && useradd -r -u $exec_user_id -g $username $username -m \
    && curl -fsSL https://getmic.ro | bash \
    && mv micro /usr/bin \
    && mkdir -p /srv/random-fit \
    && chown $username:$username /srv/random-fit -R
USER $username:$username
RUN go install github.com/go-delve/delve/cmd/dlv@v1.26.1 \
    && go install golang.org/x/text/cmd/gotext@v0.36.0
WORKDIR /srv/random-fit

FROM base AS build
ARG username
USER root:root
COPY . /srv/random-fit
RUN chown $username:$username /srv/random-fit -R
USER $username
RUN go mod tidy \
    && go mod vendor \
    && go generate github.com/vcsfrl/random-fit \
    && go generate github.com/vcsfrl/random-fit/cmd/translations \
    && go build -o ./bin/app ./main.go

FROM base AS prod
ARG username
USER root:root
COPY --from=build /srv/random-fit/bin/app /srv/random-fit/bin/app
RUN chown $username:$username /srv/random-fit -R
USER $username

