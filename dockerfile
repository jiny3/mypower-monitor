FROM --platform=$BUILDPLATFORM golang:latest AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=$TARGETOS 
ENV GOARCH=$TARGETARCH

RUN go mod tidy && go build -o bin/ucasnj-smi .

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/bin/ /app/bin/
COPY ./static ./static

EXPOSE 8080

ENTRYPOINT ["bin/ucasnj-smi", "server"]