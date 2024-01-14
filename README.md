# MMP - Agent
The backend for MMP


## configuration

```toml	
port = 8000 # port to run the server on
server_hostname = "localhost"
library_path = "/library" # path to the stl library
max_render_workers = 5 # max number of workers to render the 3d model images in parallel shouldn't exceed the number of cpu cores
file_blacklist = [".potato",".example"] # list of files to ignore when searching for stl and assets files in the library_path
model_render_color = "#167DF0" # color to render the 3d model
model_background_color =  "#FFFFFF"  # color to render the 3d model background
thingiverse_token = "your_thingiverse_token" # thingiverse token to allow the import of thingiverse projects

```

## How to contribute
You can contribute by opening an issue / pull request to this repo

You can have an overview of the development status in [here](https://github.com/orgs/Maker-Management-Platform/projects/1)

## Join us for discussion
![Discord Shield](https://discordapp.com/api/guilds/1013417395777450034/widget.png?style=shield)

Join discord if you need any support https://discord.gg/SqxKE3Ve4Z


