go-omxplayer
============
> A basic module for controlling
> [OMXPlayer](https://github.com/popcornmix/omxplayer) from
> [Go](https://golang.org/)

Table of Contents
-----------------
- [Example](#example)
- [References](#references)
- [Licence](#licence)

Example
-------
```go
omx := NewOMXPlayer()
omx.Open("file.mp4")
omx.Pause()
omx.Play()
omx.Quit()
```

References
----------
This module was originally published as a package of [pi-player](https://github.com/17xande/pi-player/tree/master/pkg/piplayer).  
This standalone version has been re-written, and some parts have been inspired by another great [OMXPlayer library](https://github.com/jleight/omxplayer)

License
-------
The go-omxplayer module is available under the MIT license.
