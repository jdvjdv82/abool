package abool

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/segmentio/encoding/json"
)

func TestDefaultValue(t *testing.T) {
	t.Parallel()
	v := New()
	if v.IsSet() {
		t.Fatal("Empty value of AtomicBool should be false")
	}

	v = NewBool(true)
	if !v.IsSet() {
		t.Fatal("NewValue(true) should be true")
	}

	v = NewBool(false)
	if v.IsSet() {
		t.Fatal("NewValue(false) should be false")
	}
}

func TestIsNotSet(t *testing.T) {
	t.Parallel()
	v := New()

	if v.IsSet() == v.IsNotSet() {
		t.Fatal("AtomicBool.IsNotSet() should be the opposite of IsSet()")
	}
}

func TestSetUnSet(t *testing.T) {
	t.Parallel()
	v := New()

	v.Set()
	if !v.IsSet() {
		t.Fatal("AtomicBool.Set() failed")
	}

	v.UnSet()
	if v.IsSet() {
		t.Fatal("AtomicBool.UnSet() failed")
	}
}

func TestSetTo(t *testing.T) {
	t.Parallel()
	v := New()

	v.SetTo(true)
	if !v.IsSet() {
		t.Fatal("AtomicBool.SetTo(true) failed")
	}

	v.SetTo(false)
	if v.IsSet() {
		t.Fatal("AtomicBool.SetTo(false) failed")
	}

	if set := v.SetToIf(true, false); set || v.IsSet() {
		t.Fatal("AtomicBool.SetTo(true, false) failed")
	}

	if set := v.SetToIf(false, true); !set || !v.IsSet() {
		t.Fatal("AtomicBool.SetTo(false, true) failed")
	}
}

func TestRace(t *testing.T) {
	repeat := 10000
	var wg sync.WaitGroup
	wg.Add(repeat * 3)
	v := New()

	// Writer
	go func() {
		for i := 0; i < repeat; i++ {
			v.Set()
			wg.Done()
		}
	}()

	// Reader
	go func() {
		for i := 0; i < repeat; i++ {
			v.IsSet()
			wg.Done()
		}
	}()

	// Writer
	go func() {
		for i := 0; i < repeat; i++ {
			v.UnSet()
			wg.Done()
		}
	}()
	wg.Wait()
}

func TestJSONCompatibleWithBuiltinBool(t *testing.T) {
	for _, value := range []bool{true, false} {
		// Test bool -> bytes -> AtomicBool

		// act 1. bool -> bytes
		buf, err := json.Marshal(value)
		if err != nil {
			t.Fatalf("json.Marshal(%t) failed: %s", value, err)
		}

		// act 2. bytes -> AtomicBool
		//
		// Try to unmarshall the JSON byte slice
		// of a normal boolean into an AtomicBool
		//
		// Create an AtomicBool with the oppsite default to ensure the unmarshal process did the work
		ab := NewBool(!value)
		err = json.Unmarshal(buf, ab)
		if err != nil {
			t.Fatalf(`json.Unmarshal("%s", %T) failed: %s`, buf, ab, err)
		}
		// assert
		if ab.IsSet() != value {
			t.Fatalf("Expected AtomicBool to represent %t but actual value was %t", value, ab.IsSet())
		}

		// Test AtomicBool -> bytes -> bool

		// act 3. AtomicBool -> bytes
		buf, err = json.Marshal(ab)
		if err != nil {
			t.Fatalf("json.Marshal(%T) failed: %s", ab, err)
		}

		// using the opposite value for the same reason as the former case
		b := ab.IsNotSet()
		// act 4. bytes -> bool
		err = json.Unmarshal(buf, &b)
		if err != nil {
			t.Fatalf(`json.Unmarshal("%s", %T) failed: %s`, buf, &b, err)
		}
		// assert
		if b != ab.IsSet() {
			t.Fatalf(`json.Unmarshal("%s", %T) didn't work, expected %t, got %t`, buf, ab, ab.IsSet(), b)
		}
	}
}

func TestUnmarshalJSONErrorNoWrite(t *testing.T) {
	for _, val := range []bool{true, false} {
		ab := NewBool(val)
		oldVal := ab.IsSet()
		buf := []byte("invalid-json")
		err := json.Unmarshal(buf, ab)
		if err == nil {
			t.Fatalf(`Error expected from json.Unmarshal("%s", %T)`, buf, ab)
		}
		if oldVal != ab.IsSet() {
			t.Fatal("Failed json.Unmarshal modified the value of AtomicBool which is not expected")
		}
	}
}

// Benchmark Read

func BenchmarkMutexRead(b *testing.B) {
	var m sync.RWMutex
	var v bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.RLock()
		_ = v
		m.RUnlock()
	}
}

func BenchmarkAtomicValueRead(b *testing.B) {
	var v atomic.Value
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Load() != nil
	}
}

func BenchmarkAtomicBoolRead(b *testing.B) {
	v := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.IsSet()
	}
}

// Benchmark Write

func BenchmarkMutexWrite(b *testing.B) {
	var m sync.RWMutex
	var v bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.RLock()
		v = true
		m.RUnlock()
	}
	b.StopTimer()
	_ = v
}

func BenchmarkAtomicValueWrite(b *testing.B) {
	var v atomic.Value
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Store(true)
	}
}

func BenchmarkAtomicBoolWrite(b *testing.B) {
	v := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Set()
	}
}

// Benchmark CAS

func BenchmarkMutexCAS(b *testing.B) {
	var m sync.RWMutex
	var v bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Lock()
		if !v {
			v = true
		}
		m.Unlock()
	}
}

func BenchmarkAtomicBoolCAS(b *testing.B) {
	v := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.SetToIf(false, true)
	}
}
