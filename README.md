## Seneca 

Creates animated GIFs from videos.

<img src="logo.png" width="289" height="309" alt="seneca animated gif logo"/>

## Dependencies

* [Go](http://golang.org/)
* [ffmpeg](http://www.ffmpeg.org/)
* [pipe](http://labix.org/pipe)

## Usage

TBD

## Installation

TBD

## Sample

```bash
$ seneca -video-infile=./goproplane.mp4 -scale 280:_
         -fps 18 -from 00:00:39 -length 9s -speed=slower
```
![animated gif](http://i.imgur.com/4VdXgx3.gif)

## License

* Code is released under Apache license. See [LICENSE][license] file.
* The license for code under the `vendor` subdirectory remains under the purview of their respective owners.
* The [logo](http://commons.wikimedia.org/wiki/File:Nuremberg_chronicles_f_105r_1.png) above is from the public domain.


[license]: https://github.com/javouhey/seneca/blob/master/LICENSE
