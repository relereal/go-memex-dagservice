package dagservice

import (
	"context"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/linking"
	"github.com/ipld/go-ipld-prime/node/basicnode"

	"github.com/relereal/go-memex-blockservice"
)

type Dagservice struct {
	blockservice *blockservice.Blockservice
	lsys         linking.LinkSystem
	lp           datamodel.LinkPrototype
}

func NewDagservice(bs *blockservice.Blockservice, lsys linking.LinkSystem, lp datamodel.LinkPrototype) *Dagservice {
	return &Dagservice{
		blockservice: bs,
		lsys:         lsys,
		lp:           lp,
	}
}

func (ds *Dagservice) Store(ctx context.Context, node ipld.Node) (datamodel.Link, error) {
	link, err := ds.lsys.Store(
		linking.LinkContext{},
		ds.lp,
		node,
	)
	return link, err
}

func (ds *Dagservice) Load(ctx context.Context, link datamodel.Link) (ipld.Node, error) {
	np := basicnode.Prototype.Any
	node, err := ds.lsys.Load(
		linking.LinkContext{},
		link,
		np,
	)
	return node, err
}
