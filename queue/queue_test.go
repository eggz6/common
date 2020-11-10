package queue

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

// QueueSuite ...
type QueueSuite struct {
	suite.Suite
	mockErr error
}

// TestUnitClient go test 执行入口
func TestQueue(t *testing.T) {
	suite.Run(t, new(QueueSuite))
}

// SetupTest 执行用例前初始化
func (s *QueueSuite) SetupTest() {
	s.mockErr = errors.New("mock error")
}

func (s *QueueSuite) Test_newQueue() {
	for _, tt := range []struct {
		name string
		want queue
	}{
		{
			name: "empty",
			want: queue{size: 10, elements: make([]*element, 10)},
		},
	} {
		s.Run(tt.name, func() {
			res := newQueue(10)
			s.Equal(tt.want, *res)
		})
	}

}

func (s *QueueSuite) Test_Empty() {
	for _, tt := range []struct {
		name string
		want bool
	}{
		{
			name: "empty",
			want: true,
		},
	} {
		s.Run(tt.name, func() {
			que := newQueue(10)
			res := que.Empty()
			s.Equal(tt.want, res)
		})
	}

}

func (s *QueueSuite) Test_Full() {
	for _, tt := range []struct {
		name string
		want bool
	}{
		{
			name: "full",
			want: false,
		},
	} {
		s.Run(tt.name, func() {
			que := newQueue(10)
			res := que.Full()
			s.Equal(tt.want, res)
		})
	}

}

func (s *QueueSuite) Test_PopToEmpty() {
	//dep := func() *queue {
	//return newQueue(10)
	//}

	for _, tt := range []struct {
		name  string
		want  interface{}
		ok    bool
		empty bool
		dep   func() *queue
	}{
		{
			name: "pop_empty",
			dep: func() *queue {
				que := newQueue(10)
				for i := 0; i < 11; i++ {
					que.Append(i)
				}

				return que
			},
			want:  []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			empty: true,
		},
	} {
		s.Run(tt.name, func() {
			que := tt.dep()
			data := make([]int, 0)
			res, ok := que.Pop()
			for ok {
				i, _ := res.(int)
				data = append(data, i)
				res, ok = que.Pop()
			}

			s.Equal(tt.empty, que.Empty())
			s.Equal(tt.want, data)
		})
	}
}

func (s *QueueSuite) Test_Pop() {
	//dep := func() *queue {
	//return newQueue(10)
	//}

	for _, tt := range []struct {
		name string
		want interface{}
		ok   bool
		dep  func() *queue
	}{
		{
			name: "empty",
			dep:  func() *queue { return newQueue(10) },
			ok:   false,
		},
		{
			name: "pop_one",
			dep: func() *queue {
				que := newQueue(10)
				for i := 0; i < 11; i++ {
					que.Append(i)
				}

				return que
			},
			want: 0,
			ok:   true,
		},
	} {
		s.Run(tt.name, func() {
			que := tt.dep()
			res, ok := que.Pop()
			s.Equal(tt.ok, ok)
			s.Equal(tt.want, res)
		})
	}
}

func (s *QueueSuite) Test_Append() {
	//dep := func() *queue {
	//return newQueue(10)
	//}

	for _, tt := range []struct {
		name string
		want []int
		full bool
		dep  func() *queue
	}{
		{
			name: "append_full",
			dep: func() *queue {
				que := newQueue(10)
				for i := 0; i < 11; i++ {
					que.Append(i)
				}

				return que
			},
			want: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			full: true,
		},
	} {
		s.Run(tt.name, func() {
			que := tt.dep()
			full := que.Full()
			s.Equal(tt.full, full)
			res := []int{}
			for _, ele := range que.elements {
				i, _ := ele.data.(int)
				res = append(res, i)
			}
			s.Equal(tt.want, res)
		})
	}
}

func (s *QueueSuite) Test_Snake() {
	//dep := func() *queue {
	//return newQueue(10)
	//}

	for _, tt := range []struct {
		name string
		want []int
		ele  []int
		pop  int
		dep  func() *queue
	}{
		{
			name: "append_full",
			dep: func() *queue {
				que := newQueue(10)
				for i := 0; i < 10; i++ {
					que.Append(i)
				}

				return que
			},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 11},
			ele:  []int{11, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			pop:  0,
		},
	} {
		s.Run(tt.name, func() {
			que := tt.dep()
			i, _ := que.Pop()
			s.Equal(tt.pop, i)
			que.Append(11)

			data := make([]int, 0)
			res, ok := que.Pop()
			for ok {
				i, _ := res.(int)
				data = append(data, i)
				res, ok = que.Pop()
			}

			eles := []int{}
			for _, ele := range que.elements {
				i, _ := ele.data.(int)
				eles = append(eles, i)
			}

			s.Equal(tt.ele, eles)
			s.Equal(tt.want, data)
		})
	}
}
