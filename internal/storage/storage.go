package storage

import (
  "err"
)

var (
  ErrURLNotFound = err.New("url not found")
  ErrURLExists = err.New("url exists")
)
