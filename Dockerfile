FROM ubuntu:22.04 As builder
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y locales locales-all curl && \
   locale-gen ja_JP.UTF-8 &&\
    apt-get clean &&\
    rm -rf /var/lib/apt/lists/*
RUN curl -OL https://dl.google.com/go/go1.18.2.linux-amd64.tar.gz &&\
    tar -C /usr/local -xzf go1.18.2.linux-amd64.tar.gz &&\
    rm -rf go1.18.2.linux-amd64.tar.gz
ENV PATH $PATH:/usr/local/go/bin
WORKDIR /app
ADD . .
RUN go build -o app
FROM ubuntu:22.04
LABEL version="0.1.0"
LABEL docker_arch="amd64"
LABEL go_version="1.18.2"
LABEL name="youtubedown"
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y locales locales-all curl youtube-dl && \
   locale-gen ja_JP.UTF-8 &&\
    apt-get clean &&\
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
RUN mkdir /.cache ./download
RUN chmod 777 /app /.cache ./download
COPY ./html html
COPY ./run.sh ./
RUN chmod 777 run.sh
COPY --from=builder /app/app /app
CMD ["./app"]
