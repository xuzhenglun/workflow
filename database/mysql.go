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
		Pass boolean default false not null,
		Done boolean default false not null,
		JustDone char(10)
	)`)
		if err != nil {
			log.Println(err)
		}
		return db
	}
}

func (this Mysql) AddRow(args ...map[string]string) error {

	SQL1 := `INSERT INTO events SET  `
	var arg1 []interface{}
	for k, v := range args[0] {
		SQL1 = SQL1 + k + " = ?" + " ,"
		arg1 = append(arg1, v)
	}
	SQL2 := `INSERT INTO process SET  `
	var arg2 []interface{}
	for k, v := range args[1] {
		SQL2 = SQL2 + k + " = ?" + " ,"
		arg2 = append(arg2, v)
	}

	this.mux.Lock()
	defer this.mux.Unlock()

	tx, err := this.Db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	stms, err := tx.Prepare(SQL1[:len(SQL1)-2])
	if err != nil {
		log.Println(err)
		return err
	}
	r, err := stms.Exec(arg1...)
	var row string
	if err == nil {
		r64, _ := r.LastInsertId()
		row = strconv.FormatInt(r64, 10)
	} else {
		log.Println(err)
		return err
	}

	SQL2 = SQL2 + "Eid = " + row
	stms, err = tx.Prepare(SQL2)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stms.Exec(arg2...)
	if err != nil {
		log.Println(err)
		return err
	}

	return Commit(tx)
}

func (this Mysql) ModifyRow(args ...map[string]string) error {
	log.Println(args[0], args[1])
	SQL1 := `UPDATE events SET `
	var arg1 []interface{}
	for k, v := range args[0] {
		SQL1 = SQL1 + k + " = ?" + " ,"
		switch v {
		case "true":
			arg1 = append(arg1, true)
		case "false":
			arg1 = append(arg1, false)
		default:
			arg1 = append(arg1, v)
		}
	}
	SQL2 := `UPDATE process SET `
	var arg2 []interface{}
	for k, v := range args[1] {
		if k != "id" {
			SQL2 = SQL2 + k + " = ?" + " ,"
			switch v {
			case "true":
				arg2 = append(arg2, true)
			case "false":
				arg2 = append(arg2, false)
			default:
				arg2 = append(arg2, v)
			}
		}
	}

	SQL1 = SQL1[:len(SQL1)-1] + "WHERE Id = (SELECT Eid FROM process WHERE Id = ?)"
	SQL2 = SQL2[:len(SQL2)-1] + "WHERE Id = ?"
	arg1 = append(arg1, args[1]["id"])
	arg2 = append(arg2, args[1]["id"])

	this.mux.Lock()
	defer this.mux.Unlock()

	tx, err := this.Db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	if len(arg1) > 1 {
		stms, err := tx.Prepare(SQL1)
		if err != nil {
			log.Println(SQL1, arg1)
			log.Println(err)
			return err
		}
		_, err = stms.Exec(arg1...)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	if len(arg2) > 1 {
		stms, err := tx.Prepare(SQL2)
		if err != nil {
			log.Println(SQL2, arg2)
			log.Println(err)
			return err
		}
		_, err = stms.Exec(arg2...)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return Commit(tx)
}

func (this Mysql) FindRow(id string, needArgs ...string) (string, error) {
	SQL := `SELECT ` + strings.Join(needArgs, ",") + ` FROM process,events WHERE process.Eid = events.Id AND process.Id = ?`
	this.mux.Lock()
	defer this.mux.Unlock()

	rows, err := this.Db.Query(SQL, id)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer rows.Close()
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
	fmt.Println(string(jsonData))
	return string(jsonData), nil
}

func (this Mysql) DeleteRow(id ...string) error {
	SQL1 := `DELETE FROM events WHERE Id = (SELECT Eid FROM process WHERE Id = ?)`
	SQL2 := `DELETE FROM process WHERE Id = ?`
	var ID []interface{}
	for _, v := range id {
		ID = append(ID, v)
	}

	this.mux.Lock()
	defer this.mux.Unlock()

	tx, err := this.Db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	stms, err := tx.Prepare(SQL1)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stms.Exec(ID...)
	if err != nil {
		log.Println(err)
		return err
	}
	stms, err = tx.Prepare(SQL2)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stms.Exec(ID...)
	if err != nil {
		log.Println(err)
		return err
	}

	return Commit(tx)
	return nil
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

func (this Mysql) GetJustDone(id string) (string, string, error) {
	SQL := "SELECT JustDone,Pass FROM events,process WHERE events.Id=process.Eid AND process.Id=" + id

	this.mux.Lock()
	defer this.mux.Unlock()
	rows, err := this.Db.Query(SQL)
	defer rows.Close()
	if err != nil {
		return "", "", err
	}

	var father, pass string
	for rows.Next() {
		err = rows.Scan(&father, &pass)
	}
	if err != nil {
		log.Println(err)
		return "", "", err
	} else {
		return father, pass, nil
	}
}

func Commit(tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		log.Println(err)
		if err := tx.Rollback(); err != nil {
			log.Println(err)
			return err
		}
		return err
	}
	return nil
}

func (this Mysql) GetList(action string, pass string) (string, error) {
	SQL := "SELECT process.Id FROM process,events WHERE process.Eid=events.Id AND process.Pass=" + pass + " AND JustDone = \"" + action + "\""
	this.mux.Lock()

	idrow, err := this.Db.Query(SQL)
	if err != nil {
		this.mux.Unlock()
		return "", err
	}

	var ids []string
	for idrow.Next() {
		var buff string
		idrow.Scan(&buff)
		ids = append(ids, buff)
	}

	this.mux.Unlock()

	jsonid, err := json.Marshal(ids)
	if err == nil {
		return string(jsonid), nil
	} else {
		log.Println(err)
		return "", err
	}
}
