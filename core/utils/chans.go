package utils

import (
	"log/slog"
	"sync"
)

func MergeWait[T any](cs ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	out := make(chan T)

	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan T) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func MergePump[T any](out chan T, chans ...<-chan T) {
	o := MergeWait(chans...)
	for v := range o {
		slog.Debug("Merging", v)
		out <- v
	}
}

func MergeSliceWait[T any](cs ...<-chan []T) []T {
	rtn := make([]T, 0)
	out := MergeWait(cs...)
	for v := range out {
		rtn = append(rtn, v...)
	}
	return rtn
}

func Jobber[T any](f func() (T, error)) (out chan T, job func() error) {
	out = make(chan T)
	return out, func() error {
		defer close(out)
		v, err := f()
		if err != nil {
			return err
		}
		out <- v
		return nil
	}
}
