// go:generate mockgen -source=db.go -destination=mock_db.go -package=db
package db

type DB interface {
	Get(key string) string
	Set(key string, value string)
}
