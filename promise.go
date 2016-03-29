package fp

import (
	"reflect"
)

// if ok is false, the promise is closed without receiving any.
func (this Promise) Recv() (a AnyVal, ok bool) {
	av, ok := this.vq.Recv()
	a = av.Interface()
	return
}

func (this Promise) IsEmpty() bool {
	return 0 == len(this.q)
}

func (this *Promise) close() {
	this.m.Lock()
	defer this.m.Unlock()
	close(this.q)
	this.closed = true
}

func (this *Promise) Then(fns ...Func1) (p *Promise) {
	return Future(func() (AnyVal, error) {
		res := <-this.q
		if this.HasError() {
			if 1 < len(fns) {
				return fns[1](this.Error())
			}
		} else {
			if 0 < len(fns) {
				return fns[0](res)
			}
		}
		return nil, nil
	})
}

func (this *Promise) Q() ChanAny {
	return this.q
}

func (this *Promise) send(a AnyVal) {
	this.q <- a
}

func (this Promise) HasError() bool {
	return nil != this.err
}

func (this Promise) Error() error {
	return this.err
}

func Fcall(v AnyVal) *Promise {
	p := makepromise()
	p.q <- v
	close(p.q)
	return p
}

func makepromise(chanlen ...int) (p *Promise) {
	p = new(Promise)
	if 0 < len(chanlen) {
		if 0 == chanlen[0] {
			panic("does not support chanlen=0. could break promise.")
		}
		p.q = make(ChanAny, chanlen[0])
	} else {
		p.q = make(ChanAny, 1)
	}
	p.vq = reflect.ValueOf(p.q)
	return
}
