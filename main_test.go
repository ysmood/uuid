package uuid_test

import (
	"fmt"
	"runtime"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ysmood/uuid"
)

func Example_basic() {
	id := uuid.New()

	fmt.Println(len(id.Bin()))

	fmt.Println(len(id.Hex()))

	// Output:
	// 16
	// 32
}

func TestBytes(t *testing.T) {
	id := uuid.New()

	expect := `uuid-\d\d\d\d_\d\d_\d\dT\d\d:\d\d:\d\d-[0-9a-f]{4}-[0-9a-f]{6}`
	assert.Regexp(t, expect, id.String())

	prefix := uuid.Parse(id.Bin()).Time

	assert.Len(t, id.Bin(), 16)
	assert.True(t, time.Since(prefix) < time.Millisecond)

	assert.Panics(t, func() {
		id := uuid.New()
		id.Namespace = make([]byte, 20)
		id.Bin()
	})

	assert.Panics(t, func() {
		id := uuid.New()
		id.Machine = make([]byte, 20)
		id.Bin()
	})

	assert.Panics(t, func() {
		id := uuid.New()
		id.Noise = make([]byte, 20)
		id.Bin()
	})
}

func TestOrder(t *testing.T) {
	list := []string{}
	sorted := []string{}
	for range make([]struct{}, 100) {
		time.Sleep(time.Millisecond)
		id := uuid.New().Hex()
		list = append(list, id)
		sorted = append(sorted, id)
	}

	sort.Strings(sorted)

	assert.Equal(t, list, sorted)
}

func TestCollision(t *testing.T) {
	for range make([]struct{}, 10) {
		dict := sync.Map{}
		wg := sync.WaitGroup{}
		wg.Add(runtime.NumCPU())

		for range make([]struct{}, runtime.NumCPU()) {
			go func() {
				for range make([]struct{}, 100000) {
					id := uuid.New().Hex()

					if _, has := dict.Load(id); has {
						panic("collision")
					}

					dict.Store(id, struct{}{})
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
