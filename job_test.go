package cron

import (
	"fmt"
	"testing"
	"time"
)

func TestJob(t *testing.T) {
	var (
		j   *Job
		cnt int
		err error
	)

	q := `@ 3s`
	ch := make(chan bool, 1)

	if j, err = newJob(q, func() {
		cnt++
	}); err != nil {
		t.Fatal(err)
	}

	j.Run()
	j.Close()

	if cnt != 1 {
		t.Fatalf("invalid count, expected %d and received %d", 1, cnt)
	}

	q = `every 1s`
	if j, err = newJob(q, func() {
		if cnt++; cnt < 4 {
			return
		}

		ch <- true
	}); err != nil {
		t.Fatal(err)
	}
	defer j.Close()
	go j.Run()

	select {
	case <-ch:
	}

	if cnt != 4 {
		t.Fatalf("invalid count, expected %d and received %d", 4, cnt)
	}

}

func TestTimeJob(t *testing.T) {
	var (
		j   *Job
		cnt int
		err error
	)

	// Get the current time
	now := time.Now()
	// Increment that time by one minute
	now = now.Add(time.Minute)
	// Set query dynamically for the current time
	q := fmt.Sprintf(`@ %s`, now.Format(timeFmt))

	// Job should run a minute from now
	if j, err = newJob(q, func() {
		cnt++
	}); err != nil {
		t.Fatal(err)
	}

	j.Run()

	if cnt != 1 {
		t.Fatalf("invalid count, expected %d and received %d", 1, cnt)
	}
}
