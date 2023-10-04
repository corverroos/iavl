package iavl

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"testing"

	"github.com/cosmos/iavl-bench/bench"
	"github.com/cosmos/iavl/v2/testutil"
	"github.com/stretchr/testify/require"
)

func TestMinRightToken(t *testing.T) {
	cases := []struct {
		name string
		a    string
		b    string
		want string
	}{
		{name: "naive", a: "alphabet", b: "elephant", want: "e"},
		{name: "trivial substring", a: "bird", b: "bingo", want: "bir"},
		{name: "longer", a: "bird", b: "birdy", want: "birdy"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got1 := MinRightToken([]byte(tc.a), []byte(tc.b))
			got2 := MinRightToken([]byte(tc.b), []byte(tc.a))
			if string(got1) != tc.want {
				t.Errorf("MinRightToken(%q, %q) = %q, want %q", tc.a, tc.b, got1, tc.want)
			}
			require.Equal(t, got1, got2)
		})
	}
}

func TestMinRightToken_Gen(t *testing.T) {
	seed := rand.Int()
	t.Logf("seed: %d", seed)
	r := rand.New(rand.NewSource(int64(seed)))
	for i := 0; i < 1_000_000; i++ {
		lenA := r.Intn(20)
		lenB := r.Intn(20)
		a := make([]byte, lenA)
		b := make([]byte, lenB)
		r.Read(a)
		r.Read(b)
		res := MinRightToken(a, b)
		resSwap := MinRightToken(b, a)
		require.Equal(t, res, resSwap)

		compare := bytes.Compare(a, b)
		switch {
		case compare < 0:
			require.Less(t, a, res)
			require.GreaterOrEqual(t, b, res)
		case compare > 0:
			require.Less(t, b, res)
			require.GreaterOrEqual(t, a, res)
		default:
			require.Equal(t, a, res)
			require.Equal(t, b, res)
		}
	}
}

func TestTreeSanity(t *testing.T) {
	outDir := "/tmp/tree_viz"
	require.NoError(t, os.RemoveAll(outDir))
	require.NoError(t, os.Mkdir(outDir, 0755))

	//gen := bench.ChangesetGenerator{
	//	Seed:             77,
	//	KeyMean:          4,
	//	KeyStdDev:        1,
	//	ValueMean:        50,
	//	ValueStdDev:      15,
	//	InitialSize:      1000,
	//	FinalSize:        10000,
	//	Versions:         5,
	//	ChangePerVersion: 10,
	//	DeleteFraction:   0.2,
	//}
	opts := testutil.NewTreeBuildOptions()
	//itr, err := gen.Iterator()
	var err error
	itr := opts.Iterator
	require.NoError(t, err)
	tree := NewTree(nil, NewNodePool())
	for ; itr.Valid(); err = itr.Next() {
		require.NoError(t, err)
		nodes := itr.Nodes()
		for ; nodes.Valid(); err = nodes.Next() {
			require.NoError(t, err)
			node := nodes.GetNode()
			if node.Delete {
				_, _, err := tree.Remove(node.Key)
				require.NoError(t, err)
			} else {
				_, err := tree.Set(node.Key, node.Value)
				require.NoError(t, err)
			}
		}
		switch itr.Version() {
		case 1:
			rehashTree(tree.root)
			require.Equal(t, "48c3113b8ba523d3d539d8aea6fce28814e5688340ba7334935c1248b6c11c7a", hex.EncodeToString(tree.root.hash))
			fmt.Printf("version=%d, hash=%x size=%d\n", itr.Version(), tree.root.hash, tree.root.size)
		case 150:
			rehashTree(tree.root)
			require.Equal(t, "876e9b511761011c273b641d4a43f510568760203a43b07d5cc3ff7b9eb8dbfb", hex.EncodeToString(tree.root.hash))
			fmt.Printf("version=%d, hash=%x size=%d\n", itr.Version(), tree.root.hash, tree.root.size)
			return
		}

		//f, err := os.Create(fmt.Sprintf("%s/version%d.dot", outDir, itr.Version()))
		//require.NoError(t, err)
		//g := writeDotGraph(tree.root, &dot.Graph{})
		//_, err = f.Write([]byte(g.String()))
		//require.NoError(t, err)
	}
}

func TestTokenizedTree(t *testing.T) {
	// can a total order of (sortKey, height) be built to satisfy a traversal order of the tree?
	// in-order seems the most possible.

	outDir := "/tmp/tree_viz"
	require.NoError(t, os.RemoveAll(outDir))
	require.NoError(t, os.Mkdir(outDir, 0755))

	var inOrder func(node *Node) []*Node
	inOrder = func(node *Node) (nodes []*Node) {
		if node == nil {
			return nil
		}
		nodes = append(nodes, inOrder(node.leftNode)...)
		nodes = append(nodes, node)
		nodes = append(nodes, inOrder(node.rightNode)...)
		return nodes
	}

	var preOrder func(node *Node) []*Node
	preOrder = func(node *Node) (nodes []*Node) {
		if node == nil {
			return nil
		}

		nodes = append(nodes, node)
		nodes = append(nodes, preOrder(node.leftNode)...)
		nodes = append(nodes, preOrder(node.rightNode)...)
		return nodes
	}

	gen := bench.ChangesetGenerator{
		Seed:             77,
		KeyMean:          4,
		KeyStdDev:        1,
		ValueMean:        10,
		ValueStdDev:      1,
		InitialSize:      1000,
		FinalSize:        10000,
		Versions:         20,
		ChangePerVersion: 10,
		DeleteFraction:   0.2,
	}
	itr, err := gen.Iterator()
	require.NoError(t, err)
	tree := NewTree(nil, NewNodePool())
	//tree.emitDotGraphs = true

	step := 0
	for ; itr.Valid(); err = itr.Next() {
		if itr.Version() > 1 {
			break
		}
		require.NoError(t, err)
		nodes := itr.Nodes()
		for ; nodes.Valid(); err = nodes.Next() {
			//if i > 7 {
			//	break
			//}

			require.NoError(t, err)
			node := nodes.GetNode()
			strKey := hex.EncodeToString(node.Key)
			bzKey := []byte(strKey)

			step++
			if node.Delete {
				_, _, err := tree.Remove(bzKey)
				require.NoError(t, err)
			} else {
				_, err := tree.Set(bzKey, node.Value)
				require.NoError(t, err)
			}
			for j, dg := range tree.dotGraphs {
				f, err := os.Create(fmt.Sprintf("%s/step_%04d_%d.dot", outDir, step, j))
				require.NoError(t, err)
				_, err = f.Write([]byte(dg.String()))
				require.NoError(t, err)
				require.NoError(t, f.Close())
			}
			tree.dotGraphs = nil

			// verify the tree at every step

			orderedNodes := preOrder(tree.root)
			sort.Slice(orderedNodes, func(i, j int) bool {
				a := orderedNodes[i]
				b := orderedNodes[j]
				res := bytes.Compare(a.sortKey, b.sortKey)
				// order by (key ASC, height DESC)
				if res != 0 {
					return res < 0
				} else {
					//fmt.Printf("resolve collision sortKey=%s\n", orderedNodes[j].sortKey)
					// height DESC
					// return a.subtreeHeight > b.subtreeHeight

					// or, more specifically below.  sortKey collisions may only occur between leaf and branch nodes.
					// in this case choose the leaf node first for in-order traversal.
					switch {
					case a.isLeaf() && b.isLeaf():
						panic("invariant violated: two leaves with same sortKey")
					case a.isLeaf() && !b.isLeaf():
						return false
					case !a.isLeaf() && b.isLeaf():
						return true
					default:
						panic("invariant violated: two branches with same sortKey")
					}
				}
			})

			inOrderNodes := inOrder(tree.root)
			var lastNode *Node
			for i, n := range inOrderNodes {
				//fmt.Printf("node: %s, %s, %d\n", string(n.key), string(n.sortKey), n.subtreeHeight)
				if lastNode != nil {
					require.LessOrEqual(t, lastNode.key, n.key)
					if bytes.Equal(lastNode.key, n.key) {
						// in-order assertion
						require.Greater(t, lastNode.subtreeHeight, n.subtreeHeight)
					}
					require.Equalf(t, n.key, orderedNodes[i].key, "expected (%s, %d), got (%s, %d)",
						string(n.key), n.subtreeHeight, string(orderedNodes[i].key), orderedNodes[i].subtreeHeight)
					require.Equalf(t, n.subtreeHeight, orderedNodes[i].subtreeHeight,
						"heights don't match, node.key: %s step=%d", n.key, step)
					require.Equal(t, n, orderedNodes[i])
				}
				lastNode = n
			}
		}
	}
}
