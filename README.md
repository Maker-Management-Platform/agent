# MMP v2 

Under development, so expect an ugly UI and a lot of missing features from previous versions.   
The v2 version features a new filesystem abstraction with the mindset that everything is an asset [https://github.com/Maker-Management-Platform/docs/issues/30], and the structure is just a three of assets that you can drill through.   
The frontend is now embeded in the same docker, dropping react in favour of htmx and templ. This aims to simplify development ans as a bonus easier integration with NAS OS's.   
There are also plans to distribute MMP as an installable application using wails.   


This version is highly experimental, so feedback is highly important.

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

