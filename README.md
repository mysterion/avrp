# avrp
single executable program for aframe-vr-player with update feature

# usage
```
  -dir string
        path to video files

  -update
        checks & downloads the latest version of 'aframe-vr-player'

  -dev
        starts in dev mode, serves 'index.html' from current directory

  -reset
        removes all configs, thumbnails & 'aframe-vr-player' files
```
# installation
## go cli
```sh
# install ffmpeg and add it to path
go install github.com/mysterion/avrp@latest
```
## standalone binary
Download the `with-ffmpeg` version for the thumbnails feature
[https://github.com/mysterion/avrp/releases](https://github.com/mysterion/avrp/releases)

