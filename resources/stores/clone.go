package stores

import (
	pb "github.com/Saser/strecku/api/v1"
	"google.golang.org/protobuf/proto"
)

func Clone(store *pb.Store) *pb.Store {
	return proto.Clone(store).(*pb.Store)
}
