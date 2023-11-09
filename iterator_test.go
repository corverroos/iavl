package iavl_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/iavl/v2"
	"github.com/stretchr/testify/require"
)

func Test_Iterator(t *testing.T) {
	pool := iavl.NewNodePool()
	sql, err := iavl.NewInMemorySqliteDb(pool)
	require.NoError(t, err)

	tree := iavl.NewTree(sql, pool, iavl.TreeOptions{StateStorage: true})
	set := func(key string, value string) {
		_, err := tree.Set([]byte(key), []byte(value))
		require.NoError(t, err)
	}
	set("a", "1")
	set("b", "2")
	set("c", "3")
	set("d", "4")
	set("e", "5")
	set("f", "6")
	set("g", "7")

	cases := []struct {
		name          string
		start, end    []byte
		inclusive     bool
		ascending     bool
		expectedCount int
		expectedStart []byte
		expectedEnd   []byte
	}{
		{
			name:          "all",
			start:         nil,
			end:           nil,
			ascending:     true,
			expectedCount: 7,
			expectedStart: []byte("a"),
			expectedEnd:   []byte("g"),
		},
		{
			name:          "b start",
			start:         []byte("b"),
			end:           nil,
			ascending:     true,
			expectedCount: 6,
			expectedStart: []byte("b"),
			expectedEnd:   []byte("g"),
		},
		{
			name:          "ab start",
			start:         []byte("ab"),
			end:           nil,
			ascending:     true,
			expectedCount: 6,
			expectedStart: []byte("b"),
			expectedEnd:   []byte("g"),
		},
		{
			name:          "c end inclusive",
			start:         nil,
			end:           []byte("c"),
			ascending:     true,
			inclusive:     true,
			expectedCount: 3,
			expectedStart: []byte("a"),
			expectedEnd:   []byte("c"),
		},
		{
			name:          "d end exclusive",
			start:         nil,
			end:           []byte("d"),
			ascending:     true,
			inclusive:     false,
			expectedCount: 3,
			expectedStart: []byte("a"),
			expectedEnd:   []byte("c"),
		},
		{
			name:          "ce end inclusive",
			start:         nil,
			end:           []byte("c"),
			ascending:     true,
			inclusive:     true,
			expectedCount: 3,
			expectedStart: []byte("a"),
			expectedEnd:   []byte("c"),
		},
		{
			name:          "ce end exclusive",
			start:         nil,
			end:           []byte("ce"),
			ascending:     true,
			inclusive:     false,
			expectedCount: 3,
			expectedStart: []byte("a"),
			expectedEnd:   []byte("c"),
		},
		{
			name:          "b to e",
			start:         []byte("b"),
			end:           []byte("e"),
			inclusive:     true,
			ascending:     true,
			expectedCount: 4,
			expectedStart: []byte("b"),
			expectedEnd:   []byte("e"),
		},
		{
			name:          "all desc",
			start:         nil,
			end:           nil,
			ascending:     false,
			expectedCount: 7,
			expectedStart: []byte("g"),
			expectedEnd:   []byte("a"),
		},
		{
			name:          "f start desc",
			start:         []byte("f"),
			end:           nil,
			ascending:     false,
			expectedCount: 6,
			expectedStart: []byte("f"),
			expectedEnd:   []byte("a"),
		},
		{
			name:          "fe start desc",
			start:         []byte("fe"),
			end:           nil,
			ascending:     false,
			expectedCount: 6,
			expectedStart: []byte("f"),
			expectedEnd:   []byte("a"),
		},
		{
			name:          "c stop desc",
			start:         nil,
			end:           []byte("c"),
			ascending:     false,
			expectedCount: 4,
			expectedStart: []byte("g"),
			expectedEnd:   []byte("d"),
		},
		{
			name:          "ce stop desc",
			start:         nil,
			end:           []byte("ce"),
			ascending:     false,
			expectedCount: 4,
			expectedStart: []byte("g"),
			expectedEnd:   []byte("d"),
		},
		{
			name:          "f to c desc",
			start:         []byte("f"),
			end:           []byte("c"),
			ascending:     false,
			expectedCount: 3,
			expectedStart: []byte("f"),
			expectedEnd:   []byte("d"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				itr *iavl.Iterator
				err error
			)
			if tc.ascending {
				itr, err = tree.Iterator(tc.start, tc.end, tc.inclusive)
			} else {
				itr, err = tree.ReverseIterator(tc.start, tc.end)
			}
			require.NoError(t, err)

			cnt := 0
			for ; itr.Valid(); itr.Next() {
				if cnt == 0 {
					require.Equal(t, tc.expectedStart, itr.Key())
				}
				fmt.Printf("%s %s\n", itr.Key(), itr.Value())
				require.NoError(t, itr.Error())
				cnt++
			}
			require.Equal(t, tc.expectedCount, cnt)
			require.Equal(t, tc.expectedEnd, itr.Key())
			require.False(t, itr.Valid())
			require.NoError(t, itr.Close())
		})
	}
}
