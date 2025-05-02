FROM golang:1.24-bookworm AS dev
ARG username
ARG exec_user_id
RUN groupadd -g $exec_user_id -o $username
RUN useradd -r -u $exec_user_id -g $username $username -m
RUN curl https://getmic.ro | bash
RUN mv micro /usr/bin
RUN mkdir -p /srv/random-fit
RUN chown $username:$username /srv/random-fit -R
USER $username:$username
RUN go install github.com/divan/expvarmon@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
WORKDIR /srv/random-fit

FROM golang:1.24-bookworm AS build
RUN mkdir /srv/random-fit
WORKDIR /srv/random-fit
COPY . /srv/random-fit
RUN go mod tidy
RUN go mod vendor
RUN go generate github.com/vcsfrl/random-fit/cmd;
RUN go build -o ./bin/app ./cmd/main.go;

FROM golang:1.24-bookworm AS prod
COPY --from=build /srv/random-fit/bin/app /srv/random-fit/bin/app
RUN mkdir -p /srv/random-fit/data
RUN curl https://getmic.ro | bash
RUN mv micro /usr/bin
WORKDIR /srv/random-fit