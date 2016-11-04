package main

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
)

func boltWriteRoution() {
	db, err := bolt.Open("sm.db", 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

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

func dbWriteDb2(db *bolt.DB, cnt dbcount) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("db2"))
		if err != nil {
			return err
		}

		val, err := json.Marshal(cnt)
		if err != nil {
			return err
		}

		return b.Put([]byte(cnt.T.Format(timekeyformat)), val)
	})
}

func dbWritePing(db *bolt.DB, ping *pingrsp) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(fmt.Sprintf(`ping-%s`, ping.Tar)))
		if err != nil {
			return err
		}

		return b.Put([]byte(ping.T.Format(timekeyformat)), []byte(fmt.Sprintf(`%d`, ping.Ms)))
	})
}
