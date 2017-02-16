package timertask

import (
	"fmt"
	"sync"
	"time"

	"github.com/MoreZero/timertask/heap"
)

////taskbase//////////////////////////////////////////////////////
type Position interface {
	GetHeapPosition() (nodePosition, slicePosition int)
	SetHeapPosition(nodePosition, slicePosition int)
}
type HeaptaskBase struct {
	nodePosition  int
	slicePosition int
}

func (this *HeaptaskBase) GetHeapPosition() (nodePosition, slicePosition int) {
	return this.nodePosition, this.slicePosition
}
func (this *HeaptaskBase) SetHeapPosition(nodePosition, slicePosition int) {
	this.nodePosition, this.slicePosition = nodePosition, slicePosition
}

///堆节点切片////////////////////////////////////////////////
type slicenode struct {
	slice    []TimerTask
	stamp    int64
	position int
}

///////////////////////////////////////////////////////////////////////
type timerHeap []*slicenode

func (this timerHeap) Len() int { return len(this) }

func (this timerHeap) Less(i, j int) bool {
	return this[i].stamp < this[j].stamp
}

func (this timerHeap) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
	this[i].position, this[j].position = this[j].position, this[i].position
}

func (this *timerHeap) Push(x interface{}) {
	*this = append(*this, x.(*slicenode))
}

func (this *timerHeap) Pop() interface{} {
	old := *this
	n := len(old)
	x := old[n-1]
	*this = old[0 : n-1]
	return x
}

/////////////////////////////////////

const (
	S_STOP = iota
	S_FLUSH
)

type HeapTimer struct {
	heap       timerHeap
	nodemap    map[int64]*slicenode
	heapmax    int
	signal     chan int8
	signalflag int8
	lock       sync.Mutex
}

//budget 预堆大小, 0则默认为256
//max	堆最大大小 0则不限制
func NewHeapTimer(budget, max int) TimerHub {
	heapTimer := &HeapTimer{
		heapmax: max,
		signal:  make(chan int8, 1),
	}
	if budget == 0 {
		budget = 256
	}
	heapTimer.heap = make([]*slicenode, 0, budget)
	heapTimer.nodemap = make(map[int64]*slicenode)

	return heapTimer
}

func (this *HeapTimer) AddTask(task TimerTask) error {
	this.lock.Lock()
	alarmtime := task.GetAlarmtime()
	node, ok := this.nodemap[alarmtime]
	if ok {
		node.slice = append(node.slice, task)
		task.(Position).SetHeapPosition(node.position, len(node.slice)-1)
		this.lock.Unlock()
		return nil
	}
	node = &slicenode{
		slice: []TimerTask{task},
		stamp: alarmtime,
	}
	position := heap.Push(&this.heap, node)
	this.nodemap[alarmtime] = node
	node.position = position
	task.(Position).SetHeapPosition(position, 0)
	if this.signalflag == 0 {
		this.signalflag = 1
		this.signal <- S_FLUSH
	}
	this.lock.Unlock()
	return nil
}

func (this *HeapTimer) DelTask(task TimerTask) error {
	return nil
}

func (this *HeapTimer) Stop() {
	this.signal <- S_STOP
}

func (this *HeapTimer) Running() error {
	var timeout int64
	var sig int8
	var now int64
	var alarmnow time.Time
	for {
		fmt.Println("timeout count:", countcode, "node:", this.heap.Len())
		now = time.Now().Unix()
		this.lock.Lock()
		this.signalflag = 0
		if this.heap.Len() == 0 {
			timeout = 99999999999
		} else {
			timeout = this.heap[0].stamp - now
		}
		this.lock.Unlock()
		select {
		case sig = <-this.signal:
			if sig == S_FLUSH {
				break
			}
			if sig == S_STOP {
				return nil
			}
		case alarmnow = <-time.After(time.Duration(timeout) * time.Second):
			this.HandleTimeout(alarmnow.Unix())
		}
	}
}

/*
const ( //定时命令
	F_SET_NEW_INTERVAL = iota
	F_DELETE_TIMER
	F_CONTINUE
)
*/
var countcode int

func (this *HeapTimer) HandleTimeout(now int64) {

	for {
		this.lock.Lock()
		if this.heap.Len() == 0 {
			this.lock.Unlock()
			return
		}
		node := this.heap[0]
		if node == nil || node.stamp > now {
			this.lock.Unlock()
			return
		}
		heap.Pop(&this.heap)
		delete(this.nodemap, node.stamp)
		this.lock.Unlock()
		length := len(node.slice)
		for i := 0; i < length; i++ {
			countcode++
			task := node.slice[i]
			call := func() {
				flag := task.HandleTimeout(now)
				switch flag {
				case F_CONTINUE, F_SET_NEW_INTERVAL:
					this.AddTask(task)
				case F_DELETE_TIMER:
				}
			}
			if task.GetMode() == M_ASYNC {
				go call()
			} else {
				call()
			}

		}

	}

}
