FROM golang:1.19

ENV TERM=xterm-256color

SHELL ["/bin/bash", "-c"]

RUN apt-get update && apt-get install -y vagrant make

RUN wget -qO- \
        https://github.com/charmbracelet/gum/releases/download/v0.10.0/gum_0.10.0_Linux_x86_64.tar.gz \
        | tar -xz -C /usr/bin/ gum

RUN curl -sSL https://get.docker.com | bash

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./files/bashrc /root/.bashrc
COPY . .
RUN make build
RUN make install

CMD ["violet"]
