package main

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
)

func boltWriteRoution() {
	for {
		select {
		case cnt := <-cntqueue:
			err := dbWriteDb2(cnt)
			if err != nil {
				fmt.Println(err)
			}
		case prsp := <-pingqueue:
			err := dbWritePing(prsp)
			if err != nil {
				fmt.Println(err)
			}
		}

	}

}

func dbWriteDb2(cnt dbcount) error {
	db, err := bolt.Open("sm.db", 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("db2"))
		if err != nil {
			return err
		}

		val := ""
		val, err = json.Marshal(cnt)
		if err != nil {
			return err
		}

		err = b.Put([]byte(cnt.T.Format(`20060102150405`)), val)
		return err
	})
}

func dbWritePing(ping *pingrsp) error {
	db, err := bolt.Open("sm.db", 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(fmt.Sprintf(`ping-%s`, ping.Tar)))
		if err != nil {
			return err
		}

		err = b.Put([]byte(ping.T.Format(`20060102150405`)), fmt.Sprintf(`%d`, ping.Ms))
		return err
	})
}
