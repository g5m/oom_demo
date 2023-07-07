FROM golang:latest
WORKDIR /app
COPY . .
EXPOSE 8080
RUN export GOPROXY="https://goproxy.cn,direct"
RUN export GOSUMDB="off"
RUN export GO111MODULE=on
#RUN go mod tidy
#RUN go build -o oom_demo .
#RUN GOOS=linux GOARCH=amd64 go build
CMD [ "./oom_demo" ]