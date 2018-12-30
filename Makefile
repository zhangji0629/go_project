
GOPATH:=$(CURDIR)

LINUXARG=CGO_ENABLED=0 GOOS=linux GOARCH=amd64

BUILDARG=-ldflags " -s -X main.buildtime=`date '+%Y-%m-%d_%H:%M:%S'` -X main.githash=`git rev-parse HEAD`"

export GOPATH

dep:
	cd src; glide install; cd -


