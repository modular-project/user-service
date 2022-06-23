package cookiestorage

// import (
// 	"os"
// 	"sync"
// 	"time"

// 	"users-service/storage"

// 	"github.com/wader/gormstore/v2"
// )

// var (
// 	once sync.Once
// 	db   *gormstore.Store
// 	key  []byte
// )

// func DB() *gormstore.Store {
// 	return db
// }

// func NewCookieStore() error {
// 	var err error
// 	once.Do(func() {
// 		loadKey()
// 		db = gormstore.New(storage.DB(), key)
// 		quit := make(chan struct{})
// 		go db.PeriodicCleanup(24*time.Hour, quit)
// 	})
// 	return err
// }

// func loadKey() {
// 	keys, _ := os.LookupEnv("RGE_COOKIE_KEY")
// 	key = []byte(keys)
// }
