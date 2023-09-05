package learningcrud

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func RunHelloWorld() {
	println("Hello World")
}

func DisplayNumber(number int) {
	fmt.Println("Display ", number)
}

func TestManyGoroutine(t *testing.T) {
	for i := 0; i < 100000; i++ {
		go DisplayNumber(i)
	}
	time.Sleep(5 * time.Second)
}

func TestGoroutine(t *testing.T) {
	//declare goroutine 2Kb vs Thread 1000Kb/1Mb
	go RunHelloWorld()
	println("Ups")
	time.Sleep(1 * time.Second)
}

func TestGoroutineChannel(t *testing.T) {
	channel := make(chan string)
	defer close(channel)
	//Becareful all channel need send or receive when no one accepted result or send will deadlock
	go func() {
		time.Sleep(2 * time.Second)
		//send data
		channel <- "Billy"
		fmt.Println("Finish Send Channel")
	}()
	//get data from channel
	result := <-channel
	fmt.Println(result)
	time.Sleep(5 * time.Second)
}

// Not using Pointer
func RunningChannelAsFunction(channel chan string) {
	time.Sleep(2 * time.Second)
	//send data
	channel <- "Billy"
	fmt.Println("Finish Send Channel")

}

func TestGoroutineChannelAsFunction(t *testing.T) {
	channel := make(chan string)
	defer close(channel)
	//Becareful all channel need send or receive when no one accepted result or send will deadlock
	go RunningChannelAsFunction(channel)
	//get data from channel
	result := <-channel
	fmt.Println(result)
	time.Sleep(5 * time.Second)
}

// get in channel
func RunningChannelJustIn(channel chan<- string) {
	time.Sleep(2 * time.Second)
	//send data
	channel <- "Billy"
	fmt.Println("Finish Send Channel")

}

// get out channel
func RunningChannelJustOut(channel <-chan string) {
	//get data
	result := <-channel
	fmt.Println(result)

}

// can tell channel just for in data
func TestGoroutineChannelJustInAndOut(t *testing.T) {
	channel := make(chan string)
	defer close(channel)
	//Becareful all channel need send or receive when no one accepted result or send will deadlock
	go RunningChannelJustIn(channel)
	//get data from channel
	go RunningChannelJustOut(channel)

	time.Sleep(5 * time.Second)
}

func RunningMany() {

}

// channel just will get 1 data if wanna get 2 or more using buffered channel
// but buffered channel must have capacity no miminum and max
// if using buffered not will deadlock if data not get used
func TestGoroutineBufferedChannel(t *testing.T) {
	channel := make(chan string, 3) //make channel buffer
	defer close(channel)

	go func() {
		channel <- "Billy"
		channel <- "iqbal"
		fmt.Println("buffer data size", len(channel)) //show data size
		channel <- "ridhwan"
	}()
	go func() {
		fmt.Println(<-channel)
		fmt.Println(<-channel)
		fmt.Println(<-channel)
	}()
	time.Sleep(2 * time.Second)
	fmt.Println("Done")
	fmt.Println("buffer size", cap(channel)) //show buffer size
}

// Range will do iteration data in to channel when channel is close iteration will stop
func TestGoroutineRangeChannel(t *testing.T) {
	channel := make(chan string) //make channel buffer
	//more good best practice close when done
	go func() {
		for i := 0; i < 10; i++ {
			channel <- "Perulangan ke " + strconv.Itoa(i)
		}
		// always close when is done do process if not will deadlock
		close(channel)
	}()

	for data := range channel {
		fmt.Println(data)
	}
	fmt.Println("end")
}

func GiveMeResponse(channel chan<- string) {
	time.Sleep(2 * time.Second)
	channel <- "data kuy"
}

// if we need more than one channel because on rangechannel is not posible using multiple channel
// this will get faster come data will be serve result will be random because listen to 2 channel
func TestGoroutineSelectChannel(t *testing.T) {
	channel1 := make(chan string)
	channel2 := make(chan string)

	defer close(channel1)
	defer close(channel2)
	//more good best practice close when done
	counter := 0
	go GiveMeResponse(channel1)
	go GiveMeResponse(channel2)
	for {
		select {
		case data := <-channel1:
			fmt.Println("data channel 1 ", data)
			counter++
		case data := <-channel2:
			fmt.Println("data channel 2 ", data)
			counter++
		default:
			fmt.Println("Waiting Data")
		}

		if counter == 2 {
			break
		}
	}
}

// add default but we select channel empty
func TestGoroutineDefaultSelectChannel(t *testing.T) {
	channel1 := make(chan string)
	channel2 := make(chan string)
	//more good best practice close when done
	counter := 0
	go GiveMeResponse(channel1)
	go GiveMeResponse(channel2)

	for {
		select {
		case data := <-channel1:
			fmt.Println("data channel 1 ", data)
			counter++
		case data := <-channel2:
			fmt.Println("data channel 2 ", data)
			counter++
		default:
			fmt.Println("Waiting Data")
		}
		if counter == 2 {
			break
		}
	}
}

// sharing variable using by multiple goroutine
// WILL BE RACE CONDITION
func TestGoroutineIssueRaceCondition(t *testing.T) {

	//example the issue x will be use multiple goroutine
	x := 0
	for i := 1; i <= 1000; i++ {
		go func() {
			for j := 1; j <= 100; j++ {
				x = x + 1
			}
		}()
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Counter : ", x)
}

// Mutual Exclusion (Solution Race Condition)
// little bit performance impact
// to handling issue race condition will be lock the variable
// just one goroutine will be give permission to lock and unlock if goroutine unlock the mutex next goroutine can access mutex
func TestGoroutineMutex(t *testing.T) {
	//example the issue x will be use multiple goroutine
	x := 0
	var mutex sync.Mutex
	for i := 1; i <= 1000; i++ {
		go func() {
			for j := 1; j <= 100; j++ {
				mutex.Lock()
				x = x + 1
				mutex.Unlock()
			}
		}()
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Counter : ", x)
}

type BankAccount struct {
	RWMutex sync.RWMutex
	Balance int
}

func (account *BankAccount) AddBalance(amount int) {
	account.RWMutex.Lock()
	account.Balance = account.Balance + amount
	account.RWMutex.Unlock()
}

func (account *BankAccount) GetBalance() int {
	account.RWMutex.RLock()
	balance := account.Balance
	account.RWMutex.RUnlock()
	return balance
}

// Mutual Exclusion Read write
// handling read and write data
func TestGoroutineMutexReadAndWrite(t *testing.T) {
	//example the issue x will be use multiple goroutine
	account := BankAccount{}

	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				account.AddBalance(1)
				fmt.Println(account.GetBalance())
			}
		}()
	}
	time.Sleep(5 * time.Second)
}

type UserBalance struct {
	sync.Mutex
	Name    string
	Balance int
}

func (user *UserBalance) Lock() {
	user.Mutex.Lock()
}

func (user *UserBalance) Unlock() {
	user.Mutex.Unlock()
}

func (user *UserBalance) ChangeAmount(amount int) {
	user.Balance = user.Balance + amount
}

func Transfer(user1 *UserBalance, user2 *UserBalance, amount int) {
	user1.Lock()
	fmt.Println("Lock " + user1.Name)
	user1.ChangeAmount(-amount)

	time.Sleep(2 * time.Second)

	user2.Lock()
	fmt.Println("Lock " + user2.Name)
	user2.ChangeAmount(amount)

	time.Sleep(2 * time.Second)

	user1.Unlock()
	user2.Unlock()

}

// Becarefull DeadLock using goroutine
// Because waiting locking mutex
func TestGoroutineIssueDeadlock(t *testing.T) {
	//example the issue x will be use multiple goroutine

	user1 := UserBalance{
		Name:    "billy",
		Balance: 1000000,
	}
	user2 := UserBalance{
		Name:    "ridhwan",
		Balance: 1000000,
	}

	//user 1 get lock
	go Transfer(&user1, &user2, 100000)
	//user 2 get lock
	go Transfer(&user2, &user1, 200000)
	//will be lock

	time.Sleep(10 * time.Second)

	fmt.Println("User 1 balance ", user1.Name, "----", user1.Balance)
	fmt.Println("User 2 balance ", user2.Name, "----", user2.Balance)
}

func RunAsynchronous(group *sync.WaitGroup) {
	defer group.Done()

	group.Add(1)
	fmt.Println("Hello")
	time.Sleep(1 * time.Second)
}

// wait group
// wait until process done
// Add(int) add goroutine and is Done() when goroutine done
// Wait() will waiting until added coroutine size add before
// BECAREFULL RUNASYNCHRONOUS NOT WAITGROUP NOT DONE WILL BE DEADLOCK
func TestGoroutineWaitGroup(t *testing.T) {
	group := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		go RunAsynchronous(group)
	}
	group.Wait()
	fmt.Println("Complete")
}

var counter = 0

func OnlyOnce() {
	counter++
}

// function will be once executing
// when many goroutine access the first execution and second execution will be ignore
func TestGoroutineOnce(t *testing.T) {
	var once sync.Once
	var group sync.WaitGroup

	for i := 0; i < 100; i++ {
		go func() {
			group.Add(1)
			once.Do(OnlyOnce)
			group.Done()
		}()
	}
	group.Wait()
	fmt.Println(counter)

}

// POOL DESIGN PATTERN
// USING FOR CONCURRENCY
// Pool for save data and we can get data when we wanna use
// POOL save problem race condition
func TestGoroutinePool(t *testing.T) {
	pool := sync.Pool{
		//make default value to change null value
		New: func() interface{} {
			return "new"
		},
	}

	pool.Put("Billy")
	pool.Put("test")
	pool.Put("helele")

	for i := 0; i < 20; i++ {
		go func() {
			//becarefull when get data will gone from pool
			data := pool.Get()
			fmt.Println(data)
			time.Sleep(1 * time.Second)
			pool.Put(data)

		}()
	}
	time.Sleep(11 * time.Second)
}

// Map for save from save condition
func TestGoroutineMap(t *testing.T) {
	data := sync.Map{}
	waitGroup := sync.WaitGroup{}

	var addToMap = func(value int) {
		defer waitGroup.Done()
		waitGroup.Add(1)
		data.Store(value, value)
	}

	for i := 0; i < 100; i++ {
		go addToMap(i)
	}
	waitGroup.Wait()

	data.Range(func(key, value interface{}) bool {
		fmt.Println(key, ":", value)
		return true
	})

}

// Lock base on condition
// when cond wall in wait will stop
// if Signal() just give one goroutine wait just one go if multiple use Broadcast
func TestGoroutineCond(t *testing.T) {
	cond := sync.NewCond(&sync.Mutex{})
	waitGroup := sync.WaitGroup{}

	var waitCondition = func(value int) {
		defer waitGroup.Done()
		cond.L.Lock()
		cond.Wait()
		fmt.Println("Done ", value)
		cond.L.Unlock()
	}
	for i := 0; i < 10; i++ {
		waitGroup.Add(1)
		go waitCondition(i)
	}

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			cond.Signal()
		}
	}()

	waitGroup.Wait()

}

// Atomic to using primitif data in concurrent
// without using mutex
// doc https//golang.org/pkg/sync/atomic/
func TestGoroutineAtomic(t *testing.T) {
	waitGroup := sync.WaitGroup{}
	var counter int64 = 0

	for i := 0; i < 100; i++ {
		waitGroup.Add(1)
		go func() {
			for j := 0; j < 100; j++ {
				atomic.AddInt64(&counter, 1)
			}
			waitGroup.Done()
		}()
	}

	waitGroup.Wait()
	fmt.Println("Counter", counter)

}

// Represent on activity if timer expired and will be send by channel
func TestGoroutineTimer(t *testing.T) {
	// timer := time.NewTimer(5 * time.Second)
	fmt.Println(time.Now())
	// time := <-timer.C
	// fmt.Println(time)

	//simplfy
	//we just need channel not need data timer
	channel := time.After(5 * time.Second)

	fmt.Println(<-channel)

	group := sync.WaitGroup{}

	group.Add(1)

	//wanna running function after timer
	time.AfterFunc(5*time.Second, func() {
		fmt.Println("Excuted")
		group.Done()
	})

	group.Wait()

}

// represent repeat, if ticker expired will send channel
// every 5 second send channel
func TestGoroutineTicker(t *testing.T) {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		time.Sleep(5 * time.Second)
		ticker.Stop()
	}()

	for tick := range ticker.C {
		fmt.Println(tick)
	}
}

func TestGoroutineSimpleTicker(t *testing.T) {
	ticker := time.Tick(1 * time.Second)

	for tick := range ticker {
		fmt.Println(tick)
	}
}

// wanna know how many thread and we can change
// default thread == cpu
func TestGOMAXPROCS(t *testing.T) {
	fmt.Println(runtime.GOMAXPROCS(-1))
	fmt.Println(runtime.NumCPU())
	fmt.Println(runtime.NumGoroutine())
}
