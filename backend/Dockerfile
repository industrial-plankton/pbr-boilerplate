###
# Author: Ethan Dinnen <ethanjdinnen@gmail.com>

# Since: 08/2018
#
# Docker recipe to build and run into a container
###
FROM golang:latest

###
# Build arguments
###

ARG APP_ENV
ARG GO_ENV

###
# Environment Variables
###

ENV GO_DEBUG="1"

WORKDIR /go/src/app
COPY . .
RUN go install -v ./...

# Expose port 3030 to the outside world
EXPOSE 3030

CMD ["backend"]
