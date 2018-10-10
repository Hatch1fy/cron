package cron

import (
	"testing"
	"time"
)

func TestJobAfter(t *testing.T) {
	var count int
	New(func() {
		count++
	}).After(time.Millisecond * 500)

	time.Sleep(time.Millisecond * 100)
	if count != 0 {
		t.Fatalf("invalid count, expected %d and received %d", 0, count)
	}

	time.Sleep(time.Millisecond * 600)
	if count != 1 {
		t.Fatalf("invalid count, expected %d and received %d", 1, count)
	}

	time.Sleep(time.Millisecond * 600)
	if count != 1 {
		t.Fatalf("invalid count, expected %d and received %d", 1, count)
	}
}

func TestJobAt(t *testing.T) {
	var count int
	now := time.Now()
	future := now.Add(time.Second)

	New(func() {
		count++
	}).At(future)

	time.Sleep(time.Millisecond * 100)
	if count != 0 {
		t.Fatalf("invalid count, expected %d and received %d", 0, count)
	}

	time.Sleep(time.Millisecond * 1200)
	if count != 1 {
		t.Fatalf("invalid count, expected %d and received %d", 1, count)
	}

	time.Sleep(time.Millisecond * 600)
	if count != 1 {
		t.Fatalf("invalid count, expected %d and received %d", 1, count)
	}
}
