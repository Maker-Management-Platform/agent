# MMP v2 

Under development, so expect an ugly UI and a lot of missing features from previous versions.   
The v2 version features a new filesystem abstraction with the mindset that everything is an asset [https://github.com/Maker-Management-Platform/docs/issues/30], and the structure is just a three of assets that you can drill through.   
The frontend is now embeded in the same docker, dropping react in favour of htmx and templ. This aims to simplify development ans as a bonus easier integration with NAS OS's.   
There are also plans to distribute MMP as an installable application using wails.   


This version is highly experimental, so feedback is highly important.


## Configuration File
The config file is located in data/config.toml
### Filesystem:
You can configure multiple filesystems aka library folders of different kinds:
#### LocalFS
The standard file system
``` toml
[[library.filesystems]]
kind = 'local'
name = 'Library'
path = '/library'
```

#### GitFS
Uses git as a read only filesystem.   
The internal library used is [go-fsimpl](https://pkg.go.dev/github.com/hairyhenderson/go-fsimpl/gitfs), please refer to the documentation about environment variables and credentials.
``` toml
[[library.filesystems]]
kind = 'gitfs'
name = 'ExampleGit'

[library.filesystems.config]
url = 'https://github.com/Rat-Rig/RatRig-PrintedParts'

```

#### S3
Not yet implemented, please reach out if you can provide a set of credentials for me to develop and test wit

``` toml
tba
```

#### Azure Blob Storage
Not yet implemented, please reach out if you can provide a set of credentials for me to develop and test with

``` toml
tba
```

#### Google Storage
Not yet implemented, please reach out if you can provide a set of credentials for me to develop and test with

``` toml
tba
```

## Setting up for development
Run npm install in the frontend folder

install air https://github.com/air-verse/air   
`go install github.com/air-verse/air@latest`

run `air --  -data-folder data` in the repo folder

## Join us for discussion
![Discord Shield](https://discordapp.com/api/guilds/1013417395777450034/widget.png?style=shield)  
Join discord if you need any support https://discord.gg/SqxKE3Ve4Z


## General Documentation (old)
[here](https://github.com/Maker-Management-Platform/docs)

## Docker compose

```yaml	
---
services:
  agent:
    image: ghcr.io/maker-management-platform/agent:v2
    container_name: mmp
    volumes:
      - ./library:/library # should contain your project library
      - ./data:/data # will contain config and state files
    ports:
      - 8000:8000 
    restart: unless-stopped
```

