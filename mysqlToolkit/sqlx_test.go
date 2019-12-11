package mysqlToolkit

import (
	"testing"
	"time"
)

type User struct {
	ID         int    `db:"id"`
	Email      string `db:"email"`
	Password   string `db:"password"`
	SendTime   time.Time
	UpdateTime time.Time `db:"update_time"`
	CreateTime time.Time
}

func TestSqlConn_QueryRow(t *testing.T) {
	c, e := NewMySQL("root:12345671@tcp(localhost:3306)/galaxy")
	if e != nil {
		t.Error(e)
		return
	}

	vs := []*User{}
	e = c.QueryRows(&vs, `select * from galaxy_user`)
	if e != nil {
		t.Error(e)
		return
	}
	t.Log(vs[0])
}
