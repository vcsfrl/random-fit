# Compile stage
FROM golang:1.24-bookworm AS build-env
ARG username
ARG exec_user_id
RUN groupadd -g $exec_user_id -o $username
RUN useradd -r -u $exec_user_id -g $username $username -m

RUN mkdir -p /srv/random-fit

RUN chown $username:$username /srv/random-fit -R
USER $username:$username
RUN go install github.com/divan/expvarmon@latest

WORKDIR /srv/random-fit
