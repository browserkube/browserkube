# Video recorder

## Description
Video recording mechanism rely on ffmpeg library for grabbing and encoding video, as well as on linux syscalls to shut down recording gracefully.

List env variables:
```
#has default values
VIDEO_SIZE
FRAME_RATE
DISPLAY_NUM
CODEC
FILE_NAME
#required
SAVE_VIDEO_ENDPOINT
SESSION_ID
```
Env examples:
```
VIDEO_SIZE=1360x1020
FRAME_RATE=12
DISPLAY_NUM=99
CODEC=libx264
FILE_NAME=myvideo.mp4
SAVE_VIDEO_ENDPOINT=s3://browserkube?region=us-west-1&endpoint=localhost:9000&disableSSL=true&s3ForcePathStyle=true
SESSION_ID=a98188c7-a2ba-4ff5-abd5-08ef3c15473f
```

## Build
To build using docker use commands from **parent directory**: 
```
#windows
docker build -t <registry> -f .\recorder\Dockerfile .

```


## TODOs:
- [ ] Change docker image to distroless
- [ ] Pipe ffmpeg output directly to the storage, without using local filestore