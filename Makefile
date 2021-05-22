VERSION:=0.0.1

git:
	git add .
	git commit -m "update"
	git push

all: clientMAC serverMAC clientWindows serverWindows clientLinux serverLinux

clientMAC:
	cd client/src && GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o ../../assets/itsClientDar_v${VERSION}

serverMAC:
	cd server/src && GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o ../../assets/itsServerDar_v${VERSION}

clientWindows:
	cd client/src && GOOS=windows GOARCH=386 go build -o ../../assets/itsClient_v${VERSION}.exe main.go

serverWindows:
	cd server/src && GOOS=windows GOARCH=386 go build -o ../../assets/itsServer_v${VERSION}.exe main.go

clientLinux:
	cd client/src && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../assets/itsClient_v${VERSION}

serverLinux:
	cd server/src && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../../assets/itsServer_v${VERSION}
