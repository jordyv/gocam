# Gocam - IP camera alert tool written in Go #

Simple script to get an image from the IP camera every x seconds, calculate the unique hash with the [perception hashing](https://en.wikipedia.org/wiki/Perceptual_hashing) algorithm and compare it with the previous capture. When the configured treshold is reached, it will trigger an action.

Currently supported alerts/actions:
 - Write a line to a log file
 - Use ImageMagick's 'compare' tool to create a difference image from the snapshots

## Installation ##

 - Clone this repository
 - Run `make`
 - The binary is compiled in the `dist` folder

## Configuration ##

Create a `gocam.yaml` file in `/etc` or at your home folder. Check `gocam.example.yaml` for an example.

```yaml
cameraURL: "<< YOUR CAMERA IMAGE URL >>"
imagePath: "/home/user/camera"
verbose: false
interval: 5s
maxKeepedImageFiles: 10
treshold: 4

alertHandlers:
  - console
  - diff
  - telegram

```
