package main

import (
	"fmt"
	"github.com/MoreZero/timertask"
	"time"
)

type TestTask struct {
	timertask.HeaptaskBase
	alarmtime int64
}

func (this *TestTask) GetAlarmtime() (alarmtime int64) {
	return this.alarmtime + 5
}

func (this *TestTask) GetMode() (handleMode int) {
	return timertask.M_SYNC
}

func (this *TestTask) HandleWork(time int64) (flag int) {
	fmt.Println("handleWork", time)
	return timertask.F_CONTINUE
}

func main() {
	var hub timertask.TimerHub
	task := &TestTask{
		alarmtime: time.Now().Unix() + 10,
	}
	hub = timertask.NewHeapTimer(100, 100)
	go func() {
		time.Sleep(10 * time.Second)
		hub.AddTask(task)
	}()
	hub.Running()
}
