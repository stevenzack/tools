package alg

type Queue struct {
	list []interface{}
}

func (q *Queue) Push(v interface{}) {
	q.list = append(q.list, v)
}

func (q *Queue) Pop() interface{} {
	if len(q.list) == 0 {
		return nil
	}
	v := q.list[0]
	q.list = q.list[1:]
	return v
}

func (q *Queue) Length() int {
	return len(q.list)
}
