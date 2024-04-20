//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	var mu sync.Mutex

	processIsGoing := make(chan bool)
	timeLeft := 10 - u.TimeUsed

	go func() {
		start := time.Now()
		process()
		duration := time.Since(start)

		mu.Lock()
		u.TimeUsed += int64(duration.Seconds())
		mu.Unlock()

		processIsGoing <- u.TimeUsed < 10
	}()
		
	select {
	case processElapsed, ok := <- processIsGoing:
		if !ok || !processElapsed {
			return u.IsPremium
		}

		return true 
	case <-time.After(time.Duration(timeLeft) * time.Second):
		return u.IsPremium
	}
}

func main() {
	RunMockServer()
}
