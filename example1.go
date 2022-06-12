
/*
	This code demonstrates how to terminate an infinite loop that reads messages from
producer go routines. Logic used makes sure all the channels are closed before terminating
the loop.
	Concepts used:
	- Channels
	- Go routines
	- nil channel
	- Use fullness of closing channel
	- Getting to know when channel is closed
*/

package main
import (
	"fmt"
	"time"
)

func printMsgInMilliSec(comChan chan string)  {
	var msg string
	for i := 0; i < 5; i++ {
		msg = fmt.Sprintf("I am from printMsgInMilliSec %d", i)
		comChan <- msg
		time.Sleep(time.Millisecond * 500)
	}
	fmt.Println("printMsgInMilliSec: Closing channel")
	close(comChan)
}

func printMsgInSec(comChan chan string)  {
	var msg string
	for i := 0; i < 5; i++ {
		msg = fmt.Sprintf("I am from printMsgInSec %d", i)
		comChan <- msg
		time.Sleep(time.Second)
	}
	fmt.Println("printMsgInSec: Closing channel")
	close(comChan)
}

func main() {
	fmt.Println("main: Starts..")
	/*
		Let's create array of channels.
		Each one of those will be used by go routines to send string message
	*/
	var arrayOfChans [2] chan string
	arrayOfChans[0] = make(chan string)
	arrayOfChans[1] = make(chan string)

	/*
		Go routines those will send message over the channel
	*/
	go printMsgInMilliSec(arrayOfChans[0])
	go printMsgInSec(arrayOfChans[1])

	/*
		Below is the logic that terminates the loop once all the channels are closed
	by respective go routines.
		Observation:
			- After channels are closed if we do not assign "nil" to closed channel,
			then that closed channel error out in "case" statement. If you want to
			experience the issue then simply comment out below code -
				- arrayOfChans[0] = nil and
				- arrayOfChans[1] = nil
	*/
	for i := len(arrayOfChans); i > 0; {
		select {
			case milliMsg, err1 := <- arrayOfChans[0]:
				if err1 {
					fmt.Println(milliMsg)
				} else {
					/* First channel is closed. So we need to assign nil to closed channel */
					i--
					arrayOfChans[0] = nil
				}
			case secMsg, err2 := <- arrayOfChans[1]:
				if err2 {
					fmt.Println(secMsg)
				} else {
					/* Second channel is closed. So we need to assign nil to closed channel */
					i--
					arrayOfChans[1] = nil
				}
		}
	}

	fmt.Println("main: Done!!")
}
