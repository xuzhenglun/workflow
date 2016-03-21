package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xuzhenglun/workflow/core"
)

type Mysql struct {
	Db       *sql.DB
	mux      sync.Mutex
	Activity *core.Activity
}

type DBError struct {
	When time.Time
	What string
}

func (this DBError) Error() string {
	return fmt.Sprintf("%v: %v", this.When, this.What)
}

func NewMysql(URI string) Mysql {
	var db Mysql
	var err error
	if db.Db, err = sql.Open("mysql", URI); err != nil {
		log.Panic(err)
		panic(err)
	} else {
		log.Println("Connected to " + URI)
		err = db.CreateTable(`CREATE TABLE process (
		Id int(4) primary key auto_increment,
		Eid int(4) not null,
		Pass boolean not null,
		IsDone boolean,
		JustDone char(10)
	)`)
		if err != nil {
			log.Println(err)
		}
		return db
	}
}

func (this Mysql) AddRow(notnull []string, args map[string]interface{}) error {
	var SQL []interface{} //这个是一个坑，只有空接口切片才能传入空接口变长型参

	Prep := "INSERT events SET "

	for _, key := range notnull {
		if value, ok := args[key]; ok {
			if v, ok := value.(string); ok {
				SQL = append(SQL, v)
			} else {
				log.Println("ERROR: Wrong Args of key \"" + key + "\"")
				return DBError{When: time.Now(), What: "Wrong Args of key of \"" + key + "\""}
			}
		} else {
			return DBError{time.Now(), "Miss args which have to be not null"}
		}
	}

	Prep = Prep + strings.Join(notnull, "=? ,") + "=?"

	log.Println(Prep, SQL)

	this.mux.Lock()
	tx, err := this.Db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	stmt, err := tx.Prepare(Prep)
	if err != nil {
		log.Println(err)
		return err
	}

	r, err := stmt.Exec(SQL...)
	if err != nil {
		log.Println(err)
		return err
	}

	p, err := tx.Prepare("INSERT process SET EID=? ,IsDone=false ,Pass=false")
	if err != nil {
		log.Println(err)
	}

	id, _ := r.LastInsertId()
	ra := strconv.FormatInt(id, 10)
	if _, err := p.Exec(ra); err != nil {
		log.Println(err)
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		if err := tx.Rollback(); err != nil {
			log.Println(err)
			log.Println("Rollback Failed")
			return err
		}
		log.Println("AddRow Failed,Rollback Success")

		return err
	}
	this.mux.Unlock()
	return nil
}

func (this Mysql) ModifyRow(notnull []string, args map[string]interface{}) error {
	var SQL []interface{}

	Prep := "UPDATE events SET "
	for _, key := range notnull {
		if value, ok := args[key]; ok {
			log.Println(key)
			if v, ok := value.(string); ok {
				SQL = append(SQL, v)
			} else {
				log.Println("ERROR: Wrong Args of key \"" + key + "\"")
				return DBError{When: time.Now(), What: "Wrong Args of key of \"" + key + "\""}
			}
		} else {
			return DBError{time.Now(), "Miss args which have to be not null"}
		}
	}

	var id int

	if v, ok := args[":id"].(string); ok {
		if v, err := strconv.Atoi(v); err == nil {
			id = v
		} else {
			log.Println(err)
			return err
		}
	} else {
		log.Println("ERROR: Inv ID")
		return DBError{When: time.Now(), What: "Inv ID"}
	}

	ispass := false
	if _, ok := args["pass"]; ok {
		ispass = true
	}

	Prep1 := "UPDATE process SET Pass=? WHERE Id=" + strconv.Itoa(id)
	Prep2 := Prep + strings.Join(notnull, "=? ,") + "=? WHERE Id=(SELECT Eid FROM process WHERE Id=" + strconv.Itoa(id) + ")"

	this.mux.Lock()
	tx, err := this.Db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	stmt1, err := tx.Prepare(Prep1)
	if err != nil {
		log.Println(err)
		return err
	}

	stmt2, err := tx.Prepare(Prep2)
	if err != nil {
		log.Println(err)
		return err
	}

	var r sql.Result
	if r, err = stmt1.Exec(ispass); err != nil {
		log.Println(err)
		if err := tx.Rollback(); err != nil {
			log.Println(err)
			log.Println("Rollback Failed")
			return err
		}
		log.Println("ModifyRow Failed,Rollback Success")

		return err
	}
	if id, _ := r.RowsAffected(); id == 0 {
		return DBError{When: time.Now(), What: "Thare is not have this row"}
	}
	if r, err = stmt2.Exec(SQL...); err != nil {
		log.Println(err)
		if err := tx.Rollback(); err != nil {
			log.Println(err)
			log.Println("Rollback Failed")
			return err
		}
		log.Println("ModifyRow Failed,Rollback Success")
		return err
	}
	if id, _ := r.RowsAffected(); id == 0 {
		return DBError{When: time.Now(), What: "Thare is not have this row"}
	}

	if err := tx.Commit(); err != nil {
		log.Println(err)
		return err
	}
	this.mux.Unlock()
	return nil
}

func (this Mysql) FindRow(id int) (string, error) {
	this.mux.Lock()

	SQL := "SELECT * FROM events,process WHERE events.Id=process.Eid AND process.Id="
	rows, err := this.Db.Query(SQL + strconv.Itoa(id))

	this.mux.Unlock()

	if err != nil {
		log.Println(err)
		return "", err
	} else {
		defer rows.Close()
	}

	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}

	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)

	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return "", err
	}
	log.Println(string(jsonData))
	return string(jsonData), nil
}

func (this Mysql) DeleteRow(id int) error {
	this.mux.Lock()
	defer this.mux.Unlock()
	tx, err := this.Db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	stmt1, err := tx.Prepare("DELETE FROM events WHERE Id=(SELECT Eid FROM process WHERE Id=?)")
	stmt2, err := tx.Prepare("DELETE FROM process WHERE Id=?")
	if err != nil {
		log.Println(err)
		return err
	} else {
		var r sql.Result
		if r, err = stmt1.Exec(id); err != nil {
			log.Println(err)
		}
		if id, _ := r.RowsAffected(); id == 0 {
			return DBError{When: time.Now(), What: "Not Exist"}
		}
		if r, err = stmt2.Exec(id); err != nil {
			log.Println(err)
		}
		if id, _ := r.RowsAffected(); id == 0 {
			return DBError{When: time.Now(), What: "Not Exist"}
		}

		err := tx.Commit()
		if err != nil {
			log.Println(err)
			if err := tx.Rollback(); err != nil {
				log.Println(err)
				log.Println("Rollback Failed")
			}
			log.Println("Rollback Success")
			return err
		} else {
			return nil
		}
	}
}

func (this Mysql) CreateTable(str string) error {
	this.mux.Lock()
	defer this.mux.Unlock()

	tx, err := this.Db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	stms, err := tx.Prepare(str)
	if err != nil {
		log.Println(err)
		return err
	}

	if _, err := stms.Exec(); err != nil {
		log.Println(err)
		return err
	}
	if err := tx.Commit(); err != nil {
		log.Println(err)
		if err := tx.Rollback(); err != nil {
			log.Println(err)
			log.Println("Rollback Failed")
			return err
		}
		log.Println("Rollback Success")
	}
	return nil
}
