package main

import (
	"bufio"
	"cmp"
	"container/heap"
	"flag"
	"fmt"
	"hash/maphash"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"sort"
	"strings"
)

var k = flag.Int("k", 10, "limit to this many top values")
var approx = flag.Bool("approx", false, "approximate results via filtered space-saving algorithm (zipf)")
var filterBits = flag.Int("fss-bits", 8, "bits to use for filter hash table")
var spaceSavingFactor = flag.Float64("fss-factor", 5.0, "factor by which to scale k for space-saving")
var enableOther = flag.Bool("other", false, "include sum count of remaining values")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

type Filter struct {
	Alpha map[uint64]int
	C     map[uint64]int
}

type KVPair struct {
	Item  string
	Hash  uint64
	Count int
	Error int
	Index int
}

type MinHeap []*KVPair

func (h MinHeap) Len() int { return len(h) }
func (h MinHeap) Less(i, j int) bool {
	return cmp.Or(
		cmp.Compare(h[i].Count, h[j].Count), // <
		cmp.Compare(h[j].Error, h[i].Count), // >
		cmp.Compare(h[i].Item, h[j].Item),   // >
	) < 0
}
func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i
	h[j].Index = j
}

func (h *MinHeap) Push(x any) {
	n := len(*h)
	item := x.(*KVPair)
	item.Index = n
	*h = append(*h, item)
}

func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	old[n-1] = nil
	x.Index = -1 // for safety
	*h = old[0 : n-1]
	return x
}

func topk(reader io.Reader, k int) (*MinHeap, int, error) {
	scanner := bufio.NewScanner(reader)
	buf := make([]byte, 1024*1024) // 1MB buffer
	scanner.Buffer(buf, 1024*1024)

	items := make(map[string]int)
	var total int

	for scanner.Scan() {
		line := scanner.Text()
		items[line]++
		total++
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, fmt.Errorf("error reading file: %s", err)
	}

	h := &MinHeap{}
	heap.Init(h)

	for key, val := range items {
		heap.Push(h, &KVPair{key, 0, val, 0, -1})
		if h.Len() > k {
			heap.Pop(h)
		}
	}

	sort.Sort(sort.Reverse(h))

	return h, total, nil
}

func filteredSpaceSaving(reader io.Reader, k int) (*MinHeap, int, error) {
	scanner := bufio.NewScanner(reader)
	buf := make([]byte, 1024*1024) // 1MB buffer
	scanner.Buffer(buf, 1024*1024)

	// TODO: fixed seed for consistency?
	seed := maphash.MakeSeed()
	mask := uint64((1 << *filterBits) - 1)
	filter := Filter{
		Alpha: make(map[uint64]int),
		C:     make(map[uint64]int),
	}

	monitored := make(map[string]*KVPair)
	var total int

	h := &MinHeap{}
	heap.Init(h)

	heapK := int(float64(k) * *spaceSavingFactor)
	heapK = max(k, heapK)

	for scanner.Scan() {
		total++

		hash := maphash.Bytes(seed, scanner.Bytes())
		bucket := hash & mask

		if filter.C[bucket] > 0 {
			line := scanner.Text()
			if pair, ok := monitored[line]; ok {
				pair.Count++
				heap.Fix(h, pair.Index)
			}
			continue
		}

		hMin := 0
		if h.Len() > 0 {
			hMin = (*h)[0].Count
		}
		if filter.Alpha[bucket]+1 >= hMin {
			line := scanner.Text()
			if h.Len() >= heapK {
				pair := heap.Pop(h).(*KVPair)
				delete(monitored, pair.Item)
				filter.C[pair.Hash&mask]--
				filter.Alpha[pair.Hash&mask] = pair.Count
			}

			pair := &KVPair{line, hash, filter.Alpha[bucket] + 1, filter.Alpha[bucket], -1}
			heap.Push(h, pair)
			monitored[line] = pair

			filter.C[bucket]++
		} else {
			filter.Alpha[bucket]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, fmt.Errorf("error reading file: %s", err)
	}

	// TODO: consider reverse pop instead
	sort.Sort(sort.Reverse(h))
	*h = (*h)[:k]

	return h, total, nil
}

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			pprof.WriteHeapProfile(f)
			f.Close()
		}()
	}

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

	algo := topk
	if *approx {
		algo = filteredSpaceSaving
	}

	h, total, err := algo(reader, *k)
	if err != nil {
		log.Fatal(err)
	}

	if len(*h) == 0 {
		return
	}

	maxLen := 0
	maxCount := 0 // technically just (*h)[0].Count
	topKTotal := 0
	for _, pair := range *h {
		maxLen = max(maxLen, len(pair.Item))
		maxCount = max(maxCount, pair.Count)
		topKTotal += pair.Count
	}

	if *enableOther {
		pair := KVPair{"OTHER", 0, total - topKTotal, 0, -1}
		*h = append(*h, &pair)
		maxLen = max(maxLen, len(pair.Item))
		maxCount = max(maxCount, pair.Count)
	}

	for _, pair := range *h {
		barWidth := pair.Count * 50 / maxCount
		bar := strings.Repeat("∎", barWidth)
		fmt.Printf("%-*s  %6d  %s\n", maxLen, pair.Item, pair.Count, bar)
	}
}
