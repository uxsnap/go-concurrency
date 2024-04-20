//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"time"
)

func getTweet(stream Stream, c chan *Tweet) {
	tweet, err := stream.Next()

	if err == ErrEOF {
		close(c)
		return
	}

	c <- tweet
	go getTweet(stream, c)
}

func producer(stream Stream, c chan *Tweet) {
	go getTweet(stream, c)
}

func consumer(c chan *Tweet) {
	for {
		if c == nil {
			break
		}
		
		t, ok := <-c;

		if !ok {
			break
		}

		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	ch := make(chan *Tweet, 10)

	// Producer
	producer(stream, ch)

	// Consumer
	consumer(ch)

	fmt.Printf("Process took %s\n", time.Since(start))
}
