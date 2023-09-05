package learningcrud

import (
	"fmt"
	"strconv"
	"sync"
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
