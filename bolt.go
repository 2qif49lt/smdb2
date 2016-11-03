package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
)

func boltWriteRoution() {
	for {
		cnt := <-cntqueue

		err := dbWrite(cnt)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func dbWrite(cnt dbcount) error {
	db, err := bolt.Open("sm.db", 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		if err != nil {
			return err
		}

		key := ""
		key, err = json.Marshal(cnt)
		if err != nil {
			return err
		}

		err = b.Put([]byte(cnt.t.Format(`20060102150405`)), key)
		return err
	})
}
