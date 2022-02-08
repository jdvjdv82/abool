# ABool :bulb:

Atomic Boolean package for Go, optimized for performance yet simple to use.

Designed for cleaner code.

## Usage

```go
package abool


cond := abool.New()     // default to false

cond.Set()              // Sets to true
cond.IsSet()            // Returns true
cond.UnSet()            // Sets to false
cond.IsNotSet()         // Returns true
cond.SetTo(any)         // Sets to whatever you want
cond.SetToIf(old, new)  // Sets to `new` only if the Boolean matches the `old`, returns whether succeeded


// embedding
type Foo struct {
    cond *abool.AtomicBool  // always use pointer to avoid copy
}
```

## Benchmark:

- Go 1.6.2
- OS X 10.11.4
- Intel CPU (to be specified)


```
# Read
BenchmarkMutexRead-48                   656360493                8.691 ns/op           0 B/op          0 allocs/op
BenchmarkAtomicValueRead-48             1000000000               0.2631 ns/op          0 B/op          0 allocs/op
BenchmarkAtomicBoolRead-48              1000000000               0.2515 ns/op          0 B/op          0 allocs/op

# Write
BenchmarkMutexWrite-48                  693644023                8.755 ns/op           0 B/op          0 allocs/op
BenchmarkAtomicValueWrite-48            980383506                5.572 ns/op           0 B/op          0 allocs/op
BenchmarkAtomicBoolWrite-48             1000000000               4.468 ns/op           0 B/op          0 allocs/op

# CAS
BenchmarkMutexCAS-48                    316539304               18.47 ns/op            0 B/op          0 allocs/op
BenchmarkAtomicBoolCAS-48               1000000000               4.179 ns/op           0 B/op          0 allocs/op




```

## Special thanks to contributors

- [barryz](https://github.com/barryz)
  - Added the `Toggle` method
- [Lucas Rouckhout](https://github.com/LucasRouckhout)
  - Implemented JSON Unmarshal and Marshal interface
- [Sebastian Schicho](https://github.com/schicho)
  - Reported a regression with test case

