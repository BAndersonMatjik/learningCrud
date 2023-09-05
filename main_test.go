package learningcrud

import (
	"fmt"
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
