package host

func NewPagerResult() *PagerResult {
	return &PagerResult{
		Data: NewHostSet(),
	}
}

type PagerResult struct {
	Data    *HostSet
	Err     error
	HasNext bool
}

type Pager interface {
	Next() *PagerResult
}
