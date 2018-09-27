package cron

import (
	"testing"
)

func TestParseQueryAt(t *testing.T) {
	q := "@ 300ms"
	desc, ref, val, loc, err := parseQuery(q)
	if err != nil {
		t.Fatal(err)
	}

	if desc != "@" {
		t.Fatalf("Invalid description, expected \"%s\" and received \"%s\"", "@", desc)
	}

	if ref != "ms" {
		t.Fatalf("Invalid reference, expected \"%s\" and received \"%s\"", "ms", ref)
	}

	if val != 300 {
		t.Fatalf("Invalid value, expected \"%d\" and received \"%d\"", 300, val)
	}

	if loc != nil {
		t.Fatalf("Invalid location, expected nil and received %v", loc)
	}
}

func TestParseQueryEvery(t *testing.T) {
	q := "every 2h"
	desc, ref, val, loc, err := parseQuery(q)
	if err != nil {
		t.Fatal(err)
	}

	if desc != "every" {
		t.Fatalf("Invalid description, expected \"%s\" and received \"%s\"", "every", desc)
	}

	if ref != "h" {
		t.Fatalf("Invalid reference, expected \"%s\" and received \"%s\"", "h", ref)
	}

	if val != 2 {
		t.Fatalf("Invalid value, expected \"%d\" and received \"%d\"", 2, val)
	}

	if loc != nil {
		t.Fatalf("Invalid location, expected nil and received %v", loc)
	}
}

func TestParseMidnight(t *testing.T) {
	q := "@ midnight -07"
	desc, ref, val, loc, err := parseQuery(q)
	if err != nil {
		t.Fatal(err)
	}

	if desc != "@" {
		t.Fatalf("Invalid description, expected \"%s\" and received \"%s\"", "@", desc)
	}

	if ref != "time" {
		t.Fatalf("Invalid reference, expected \"%s\" and received \"%s\"", "00:00", ref)
	}

	if val != 0 {
		t.Fatalf("Invalid value, expected \"%d\" and received \"%d\"", 0, val)
	}

	if loc.String() != "-07" {
		t.Fatalf("invalid loc, expected %s and received %s", "-07", loc.String())
	}
}
