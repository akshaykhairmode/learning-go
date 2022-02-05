#Dockerfile
#From which image we want to build. This is basically our environment.
FROM golang:1.17.6-alpine as Build

#This will copy all the files in our repo to the inside the container at root location.
COPY . . 

#build our binary at root location. Binary name will be main. We are using go modules so gpath env variable should be empty.
RUN GOPATH= go build -o /main main.go 

#This will use scratch image which is the smallest image we can have. We will use scratch because we needed go environment only for building.
FROM scratch

#we copy our binary from build to scratch.
COPY --from=Build main main

#we tell docker what to run when this image is run and run it as executable.
ENTRYPOINT [ "/main" ]