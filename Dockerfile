FROM golang:latest AS build
WORKDIR /app
COPY . .
RUN go build -o ascii-web . 

FROM debian
COPY --from=build /app /app
WORKDIR /app
ENTRYPOINT [ "./ascii-web" ]
EXPOSE 8080
