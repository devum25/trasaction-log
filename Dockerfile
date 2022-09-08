# stage 1: compile the binary in a containerized golang environment
FROM golang:1.18 as build

# copy the src file from the host
COPY . /src
# set the working directory to the same place as copied the code
WORKDIR /src
# build the binary!
RUN CGO_ENABLED=0 GOOS=linux go build -o kvs



# stage 2: Build key-value store image proper
#
# Use a scratch image, which contains no distribution files
FROM scratch

# Copy the binary from the build container
COPY --from=build /src/kvs .

EXPOSE 8080

CMD ["/kvs"]