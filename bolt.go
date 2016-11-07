package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"strings"
	"time"
)

func boltWriteRoution() {
	for {
		select {
		case cnt := <-cntqueue:
			err := dbWriteDb2(cnt)
			if err != nil {
				fmt.Println(err)
			}
		case prsp := <-pingReduceQueue:
			err := dbWritePing(prsp)
			if err != nil {
				fmt.Println(err)
			}
		}

	}

}

func dbWriteDb2(cnt dbcount) error {
	db, err := bolt.Open(boltdbname, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

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

func dbWritePing(ping *pingReduceRsp) error {
	db, err := bolt.Open(boltdbname, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(fmt.Sprintf(`ping-%s`, ping.Tar)))
		if err != nil {
			return err
		}

		val, err := json.Marshal(ping)
		if err != nil {
			return err
		}

		return b.Put([]byte(ping.T.Format(timekeyformat)), val)
	})
}

func dbReadPing(tar string, from, to string, limit int) (error, []string) {
	db, err := bolt.Open(boltdbname, 0600, nil)
	if err != nil {
		return err, nil
	}
	defer db.Close()

	rsps := make([]string, 0)

	err = db.View(func(tx *bolt.Tx) error {
		bt := tx.Bucket([]byte(fmt.Sprintf(`ping-%s`, tar)))
		if bt == nil {
			return fmt.Errorf("ip donot exist")
		}
		cr := bt.Cursor()
		if cr == nil {
			return fmt.Errorf("cursor return nil")
		}

		for k, v := cr.Seek([]byte(from)); k != nil && bytes.Compare(k, []byte(to)) <= 0; k, v = cr.Next() {
			limit--
			if limit < 0 {
				break
			}
			rsps = append(rsps, string(v))
		}
		return nil
	})
	return err, rsps
}
func dbReadDb2(from, to string, limit int) (error, []string) {
	db, err := bolt.Open(boltdbname, 0600, nil)
	if err != nil {
		return err, nil
	}
	defer db.Close()

	rsps := make([]string, 0)

	err = db.View(func(tx *bolt.Tx) error {
		bt := tx.Bucket([]byte(`db2`))
		if bt == nil {
			return fmt.Errorf("db2 bucket dont exist")
		}
		cr := bt.Cursor()
		if cr == nil {
			return fmt.Errorf("cursor return nil")
		}

		for k, v := cr.Seek([]byte(from)); k != nil && bytes.Compare(k, []byte(to)) <= 0; k, v = cr.Next() {
			limit--
			if limit < 0 {
				break
			}
			rsps = append(rsps, string(v))
		}
		return nil
	})
	return err, rsps
}

func dbReadPingLast(limit int) (error, []string) {
	db, err := bolt.Open(boltdbname, 0600, nil)
	if err != nil {
		return err, nil
	}
	defer db.Close()

	rsps := make([]string, 0)

	tend := time.Now().Add(time.Minute * -30).Format(timekeyformat)

	err = db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, bt *bolt.Bucket) error {
			if strings.HasPrefix(string(name), "ping-") && bt != nil {
				cr := bt.Cursor()
				if cr != nil {
					tmp := limit
					for k, v := cr.Last(); k != nil && bytes.Compare(k, []byte(tend)) >= 0; k, v = cr.Prev() {
						tmp--
						if tmp < 0 {
							break
						}

						rsps = append(rsps, string(v))
					}
				}
			}
			return nil
		})
	})
	return err, rsps
}
