package ravenTree

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type ErrCollections struct {
	mutex sync.Mutex
	errs  []error
}

func (ec *ErrCollections) Add(errString string) {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()
	ec.errs = append(ec.errs, errors.New(errString))
}

func (ec *ErrCollections) Error() string {
	errsMessages := make([]string, len(ec.errs))
	ec.mutex.Lock()
	defer ec.mutex.Unlock()
	for i, e := range ec.errs {
		errsMessages[i] = fmt.Sprintf("%s\n", e.Error())
	}

	return strings.Join(errsMessages, "")
}

func (ec *ErrCollections) HasError() error {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()
	if len(ec.errs) == 0 {
		return nil
	}

	return ec
}

func (ec *ErrCollections) CleanCollection() {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()
	ec.errs = []error{}
}
