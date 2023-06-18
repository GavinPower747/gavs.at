FROM golang:1.19 as builder

COPY . ./app

WORKDIR /go/app

RUN make compile ENVIROMENT=production

FROM mcr.microsoft.com/azure-functions/base:4 as runtime-image

ENV AzureWebJobsStorage=$AzureWebJobsStorage

ENV AzureWebJobsScriptRoot=/home/app \
    AzureFunctionsJobHost__Logging__Console__IsEnabled=true

COPY --from=builder ["/go/app/functions", "/home/app"]