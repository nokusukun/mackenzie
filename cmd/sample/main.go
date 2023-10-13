package main

import (
	"errors"
	"fmt"
	"github.com/nokusukun/bingo/mackenzie"
	"log/slog"
	"os"
	"time"
)

func init() {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)
	slog.SetDefault(
		slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: lvl,
		})),
	)

}

type Point struct {
	X int
	Y int
}

func main() {

	myCache, err := mackenzie.Create[int](func(p *Point) int {
		fmt.Println("Squaring", p)
		return p.X * p.Y
	}, mackenzie.Config{Lifetime: 1 * time.Second})

	if err != nil {
		panic(err)
	}

	// Create a new value
	val, err := myCache.Get(&Point{1, 2})
	if err != nil {
		panic(err)
	}
	fmt.Println("OK!", val)

	// Get a new value, shouldn't be cached
	val, err = myCache.Get(&Point{1, 3})
	if err != nil {
		panic(err)
	}
	fmt.Println("OK!", val)

	// Get the first value, should be cached
	val, err = myCache.Get(&Point{1, 2})
	if err != nil {
		panic(err)
	}
	fmt.Println("OK!", val)

	// Get the first value's expiry, make sure that it isn't 0 and sleep for that duration
	expiry := myCache.ExpiresIn(&Point{1, 2})
	fmt.Println("Expiry:", expiry)
	fmt.Println("Sleeping for 2 seconds")
	time.Sleep(expiry)

	// Get the first value, should be expired
	val, err = myCache.Get(&Point{1, 2})
	if err != nil {
		panic(err)
	}
	fmt.Println("OK!", val)

	// Force incorrect number of arguments by passing wrong number of arguments
	val, err = myCache.Get(1, 2)
	if err != nil {
		if !errors.Is(err, mackenzie.ErrIncorrectNumberOfArguments) {
			panic(err)
		}
		fmt.Println("OK!", err)
	}
	fmt.Println("OK!", val)

	// Force incorrect type of arguments by passing wrong type of arguments
	val, err = myCache.Get(1)
	if err != nil {
		fmt.Println("OK!", err)
	}
	fmt.Println("OK!", val)
	fmt.Println("Is ErrCallInvalid? (it should be)", mackenzie.IsErrCallInvalid(err))
	fmt.Println("Is ErrMackenzie? (it should not be)", mackenzie.IsErrMackenzie(err))
}
