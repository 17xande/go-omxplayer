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
omx := NewOMXPlayer("file.mp4", false, false, nil)

waiting := make(chan string)
defer close(waiting)

// Returned errors are ignored for brevity.
_ := omx.Open(waiting)
_ := omx.Pause()
_ := omx.Play()
_ := omx.Quit()

// Wait for OMXPlayer process to end.
_ := <-waiting

```

References
----------
This module was originally published as a package of [pi-player](https://github.com/17xande/pi-player/tree/master/pkg/piplayer).  
This standalone version has been re-written, and some parts have been inspired by another great [OMXPlayer library](https://github.com/jleight/omxplayer)

License
-------
The go-omxplayer module is available under the MIT license.
