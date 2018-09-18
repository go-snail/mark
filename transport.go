package mark


type Writer interface {
	write(*event) error
}