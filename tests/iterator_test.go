package tests

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/yu31/gostructs/avl"
	"github.com/yu31/gostructs/bs"
	"github.com/yu31/gostructs/container"
	"github.com/yu31/gostructs/rb"
	"github.com/yu31/gostructs/skip"
)

func shuffleSeeds(s1 []container.Int64) []container.Int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s2 := make([]container.Int64, len(s1))
	for i := 0; i < len(s1); i++ {
		s2[i] = s1[i]
	}
	for i := len(s2) - 1; i > 0; i-- {
		num := r.Intn(i + 1)
		s2[i], s2[num] = s2[num], s2[i]
	}
	return s2
}

func TestContainer_Iterator(t *testing.T) {
	// --------- sequence in box: [22, 24, 35, 61, 64, 67, 76, 84, 87, 91, 97, 130, 133, 145, 150] ---------
	seeds := []container.Int64{22, 24, 35, 61, 64, 67, 76, 84, 87, 91, 97, 130, 133, 145, 150}

	process := func(t *testing.T, box container.Container) {
		for _, k := range shuffleSeeds(seeds) {
			box.Insert(k, int64(k*2+1))
		}

		var iter container.Iterator
		var element container.Element

		// Test case: start == nil and boundary == nil
		t.Run("case1", func(t *testing.T) {
			iter = box.Iter(nil, nil)
			require.NotNil(t, iter)
			require.True(t, iter.Valid())
			for i := 0; i < len(seeds); i++ {
				element = iter.Next()
				require.NotNil(t, element, "key: %v", seeds[i])
				require.Equal(t, element.Key(), seeds[i])
				require.Equal(t, element.Value(), int64(seeds[i]*2+1))
			}
			element = iter.Next()
			require.Nil(t, element)
			require.False(t, iter.Valid())
		})

		// Test case: start != nil && boundary == nil
		t.Run("case2", func(t *testing.T) {
			// start < first node
			iter = box.Iter(container.Int64(21), nil)
			require.NotNil(t, iter)
			require.True(t, iter.Valid())
			for i := 0; i < len(seeds); i++ {
				element = iter.Next()
				require.NotNil(t, element, "key: %v", seeds[i])
				require.Equal(t, element.Key(), seeds[i])
				require.Equal(t, element.Value(), int64(seeds[i]*2+1))
			}
			element = iter.Next()
			require.Nil(t, element)
			require.False(t, iter.Valid())

			// start == first node
			iter = box.Iter(container.Int64(22), nil)
			require.NotNil(t, iter)
			require.True(t, iter.Valid())
			for i := 0; i < len(seeds); i++ {
				element = iter.Next()
				require.NotNil(t, element, "key: %v", seeds[i])
				require.Equal(t, element.Key(), seeds[i])
				require.Equal(t, element.Value(), int64(seeds[i]*2+1))
			}
			element = iter.Next()
			require.Nil(t, element)
			require.False(t, iter.Valid())

			// start > first node && start < last node
			iter = box.Iter(container.Int64(27), nil)
			require.NotNil(t, iter)
			require.True(t, iter.Valid())
			for i := 2; i < len(seeds); i++ {
				element = iter.Next()
				require.NotNil(t, element, "key: %v", seeds[i])
				require.Equal(t, element.Key(), seeds[i])
				require.Equal(t, element.Value(), int64(seeds[i]*2+1))
			}
			element = iter.Next()
			require.Nil(t, element)
			require.False(t, iter.Valid())

			// start > first node && start < last node
			iter = box.Iter(container.Int64(62), nil)
			require.NotNil(t, iter)
			require.True(t, iter.Valid())
			for i := 4; i < len(seeds); i++ {
				element = iter.Next()
				require.NotNil(t, element, "key: %v", seeds[i])
				require.Equal(t, element.Key(), seeds[i])
				require.Equal(t, element.Value(), int64(seeds[i]*2+1))
			}
			element = iter.Next()
			require.Nil(t, element)
			require.False(t, iter.Valid())

			// start > root node && start < last node
			iter = box.Iter(container.Int64(132), nil)
			require.NotNil(t, iter)
			require.True(t, iter.Valid())
			for i := 12; i < len(seeds); i++ {
				element = iter.Next()
				require.NotNil(t, element, "key: %v", seeds[i])
				require.Equal(t, element.Key(), seeds[i])
				require.Equal(t, element.Value(), int64(seeds[i]*2+1))
			}
			element = iter.Next()
			require.Nil(t, element)
			require.False(t, iter.Valid())

			// start == last node
			iter = box.Iter(container.Int64(150), nil)
			require.NotNil(t, iter)
			require.True(t, iter.Valid())
			element = iter.Next()
			require.Equal(t, element.Key(), container.Int64(150))
			element = iter.Next()
			require.Nil(t, element)
			require.False(t, iter.Valid())

			// start > last node
			iter = box.Iter(container.Int64(156), nil)
			require.NotNil(t, iter)
			element = iter.Next()
			require.Nil(t, element)
			require.False(t, iter.Valid())
		})

		// Test case: start == nil && boundary != nil
		t.Run("case3", func(t *testing.T) {
			// boundary < first node
			iter = box.Iter(nil, container.Int64(21))
			require.NotNil(t, iter)
			require.False(t, iter.Valid())
			element = iter.Next()
			require.Nil(t, element)

			// boundary == first node
			iter = box.Iter(nil, container.Int64(22))
			require.NotNil(t, iter)
			require.False(t, iter.Valid())
			element = iter.Next()
			require.Nil(t, element)

			// boundary > first node
			iter = box.Iter(nil, container.Int64(24))
			require.NotNil(t, iter)
			require.True(t, iter.Valid())
			element = iter.Next()
			require.Equal(t, element.Key(), container.Int64(22))
			element = iter.Next()
			require.Nil(t, element)
			require.False(t, iter.Valid())

			// boundary < last node && bound > first node
			iter = box.Iter(nil, container.Int64(147))
			require.NotNil(t, iter)
			require.True(t, iter.Valid())
			for i := 0; i < len(seeds)-1; i++ {
				element = iter.Next()
				require.NotNil(t, element, "key: %v", seeds[i])
				require.Equal(t, element.Key(), seeds[i])
				require.Equal(t, element.Value(), int64(seeds[i]*2+1))
			}
			element = iter.Next()
			require.Nil(t, element)
			require.False(t, iter.Valid())

			// boundary == last node
			iter = box.Iter(nil, container.Int64(150))
			require.NotNil(t, iter)
			require.True(t, iter.Valid())
			for i := 0; i < len(seeds)-1; i++ {
				element = iter.Next()
				require.NotNil(t, element, "key: %v", seeds[i])
				require.Equal(t, element.Key(), seeds[i])
				require.Equal(t, element.Value(), int64(seeds[i]*2+1))
			}
			element = iter.Next()
			require.Nil(t, element)
			require.False(t, iter.Valid())

			// boundary > last node
			iter = box.Iter(nil, container.Int64(156))
			require.NotNil(t, iter)
			require.True(t, iter.Valid())
			for i := range seeds {
				element = iter.Next()
				require.NotNil(t, element, "key: %v", seeds[i])
				require.Equal(t, element.Key(), seeds[i])
				require.Equal(t, element.Value(), int64(seeds[i]*2+1))
			}
			element = iter.Next()
			require.Nil(t, element)
			require.False(t, iter.Valid())
		})

		// Test case: start != nil && boundary != nil
		t.Run("case4", func(t *testing.T) {
			// start < boundary && start > first node && bound < last node
			iter = box.Iter(container.Int64(68), container.Int64(132))
			require.NotNil(t, iter)
			require.True(t, iter.Valid())
			for i := 6; i < len(seeds)-3; i++ {
				element = iter.Next()
				require.NotNil(t, element, "key: %v", seeds[i])
				require.Equal(t, element.Key(), seeds[i])
				require.Equal(t, element.Value(), int64(seeds[i]*2+1))
			}
			element = iter.Next()
			require.Nil(t, element)
			require.False(t, iter.Valid())

			// start < boundary && start < first node && bound > last node
			iter = box.Iter(container.Int64(21), container.Int64(153))
			require.NotNil(t, iter)
			require.True(t, iter.Valid())
			for i := 0; i < len(seeds); i++ {
				element = iter.Next()
				require.NotNil(t, element, "key: %v", seeds[i])
				require.Equal(t, element.Key(), seeds[i])
				require.Equal(t, element.Value(), int64(seeds[i]*2+1))
			}
			element = iter.Next()
			require.Nil(t, element)
			require.False(t, iter.Valid())

			// start == boundary, start and boundary exists.
			iter = box.Iter(container.Int64(24), container.Int64(24))
			require.NotNil(t, iter)
			element = iter.Next()
			require.Nil(t, element)
			require.False(t, iter.Valid())

			// start == boundary, start and boundary not exists.
			iter = box.Iter(container.Int64(25), container.Int64(25))
			require.NotNil(t, iter)
			element = iter.Next()
			require.Nil(t, element)
			require.False(t, iter.Valid())

			// start < boundary && start < first node && bound < first node
			iter = box.Iter(container.Int64(21), container.Int64(13))
			require.NotNil(t, iter)
			element = iter.Next()
			require.Nil(t, element)
			require.False(t, iter.Valid())

			// start > boundary && start > first node
			iter = box.Iter(container.Int64(65), container.Int64(27))
			require.NotNil(t, iter)
			element = iter.Next()
			require.Nil(t, element)
			require.False(t, iter.Valid())
		})
	}

	t.Run("bstree", func(t *testing.T) {
		process(t, bs.New())
	})
	t.Run("avtree", func(t *testing.T) {
		process(t, avl.New())
	})
	t.Run("rbtree", func(t *testing.T) {
		process(t, rb.New())
	})
	t.Run("skiplist", func(t *testing.T) {
		process(t, skip.New())
	})
}