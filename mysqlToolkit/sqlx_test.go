package mysqlToolkit

import (
	"testing"
	"time"
)

type User struct {
	ID         int       `db:"id"`
	Email      string    `db:"email"`
	Password   string    `db:"password"`
	UpdateTime time.Time `db:"update_time"`
	CreateTime time.Time
}

func TestSqlConn_QueryRow(t *testing.T) {
	c, e := NewMySQL("root:12345671@tcp(localhost:3306)/galaxy")
	if e != nil {
		t.Error(e)
		return
	}

	v := User{}
	e = c.QueryRowPartial(&v, `select * from galaxy_user`)
	if e != nil {
		t.Error(e)
		return
	}
	if v.Email != `email` {
		t.Error("v.Email is not `email` , but ", v.Email)
		return
	}
	t.Log(v.ID)
	t.Log(v.UpdateTime)
}
