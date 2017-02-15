package timertask

const ( //回调模式
	M_SYNC  = iota //同步模式
	M_ASYNC        //异步模式，会新建协程处理任务
)

const ( //定时命令
	F_SET_NEW_INTERVAL = iota
	F_DELETE_TIMER
	F_CONTINUE
)

type TimerTask interface {
	GetAlarmtime() (alarmtime int64)
	GetMode() (handleMode int)
	HandleWork(time int64) (flag int)
}

type TimerHub interface {
	//GetConfig() (alarmtime *time.Time, handleMode int)
	//HandleWork(time *time.Time) (newtime *time.Time, flag int)
	AddTask(task TimerTask) error
	DelTask(task TimerTask) error
	Running() error
}

//Ignition()(alarmtime *time.Time)
