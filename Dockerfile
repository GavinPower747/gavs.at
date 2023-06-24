FROM golang:alpine as builder

COPY . ./app

WORKDIR /go/app

RUN apk add --update make
RUN make compile ENVIROMENT=production

FROM mcr.microsoft.com/azure-functions/base:4 as runtime-image

ENV AzureWebJobsScriptRoot=/home/app \
    AzureFunctionsJobHost__Logging__Console__IsEnabled=true

COPY --from=builder ["/go/app/functions", "/home/app"]