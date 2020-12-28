package users

import (
	pb "github.com/Saser/strecku/api/v1"
	"google.golang.org/protobuf/proto"
)

func Clone(user *pb.User) *pb.User {
	return proto.Clone(user).(*pb.User)
}
