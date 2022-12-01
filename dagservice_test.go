package dagservice

import (
	"context"
	"os"
	"testing"

	cid "github.com/ipfs/go-cid"
	_ "github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/fluent/qp"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	mc "github.com/multiformats/go-multicodec"
	mh "github.com/multiformats/go-multihash"

	blockservice "github.com/relereal/go-memex-blockservice"
	blockstore "github.com/relereal/go-memex-blockstore"
	datastore "github.com/relereal/go-sqlite-datastore"
)

func clearDatastore(dstore *datastore.Datastore) {
	dstore.CloseDb()
	os.RemoveAll("test")
}

func getDagservice() (*Dagservice, *datastore.Datastore) {
	// get datastore
	os.Mkdir("test", 0777)
	os.Remove("test/testdb.db")
	dstore := datastore.NewDatastore("test/testdb.db", "keystore")
	dstore.Connect()

	// get blockstore
	bstore := blockstore.NewBlockstore(dstore)

	// get blockservice
	bservice := blockservice.NewBlockservice(bstore)

	// get linksystem
	lsys := cidlink.DefaultLinkSystem()
	lsys.SetReadStorage(bservice)
	lsys.SetWriteStorage(bservice)

	// get linkprototype
	lp := cidlink.LinkPrototype{Prefix: cid.Prefix{
		Version:  1,
		Codec:    uint64(mc.DagCbor),
		MhType:   mh.SHA2_256,
		MhLength: -1, // default length
	}}

	// get dagservice
	dservice := NewDagservice(bservice, lsys, lp)

	return dservice, dstore
}

func TestDagservice(t *testing.T) {
	dservice, dstore := getDagservice()
	defer clearDatastore(dstore)

	np := basicnode.Prototype.Any
	node, err := qp.BuildMap(np, -1, func(ma datamodel.MapAssembler) {
		qp.MapEntry(ma, "testkey", qp.String("testvalue"))
	})

	link, err := dservice.Store(context.Background(), node)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	node, err = dservice.Load(context.Background(), link)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}
	testValue, err := node.LookupByString("testkey")
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	testValueStr, err := testValue.AsString()
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}
	if testValueStr != "testvalue" {
		t.Errorf("Expected testvalue, got %s", testValueStr)
	}
}
