package render

import (
	"container/heap"
	"golang.org/x/exp/shiny/screen"
)

var (
	rh *RenderableHeap
)

type RenderableHeap []Renderable

func (h RenderableHeap) Len() int           { return len(h) }
func (h RenderableHeap) Less(i, j int) bool { return h[i].GetLayer() < h[j].GetLayer() }
func (h RenderableHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *RenderableHeap) Push(x interface{}) {
	*h = append(*h, x.(Renderable))
}

func (h_p *RenderableHeap) Pop() interface{} {
	h := *h_p
	n := len(h)
	x := h[n-1]
	*h_p = h[0 : n-1]
	return x
}

func InitDrawHeap() {
	rh = &RenderableHeap{}
	heap.Init(rh)

}

func Draw(r Renderable, l int) Renderable {
	// Bind to PostDraw if this causes synchronization issues with DrawHeap
	r.SetLayer(l)
	heap.Push(rh, r)
	return r
}

func LoadSpriteAndDraw(filename string, l int) *Sprite {
	s := LoadSprite(filename)
	return Draw(s, l).(*Sprite)
}

func DrawHeap(b screen.Buffer) {
	newRh := &RenderableHeap{}
	for rh.Len() > 0 {
		r := heap.Pop(rh).(Renderable)
		r.Draw(b)
		heap.Push(newRh, r)
	}
	rh = newRh
}