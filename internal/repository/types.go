package repository

type IStorage interface {
	Get(id int) (string, bool)
	Set(link string) int
}
