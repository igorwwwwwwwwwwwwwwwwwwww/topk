package main

import (
	"bufio"
	"cmp"
	"container/heap"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

var k = flag.Int("k", 10, "limit to this many top values")
var other = flag.Bool("other", false, "include sum count of remaining values")

type KVPair struct {
	Item  string
	Count int
}

type MaxHeap []KVPair

func (h MaxHeap) Len() int { return len(h) }
func (h MaxHeap) Less(i, j int) bool {
	return cmp.Or(
		cmp.Compare(h[i].Count, h[j].Count),
		cmp.Compare(h[i].Item, h[j].Item),
	) < 0
}
func (h MaxHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x any) {
	*h = append(*h, x.(KVPair))
}

func (h *MaxHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	flag.Parse()

	reader := os.Stdin

	if flag.NArg() > 0 {
		filename := flag.Arg(0)
		f, err := os.Open(filename)
		if err != nil {
			log.Fatalf("failed to open file: %s", err)
		}
		defer f.Close()

		reader = f
	}

	items := make(map[string]int)
	var total int

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		items[line]++
		total++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}

	h := &MaxHeap{}
	heap.Init(h)

	for key, val := range items {
		heap.Push(h, KVPair{key, val})
		if h.Len() > *k {
			heap.Pop(h)
		}
	}

	if len(*h) == 0 {
		return
	}

	sort.Sort(sort.Reverse(h))

	maxLen := 0
	maxCount := 0 // technically just (*h)[0].Count
	topKTotal := 0
	for _, pair := range *h {
		maxLen = max(maxLen, len(pair.Item))
		maxCount = max(maxCount, pair.Count)
		topKTotal += pair.Count
	}

	if *other {
		pair := KVPair{"OTHER", total - topKTotal}
		*h = append(*h, pair)
		maxLen = max(maxLen, len(pair.Item))
		maxCount = max(maxCount, pair.Count)
	}

	for _, pair := range *h {
		barWidth := pair.Count * 50 / maxCount
		bar := strings.Repeat("∎", barWidth)
		fmt.Printf("%-*s  %6d  %s\n", maxLen, pair.Item, pair.Count, bar)
	}
}
