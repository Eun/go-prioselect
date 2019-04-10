package prioselect

import "reflect"

// Select does a select on the specified channels in their order
// it returns the value and the channel the value was read from
//
// v, _ := Select(time.After(time.Minute), time.After(time.Second))
// fmt.Printf("Timer ticked after %s", v.(time.Time).String())
func Select(channels ...interface{}) (value interface{}, channel interface{}) {
	size := len(channels)
	if size <= 0 {
		return nil, nil
	}
	channelsToRead := make([]reflect.SelectCase, size)

	for i, c := range channels {
		channelsToRead[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(c),
		}
	}

	defaultCase := reflect.SelectCase{
		Dir: reflect.SelectDefault,
	}

	for i := 0; i < size-1; i++ {
		chosen, v, ok := reflect.Select([]reflect.SelectCase{channelsToRead[i], defaultCase})
		if ok {
			if chosen == 0 {
				return v.Interface(), channels[i]
			}
		}
	}

	// if this is the last channel
	for {
		chosen, v, ok := reflect.Select(channelsToRead)
		if ok {
			return v.Interface(), channelsToRead[chosen].Chan.Interface()
		}
		// remove the selected channel from the slice
		channelsToRead = append(channelsToRead[:chosen], channelsToRead[chosen+1:]...)
		if len(channelsToRead) <= 0 {
			return nil, nil
		}
	}
}
