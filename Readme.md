# IntrospectedTunnelsSocket

> Introspected Tunnel for localhost

- [structure](#structure)
- [Download](#download)
- [How to use](#how-to-use)

## Structure

![structure.png](structure.png)

## Download

### Linux

- server

> curl -O https://raw.githubusercontent.com/skynocover/IntrospectedTunnelsSocket_binaryStore/master/itsServer_v0.0.3

- client

> curl -O https://raw.githubusercontent.com/skynocover/IntrospectedTunnelsSocket_binaryStore/master/itsClient_v0.0.3

### MAC

- server

> curl -O https://raw.githubusercontent.com/skynocover/IntrospectedTunnelsSocket_binaryStore/master/itsServerDar_v0.0.3

- client

curl -O https://raw.githubusercontent.com/skynocover/IntrospectedTunnelsSocket_binaryStore/master/itsClientDar_v0.0.3

### Windows

- server

> curl -O https://raw.githubusercontent.com/skynocover/IntrospectedTunnelsSocket_binaryStore/master/itsServer_v0.0.3.exe

- client

curl -O https://raw.githubusercontent.com/skynocover/IntrospectedTunnelsSocket_binaryStore/master/itsClient_v0.0.3.exe

## How to use

### 1. Build itsServer

> set .env file

```
SERVER_LISTEN=8080 # 
DOMAINS=domain.com,domain2.tk 
```

- SERVER_LISTEN: this is your port for listen request
- DOMAINS: these are your domain, you should add Aname record to the itsServer you running, itsServer will give these domain to your client

Run the itsServer on your VPS with domain

### 2. Run the itsClient on same place with your service

> set .env file

```
DOMAIN=http://localhost:8080
PROXY=http://localhost:3020
```

- DOMAIN: Your ITS Server Domain, client will connenct with socket to this URL
- PROXY: The URL your service running

### 3. Enjoy your tunnel
