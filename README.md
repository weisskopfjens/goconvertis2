# goConvertIS2
With this command line tool an a golang package you can extract visual images and infrared images from "Fluke" proprietary .IS2 files.

FLUKE FLIR Infrared cams use this fileformat.

It was hard work to reverser engineer the file formats.
- I used a hexeditor for analyzing the files.
- I learned the Stefan Boltzmann Law ;-)

There are two main versions of .is2 files.

- The older version is a binary format.
- The newer version is a zip with .is2 file extension

This tool can handle both.

This is a experimental tool. The temperature values are a little bit inaccurate. Maybe someone can solve this problem. This tool and package is only for study and demonstration purposes. It`s not an official FLUKE product. 


## Usage
```
goConvertIS2 by (c)Jens Weißkopf (github.com/weisskopfjens/goconvertis)
(*) are required parameter.
  -b float
        Background temperature. (default 20)
  -e float
        Emission factor. (default 0.95)
  -i string
        (*) A .is2 File.
  -max float
        Max. temperature. (default 70)
  -min float
        Min. temperature. (default 20)
  -oi string
        A .jpg file for infrared output. (default "ir.jpg")
  -ov string
        A .jpg file for visual output. (default "vis.jpg")
```

## Todo
- Write unit tests
- Convert multiple files