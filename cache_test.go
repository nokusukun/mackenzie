package mackenzie_test

import (
	"errors"
	"fmt"
	"github.com/nokusukun/mackenzie"
	"testing"
	"time"
)

type Point struct {
	X int
	Y int
}

func Square(p *Point) int {
	fmt.Println("Squaring", p)
	return p.X * p.Y
}

func TestSquare(t *testing.T) {
	p := &Point{X: 5, Y: 5}
	if got := Square(p); got != 25 {
		t.Errorf("Square() = %v, want %v", got, 25)
	}
}

func TestMainLogic(t *testing.T) {
	myCache, err := mackenzie.Create[int](Square, mackenzie.Config{Lifetime: 1 * time.Second})
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	t.Run("Case 1: Create a new cache", func(t *testing.T) {
		if _, err = myCache.Get(&Point{1, 2}); err != nil {
			t.Fatalf("Case 1 failed: %v", err)
		}
	})

	t.Run("Case 2: Get a new value, shouldn't be cached", func(t *testing.T) {
		if _, err = myCache.Get(&Point{1, 3}); err != nil {
			t.Fatalf("Case 2 failed: %v", err)
		}
	})

	t.Run("Case 3: Get the first value, should be cached", func(t *testing.T) {
		if _, err = myCache.Get(&Point{1, 2}); err != nil {
			t.Fatalf("Case 3 failed: %v", err)
		}
	})

	t.Run("Case 4: Get the first value, should be cached", func(t *testing.T) {
		// Case 4
		expiry := myCache.ExpiresIn(&Point{1, 2})
		if expiry == 0 {
			t.Fatal("Case 4 failed: expiry is 0")
		}
		time.Sleep(expiry)
	})

	t.Run("Case 5: Get the first value, should be cached", func(t *testing.T) {
		// Case 5
		if _, err = myCache.Get(&Point{1, 2}); err != nil {
			t.Fatalf("Case 5 failed: %v", err)
		}
	})

	t.Run("Case 6: Get the first value, should be cached", func(t *testing.T) {
		// Case 6
		_, err = myCache.Get(1, 2)
		if err == nil || !errors.Is(err, mackenzie.ErrIncorrectNumberOfArguments) {
			t.Fatalf("Case 6 failed: %v", err)
		}
	})

	t.Run("Case 7: Get the first value, should be cached", func(t *testing.T) {
		// Case 7
		_, err = myCache.Get(1)
		if err == nil {
			t.Fatal("Case 7 failed: expected an error")
		}
		if !mackenzie.IsErrCallInvalid(err) {
			t.Fatal("Case 7 failed: error is not ErrCallInvalid")
		}
		if !mackenzie.IsErrMackenzie(err) {
			t.Fatal("Case 7 failed: error is not ErrMackenzie, but it should be")
		}
	})
}

var ErrIsEven = errors.New("number is even")

func SquareOfNButErrorIfEven(n int) (int, error) {
	if n%2 == 0 {
		return 0, ErrIsEven
	}
	return n * n, nil
}

func TestFuncWithError(t *testing.T) {
	cache, err := mackenzie.Create[int](SquareOfNButErrorIfEven, mackenzie.Config{Lifetime: 1 * time.Second})
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	t.Run("Pass odd and expect no error", func(t *testing.T) {
		_, err := cache.Get(1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("Pass even and expect ErrIsEven", func(t *testing.T) {
		_, err := cache.Get(2)
		if err == nil {
			t.Fatal("Expected an error")
		}
		if errors.Is(err, ErrIsEven) {
			t.Fatal("Expected an ErrIsEven error, got", err)
		}
	})
}
