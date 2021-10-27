package filereconciler

type Handler interface {
	Reconcile() error
	Cleanup() error
}
