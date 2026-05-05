package main

type keyLogger struct {
  C <-chan byte
}

func newKeyLogger() *keyLogger {
  c := make(chan byte)
  klgr := &keyLogger{
    C: c,
  }
  go func() {
    for {
      key := readKey()
      c <- key
    }
  }()
  return klgr
}

