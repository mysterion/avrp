rm -rf build

function non_zero_exit() {
    if [[ $1 -ne 0 ]]; then
        exit 1
    fi
}

which go
non_zero_exit

mkdir build
go build -o build/ cmd/avrp/avrp.go

if [[ -z "${FFMPEG_BUILD}" ]]; then
    echo '${FFMPEG_BUILD} not set.'
    echo 'skipping ffmpeg binaries'
    mkdir -p build/thirdparty
    cp -r thirdparty/*.md build/thirdparty/
else
    mkdir -p build/thirdparty
    cp -r thirdparty/ffmpeg* thirdparty/ffprobe* build/thirdparty/
    rm build/thirdparty/*.go
    rm build/thirdparty/*.md
fi

