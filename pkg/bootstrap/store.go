package bootstrap

import (
	"os"

	"github.com/facebookgo/inject"
	"github.com/facebookgo/inmem"
	_ "github.com/lib/pq"
	"golang.ysitd.cloud/db"
)

func injectStore(graph *inject.Graph) {
	graph.Provide(
		&inject.Object{Value: db.NewOpener("postgres", os.Getenv("DB_URL"))},
		&inject.Object{Value: inmem.NewLocked(16)},
	)
}
