package coze

type Pagination[T any] interface {
	Err() error
	Items() []*T
	Current() *T
	Total() int
	HasMore() bool
}

type PageRequest struct {
	PageToken string `json:"page_token,omitempty"`
	PageNum   int    `json:"page_num,omitempty"`
	PageSize  int    `json:"page_size,omitempty"`
}

type PageResponse[T any] struct {
	HasMore bool   `json:"has_more"`
	Total   int    `json:"total"`
	Data    []*T   `json:"data"`
	LastID  string `json:"last_id,omitempty"`
	NextID  string `json:"next_id,omitempty"`
	LogID   string `json:"log_id,omitempty"`
}

type basePager[T any] struct {
	pageFetcher    PageFetcher[T]
	pageSize       int
	currentPage    *PageResponse[T]
	currentIndex   int
	currentPageNum int
	cur            *T
	err            error
}

func (p *basePager[T]) Err() error {
	return p.err
}

func (p *basePager[T]) Items() []*T {
	return ptrValue(p.currentPage).Data
}

func (p *basePager[T]) Current() *T {
	return p.cur
}

func (p *basePager[T]) Total() int {
	return ptrValue(p.currentPage).Total
}

func (p *basePager[T]) HasMore() bool {
	return ptrValue(p.currentPage).HasMore
}

// PageFetcher interface
type PageFetcher[T any] func(request *PageRequest) (*PageResponse[T], error)

// NumberPaged implementation
type NumberPaged[T any] struct {
	basePager[T]
}

func NewNumberPaged[T any](fetcher PageFetcher[T], pageSize, pageNum int) (*NumberPaged[T], error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	paginator := &NumberPaged[T]{basePager: basePager[T]{pageFetcher: fetcher, pageSize: pageSize, currentPageNum: pageNum}}
	err := paginator.fetchNextPage()
	if err != nil {
		return nil, err
	}
	return paginator, nil
}

func (p *NumberPaged[T]) fetchNextPage() error {
	request := &PageRequest{PageNum: p.currentPageNum, PageSize: p.pageSize}
	var err error
	p.currentPage, err = p.pageFetcher(request)
	if err != nil {
		return err
	}
	p.currentIndex = 0
	p.currentPageNum++
	return nil
}

func (p *NumberPaged[T]) Next() bool {
	if p.currentIndex < len(ptrValue(p.currentPage).Data) {
		p.cur = p.currentPage.Data[p.currentIndex]
		p.currentIndex++
		return true
	}
	if p.currentPage.HasMore {
		err := p.fetchNextPage()
		if err != nil {
			p.err = err
			return false
		}
		if len(p.currentPage.Data) == 0 {
			return false
		}
		p.cur = p.currentPage.Data[p.currentIndex]
		p.currentIndex++
		return true
	}
	return false
}

// TokenPaged implementation
type TokenPaged[T any] struct {
	basePager[T]
	pageToken *string
}

func NewTokenPaged[T any](fetcher PageFetcher[T], pageSize int, nextID *string) (*TokenPaged[T], error) {
	paginator := &TokenPaged[T]{basePager: basePager[T]{pageFetcher: fetcher, pageSize: pageSize}, pageToken: nextID}
	err := paginator.fetchNextPage()
	if err != nil {
		return nil, err
	}
	return paginator, nil
}

func (p *TokenPaged[T]) fetchNextPage() error {
	request := &PageRequest{PageToken: ptrValue(p.pageToken), PageSize: p.pageSize}
	var err error
	p.currentPage, err = p.pageFetcher(request)
	if err != nil {
		return err
	}
	p.currentIndex = 0
	p.pageToken = &p.currentPage.NextID
	return nil
}

func (p *TokenPaged[T]) Next() bool {
	if p.currentIndex < len(ptrValue(p.currentPage).Data) {
		p.cur = p.currentPage.Data[p.currentIndex]
		p.currentIndex++
		return true
	}
	if p.currentPage.HasMore {
		err := p.fetchNextPage()
		if err != nil {
			p.err = err
			return false
		}
		if len(p.currentPage.Data) == 0 {
			return false
		}
		p.cur = p.currentPage.Data[p.currentIndex]
		p.currentIndex++
		return true
	}
	return false
}
