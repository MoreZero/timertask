package timertask

import (
	"time"
)

const (
	SET_NEW_INTERVAL = iota
	DELETE_TIMER
	CONTINUE
)

type SingleTimer interface {
	TimerFunc(now *time.Time) (newtimer time.Duration, flag int)
}

func StartSingleTimer(st SingleTimer, interval time.Duration) {
	for {
		now := time.Now()
		newtimer, flag := st.TimerFunc(&now)
		switch flag {
		case CONTINUE:
			break
		case SET_NEW_INTERVAL:
			interval = newtimer
		case DELETE_TIMER:
			return
		}
		time.Sleep(interval)
	}
}
