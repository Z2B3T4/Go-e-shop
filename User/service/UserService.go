package service

import "project1/User/user"

type UserService struct {
	user.UnimplementedUserServiceServer
}
