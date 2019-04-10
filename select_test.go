package prioselect

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSelect(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		ch1 := make(chan int)
		ch2 := make(chan int)

		go func() {
			time.Sleep(time.Millisecond * 10)
			ch1 <- 1
		}()

		value, channel := Select(ch1, ch2)
		require.Equal(t, 1, value)
		require.Equal(t, channel, ch1)
	})

	t.Run("SimpleReverse", func(t *testing.T) {
		ch1 := make(chan int)
		ch2 := make(chan int)

		go func() {
			time.Sleep(time.Millisecond * 10)
			ch2 <- 2
		}()

		value, channel := Select(ch1, ch2)
		require.Equal(t, 2, value)
		require.Equal(t, channel, ch2)
	})

	t.Run("Close", func(t *testing.T) {
		ch1 := make(chan int)
		ch2 := make(chan int)

		go func() {
			time.Sleep(time.Millisecond * 10)
			close(ch1)
			time.Sleep(time.Millisecond * 10)
			ch2 <- 2
		}()

		value, channel := Select(ch1, ch2)
		require.Equal(t, 2, value)
		require.Equal(t, channel, ch2)
	})
	t.Run("Close All", func(t *testing.T) {
		ch1 := make(chan int)
		ch2 := make(chan int)

		go func() {
			time.Sleep(time.Millisecond * 10)
			close(ch1)
			close(ch2)
		}()

		value, channel := Select(ch1, ch2)
		require.Nil(t, value)
		require.Nil(t, channel)
	})

	t.Run("Parallel", func(t *testing.T) {
		t.Run("", func(t *testing.T) {
			ch1 := make(chan int, 1)
			ch2 := make(chan int, 1)

			ch1 <- 1
			ch2 <- 2

			value, channel := Select(ch1, ch2)
			require.Equal(t, 1, value)
			require.Equal(t, channel, ch1)
		})
		t.Run("", func(t *testing.T) {
			ch1 := make(chan int, 1)
			ch2 := make(chan int, 1)

			ch2 <- 2
			ch1 <- 1

			value, channel := Select(ch1, ch2)
			require.Equal(t, 1, value)
			require.Equal(t, channel, ch1)
		})
	})

	t.Run("No Parameters", func(t *testing.T) {
		value, channel := Select()
		require.Nil(t, value)
		require.Nil(t, channel)
	})

	t.Run("Parameter is not a channel", func(t *testing.T) {
		require.Panics(t, func() {
			value, channel := Select(1)
			require.Nil(t, value)
			require.Nil(t, channel)
		})
	})

	t.Run("Types", func(t *testing.T) {
	})
}

func traditional(a chan int, b chan int) (int, chan int) {
	select {
	case v := <-a:
		return v, a
	default:
		// continue with next select
	}
	select {
	case v := <-a:
		return v, a
	case v := <-b:
		return v, b
	}
}

func TestTraditional(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		ch1 := make(chan int)
		ch2 := make(chan int)

		go func() {
			time.Sleep(time.Millisecond * 10)
			ch1 <- 1
		}()

		value, channel := traditional(ch1, ch2)
		require.Equal(t, 1, value)
		require.Equal(t, channel, ch1)
	})

	t.Run("SimpleReverse", func(t *testing.T) {
		ch1 := make(chan int)
		ch2 := make(chan int)

		go func() {
			time.Sleep(time.Millisecond * 10)
			ch2 <- 2
		}()

		value, channel := traditional(ch1, ch2)
		require.Equal(t, 2, value)
		require.Equal(t, channel, ch2)
	})

	t.Run("Parallel", func(t *testing.T) {
		t.Run("", func(t *testing.T) {
			ch1 := make(chan int, 1)
			ch2 := make(chan int, 1)

			ch1 <- 1
			ch2 <- 2

			value, channel := traditional(ch1, ch2)
			require.Equal(t, 1, value)
			require.Equal(t, channel, ch1)
		})
		t.Run("", func(t *testing.T) {
			ch1 := make(chan int, 1)
			ch2 := make(chan int, 1)

			ch2 <- 2
			ch1 <- 1

			value, channel := traditional(ch1, ch2)
			require.Equal(t, 1, value)
			require.Equal(t, channel, ch1)
		})
	})
}

func BenchmarkSelect(b *testing.B) {
	b.Run("Simple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ch1 := make(chan int)
			ch2 := make(chan int)

			go func() {
				time.Sleep(time.Millisecond * 10)
				ch1 <- 1
			}()

			value, channel := Select(ch1, ch2)
			require.Equal(b, 1, value)
			require.Equal(b, channel, ch1)
		}
	})

	b.Run("Parallel", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ch1 := make(chan int, 1)
			ch2 := make(chan int, 1)

			ch1 <- 1
			ch2 <- 2

			value, channel := Select(ch1, ch2)
			require.Equal(b, 1, value)
			require.Equal(b, channel, ch1)
		}
	})
}

func BenchmarkTraditional(b *testing.B) {
	b.Run("Simple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ch1 := make(chan int)
			ch2 := make(chan int)

			go func() {
				time.Sleep(time.Millisecond * 10)
				ch1 <- 1
			}()

			value, channel := traditional(ch1, ch2)
			require.Equal(b, 1, value)
			require.Equal(b, channel, ch1)
		}
	})

	b.Run("Parallel", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ch1 := make(chan int, 1)
			ch2 := make(chan int, 1)

			ch1 <- 1
			ch2 <- 2

			value, channel := traditional(ch1, ch2)
			require.Equal(b, 1, value)
			require.Equal(b, channel, ch1)
		}
	})
}
