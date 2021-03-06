package main

import (
	"database/sql"
	"fmt"
	_ "github.com/sebastienboisard/godb2"
	"time"
)

func db2Roution(conn string) {
	for {
		time.Sleep(time.Minute)
		cnt, err := readCount(conn)
		if err == nil {
			cntqueue <- cnt
		}
	}
}

func readCount(conn string) (*dbcount, error) {
	//return nil, fmt.Errorf("test")

	cnt := &dbcount{}
	cnt.T = time.Now()

	if conn == "" {
		return nil, fmt.Errorf("connect string is empty")
	}

	db, err := sql.Open("db2-cli", conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	date := time.Now().Format(`0102`)
	newyear := fmt.Sprintf("%d-01-01", time.Now().Year())

	sql := fmt.Sprintf(`select * from (select count(*)  from TBL_SMSendTask),
		(select count(*)  from TBL_SMRESULT_%s where sendtime>'%s' and (Recv_Status = '0' or Recv_Status = '2')) ,
		(select count(*)  from TBL_SMRESULT_%s where sendtime>'%s' and (Recv_Status != '0' and Recv_Status != '2'));`,
		date, newyear, date, newyear)
	st, err := db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer st.Close()

	err = st.QueryRow().Scan(&cnt.Sending, &cnt.Ok, &cnt.Fail)
	if err != nil {
		return nil, err
	}

	return cnt, nil
}
