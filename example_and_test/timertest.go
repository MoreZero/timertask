package main

import (
	"fmt"
	"time"

	"github.com/MoreZero/timertask"
	"github.com/neverlee/glog"
)

/*
type TestTimer struct {
	I int
}

func (this *TestTimer) TimerFunc(now *time.Time) (newtimer time.Duration, flag int) {
	this.I++
	fmt.Println("called TestTimer::TimerFunc")
	if this.I == 5 {
		return 1 * time.Second, timertask.SET_NEW_INTERVAL
	} else if this.I == 10 {
		return 0, timertask.DELETE_TIMER
	}
	return 0, timertask.CONTINUE
}

func main() {
	timertask.StartSingleTimer(&TestTimer{}, 2*time.Second)
}
*/

type TestTask struct {
	timertask.HeaptaskBase
	alarmtime int64
	id        int
	count     int
}

func (this *TestTask) GetAlarmtime() (alarmtime int64) {
	this.alarmtime += 5
	return this.alarmtime
}

func (this *TestTask) GetMode() (handleMode int) {
	return timertask.M_ASYNC
}

func (this *TestTask) HandleTimeout(time1 int64) (flag int) {
	//fmt.Println("handleWork", time1, "now:", time.Now().Unix(), "id", this.id)
	glog.Extraln("self:", this.alarmtime, "handleWork", time1, "now:", time.Now().Unix(), "id", this.id)
	this.count += 1
	if this.count == 5 {
		return timertask.F_DELETE_TIMER
	}
	return timertask.F_CONTINUE
}

func addtask(i int, hub timertask.TimerHub) {
	task := &TestTask{
		alarmtime: time.Now().Unix() + 10,
		id:        i,
	}
	hub.AddTask(task)
}

func hubtestadd(hub timertask.TimerHub) {
	/*
		for i := 0; i < 200; i++ {
			go addtask(i, hub)
		}
		time.Sleep(1 * time.Second)
		for i := 200; i < 400; i++ {
			go addtask(i, hub)
		}*/
	for {
		for i := 0; i < 100000; i++ {
			go addtask(i, hub)
		}
		time.Sleep(30 * time.Second)
	}
}

func main() {
	var hub timertask.TimerHub
	hub = timertask.NewHeapTimer(100, 100)
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("start hub add test")
		go hubtestadd(hub)
	}()
	hub.Running()
}
