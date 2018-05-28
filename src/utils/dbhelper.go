package utils

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
)

func PrintAllKeysInDB(db *leveldb.DB)  {
	log.Println("-----------------------")
	item := db.NewIterator(nil, nil)
	for item.Next() {
		key := item.Key()
		log.Println(string(key))
	}
	log.Println("-----------------------")
}
