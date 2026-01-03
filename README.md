# avrp
single executable program for aframe-vr-player with update feature

# usage
```
Usage of avrp:
  -dev
        starts in dev mode, serves 'index.html' from current directory
  -dir string
        path to video files
  -get-ffmpeg
        downloads ffmpeg
  -no-thumb
        disable thumbnail generation for current session
  -port int
        port to serve on (default 5000) (default 5000)
  -reset
        removes all configs, thumbnails & 'aframe-vr-player' files
  -sha string
        Optional - download a specific commit of aframe-vr-player (default "latest")
  -toggle-thumb
        toggles disable/enable thumbnail generation
  -update
        checks & downloads the latest version(commit) of 'aframe-vr-player'
```

### Examples
```sh
# -flag
# --flag
# -flag=x

# deletes all configs of avrp and aframe-vr-player / factory reset
avrp --reset

# serve video directory 
avrp --dir "C:\\Users\\User\\Video\\"
# linux path
avrp --dir "/home/user/Videos/"

# downloads the latest version(commit) of aframe-vr-player
avrp --update

# downloads https://github.com/mysterion/aframe-vr-player/tree/6757b3a6266be91c57e6236505249508b3e7d9af 
# https://github.com/mysterion/aframe-vr-player/commits/main/
avrp --update --sha 6757b3a6266be91c57e6236505249508b3e7d9af


```
# installation
## go cli
```sh
# install ffmpeg and add it to path
go install github.com/mysterion/avrp@latest
```
## standalone binary
[https://github.com/mysterion/avrp/releases](https://github.com/mysterion/avrp/releases)

