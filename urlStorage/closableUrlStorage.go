package urlStorage

import "io"

type ClosableUrlStorage interface {
	UrlStorage
	io.Closer
}
