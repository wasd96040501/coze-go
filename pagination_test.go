package coze

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestData 测试数据结构
type TestData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// mockDataSource 模拟数据源
type mockDataSource struct {
	data  []*TestData
	total int
}

// newMockDataSource 创建模拟数据源
func newMockDataSource(total int) *mockDataSource {
	data := make([]*TestData, total)
	for i := 0; i < total; i++ {
		data[i] = &TestData{
			ID:   i + 1,
			Name: fmt.Sprintf("test-%d", i+1),
		}
	}
	return &mockDataSource{
		data:  data,
		total: total,
	}
}

// getNumberPageData 获取基于页码的分页数据
func (m *mockDataSource) getNumberPageData(request *PageRequest) (*PageResponse[TestData], error) {
	pageSize := request.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	startIndex := (request.PageNum - 1) * pageSize
	if startIndex >= len(m.data) {
		return &PageResponse[TestData]{
			HasMore: false,
			Total:   len(m.data),
			Data:    []*TestData{},
		}, nil
	}

	endIndex := startIndex + pageSize
	if endIndex > len(m.data) {
		endIndex = len(m.data)
	}

	return &PageResponse[TestData]{
		HasMore: endIndex < len(m.data),
		Total:   len(m.data),
		Data:    m.data[startIndex:endIndex],
	}, nil
}

// getTokenPageData 获取基于令牌的分页数据
func (m *mockDataSource) getTokenPageData(request *PageRequest) (*PageResponse[TestData], error) {
	pageSize := request.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	startIndex := 0
	if request.PageToken != "" {
		// 简单模拟：令牌就是上一页最后一条数据的 ID
		for i, item := range m.data {
			if fmt.Sprintf("%d", item.ID) == request.PageToken {
				startIndex = i + 1
				break
			}
		}
	}

	if startIndex >= len(m.data) {
		return &PageResponse[TestData]{
			HasMore: false,
			Total:   len(m.data),
			Data:    []*TestData{},
		}, nil
	}

	endIndex := startIndex + pageSize
	if endIndex > len(m.data) {
		endIndex = len(m.data)
	}

	var nextID string
	if endIndex < len(m.data) {
		nextID = fmt.Sprintf("%d", m.data[endIndex-1].ID)
	}

	return &PageResponse[TestData]{
		HasMore: endIndex < len(m.data),
		Total:   len(m.data),
		Data:    m.data[startIndex:endIndex],
		NextID:  nextID,
	}, nil
}

func TestNumberPaged(t *testing.T) {
	// 创建模拟数据源
	total := 25
	mockSource := newMockDataSource(total) // 总共25条数据
	pageSize := 10

	// 创建基于页码的分页器
	pager, err := NewNumberPaged[TestData](mockSource.getNumberPageData, pageSize, 0)
	assert.NoError(t, err)
	assert.NotNil(t, pager)

	// 测试迭代器
	t.Run("Iterator", func(t *testing.T) {
		count := 0
		for pager.Next() {
			count++
			item := pager.Current()
			assert.Equal(t, count, item.ID)
			assert.Equal(t, fmt.Sprintf("test-%d", count), item.Name)
		}
		assert.Equal(t, total, count)
		assert.False(t, pager.HasMore())
		assert.Equal(t, 25, pager.Total())
		assert.NoError(t, pager.Err())
	})

	t.Run("manual fetch next page", func(t *testing.T) {
		count := 0
		hasMore := true
		currentPage := 1
		for hasMore {
			pager, err := NewNumberPaged[TestData](mockSource.getNumberPageData, pageSize, currentPage)
			assert.Nil(t, err)
			hasMore = pager.HasMore()
			count += len(pager.Items())
			currentPage++
		}
		assert.Equal(t, total, count)
		assert.False(t, pager.HasMore())
		assert.Equal(t, total, pager.Total())
		assert.NoError(t, pager.Err())
	})
}

func TestTokenPaged(t *testing.T) {
	total := 25
	mockSource := newMockDataSource(total) // 总共25条数据
	pageSize := 10

	pager, err := NewTokenPaged[TestData](mockSource.getTokenPageData, pageSize, nil)
	assert.NoError(t, err)
	assert.NotNil(t, pager)

	t.Run("iterator", func(t *testing.T) {
		count := 0
		for pager.Next() {
			count++
			item := pager.Current()
			assert.Equal(t, count, item.ID)
			assert.Equal(t, fmt.Sprintf("test-%d", count), item.Name)
		}
		assert.Equal(t, total, count)
		assert.False(t, pager.HasMore())
		assert.Equal(t, total, pager.Total())
		assert.NoError(t, pager.Err())
	})

	t.Run("manual fetch next page", func(t *testing.T) {
		count := 0
		hasMore := true
		var nextID *string
		for hasMore {
			pager, err := NewTokenPaged[TestData](mockSource.getTokenPageData, pageSize, nextID)
			assert.Nil(t, err)
			hasMore = pager.HasMore()
			count += len(pager.Items())
			for _, item := range pager.Items() {
				if item != nil {
					nextID = ptr(strconv.Itoa(item.ID))
				}
			}
		}
		assert.Equal(t, total, count)
		assert.False(t, pager.HasMore())
		assert.Equal(t, total, pager.Total())
		assert.NoError(t, pager.Err())
	})
}

func TestPagerError(t *testing.T) {
	// 测试错误情况
	errorFetcher := func(request *PageRequest) (*PageResponse[TestData], error) {
		return nil, fmt.Errorf("mock error")
	}

	// 测试基于页码的分页器错误处理
	t.Run("NumberPaged Error", func(t *testing.T) {
		pager, err := NewNumberPaged[TestData](errorFetcher, 10, 1)
		assert.Error(t, err)
		assert.Nil(t, pager)
	})

	// 测试基于令牌的分页器错误处理
	t.Run("TokenPaged Error", func(t *testing.T) {
		pager, err := NewTokenPaged[TestData](errorFetcher, 10, nil)
		assert.Error(t, err)
		assert.Nil(t, pager)
	})
}

func TestEmptyPage(t *testing.T) {
	// 创建空数据源
	emptySource := newMockDataSource(0)

	// 测试基于页码的空分页
	t.Run("Empty NumberPaged", func(t *testing.T) {
		pager, err := NewNumberPaged[TestData](emptySource.getNumberPageData, 10, 1)
		assert.NoError(t, err)
		assert.NotNil(t, pager)
		assert.False(t, pager.Next())
		assert.Equal(t, 0, pager.Total())
		assert.False(t, pager.HasMore())
		assert.NoError(t, pager.Err())
	})

	// 测试基于令牌的空分页
	t.Run("Empty TokenPaged", func(t *testing.T) {
		pager, err := NewTokenPaged[TestData](emptySource.getTokenPageData, 10, nil)
		assert.NoError(t, err)
		assert.NotNil(t, pager)
		assert.False(t, pager.Next())
		assert.Equal(t, 0, pager.Total())
		assert.False(t, pager.HasMore())
		assert.NoError(t, pager.Err())
	})
}

func TestInvalidPageSize(t *testing.T) {
	mockSource := newMockDataSource(25)

	// 测试基于页码的无效页大小
	t.Run("Invalid PageSize NumberPaged", func(t *testing.T) {
		pager, err := NewNumberPaged[TestData](mockSource.getNumberPageData, 0, 1)
		assert.NoError(t, err)
		assert.NotNil(t, pager)
		assert.True(t, pager.Next())
		assert.Equal(t, 25, pager.Total())
	})

	// 测试基于令牌的无效页大小
	t.Run("Invalid PageSize TokenPaged", func(t *testing.T) {
		pager, err := NewTokenPaged[TestData](mockSource.getTokenPageData, 0, nil)
		assert.NoError(t, err)
		assert.NotNil(t, pager)
		assert.True(t, pager.Next())
		assert.Equal(t, 25, pager.Total())
	})
}
