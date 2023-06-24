FROM golang:latest as builder

COPY . ./app

WORKDIR /go/app

RUN make compile ENVIROMENT=production

FROM mcr.microsoft.com/azure-functions/base:4 as runtime-image

ENV AzureWebJobsScriptRoot=/home/app \
    AzureFunctionsJobHost__Logging__Console__IsEnabled=true

RUN apt install -y libc6 libc6-dev libc6-dbg

COPY --from=builder ["/go/app/functions", "/home/app"]