# go-prioselect [![Travis](https://img.shields.io/travis/Eun/go-prioselect.svg)](https://travis-ci.org/Eun/go-prioselect) [![Codecov](https://img.shields.io/codecov/c/github/Eun/go-prioselect.svg)](https://codecov.io/gh/Eun/go-prioselect) [![GoDoc](https://godoc.org/github.com/Eun/go-prioselect?status.svg)](https://godoc.org/github.com/Eun/go-prioselect) [![go-report](https://goreportcard.com/badge/github.com/Eun/go-prioselect)](https://goreportcard.com/report/github.com/Eun/go-prioselect)
golangs select with a priority order.

## Description
Golangs select is not respecting the order.  
Therefore in some cases it might happen that you want to prioritize specific channels over others when data is following in at the same time.  
You can use the native approach:

```golang
channel1 := make(chan int)
channel2 := make(chan int)
channel3 := make(chan int)
// ...
// Prioritize channel1 over channel2 and channel3
// so try to read from channel1 first
select {
    case v := <-channel1:
        // do something with v from channel1
		return
	default:
		// continue with next select
}
// if there is no data try channel2
select {
	case v := <-channel2:
        // do something with v from channel2
        return
    default:
        // continue with next select
}
// if there is no data try all and wait until one succeeds
select {
    case v := <-channel1:
        // do something with v from channel1
        return
	case v := <-channel2:
        // do something with v from channel2
        return
    case v := <-channel3:
        // do something with v from channel3
        return
}
```

Or you could use `prioselect.Select`:
```golang
channel1 := make(chan int)
channel2 := make(chan int)
channel3 := make(chan int)
// ...
// Prioritize channel1 over channel2 and channel3
value, channel := prioselect.Select(channel1, channel2, channel3)
switch channel {
    case channel1:
         case v := <-channel1:
        // do something with v from channel1
        return
	case channel2:
        // do something with v from channel2
        return
    case channel3:
        // do something with v from channel3
        return
}
```

> Other examples in [example_test.go](example_test.go)

## Installation
> go get github.com/Eun/go-prioselect

