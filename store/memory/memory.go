package memory

import "sync"

func init() {

}

var addressMemory = newAddressCacher()

type addressCacher struct {
	m sync.Map
}

func newAddressCacher() *addressCacher {
	return &addressCacher{
		m: sync.Map{},
	}
}
