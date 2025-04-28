# Compile stage
FROM golang:1.24-bookworm AS build-env
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
RUN go install github.com/charmbracelet/glow@latest

WORKDIR /srv/random-fit
