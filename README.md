# 🎥 ASCII Converter  
Convert images, GIFs, videos and webcam steam into ASCII art in the terminal or save them as files.

Ascii converter is tool that can be used to convert/preview images(png,jpg), gifs, video(mp4,avi,mov,webm) and web cameras in ascii art mode

# GIF/Image to ASCII Comparison

Here’s a side-by-side look at the conversion:

| Original | ASCII (no color) | ASCII (colored) |
|--------------|------------------|-----------------|
| ![Original](./examples/cig.gif) | ![ASCII BW](./examples/ascii_cig.gif) | ![ASCII Color](./examples/ascii_color_cig.gif) |
| ![Original](./examples/kame.gif) | ![ASCII BW](./examples/ascii_kame.gif) | ![ASCII Color](./examples/ascii_color_kame.gif) |
| ![Original](./examples/anime.jpg) | ![ASCII BW](./examples/ascii_anime.png) | ![ASCII Color](./examples/ascii_color_anime.png) |
| ![Original](./examples/test.jpg) | ![ASCII BW](./examples/ascii_test.png) | ![ASCII Color](./examples/ascii_color_test.png) |


## 📖 Usage

```bash
Usage: ascii-cli <command> [options]

Commands:
  convert   Convert image/gif/video to ASCII
  preview   Preview ASCII frames in terminal
  camera    Preview/convert camera ASCII frames


## 🎥 CLI Demo

Here’s a quick demo of the CLI converting a video into ASCII art:

![CLI Demo](./examples/video.gif)

In this example, the tool takes a video input and outputs an ASCII-rendered version in real time.

### Try it yourself
```bash
ac -i input.mp4 -o output.mp4 -C
