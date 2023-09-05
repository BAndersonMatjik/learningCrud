package learningcrud

import (
	"fmt"
	"strconv"
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
		close(channel)
	}()

	for data := range channel {
		fmt.Println(data)
	}
}
