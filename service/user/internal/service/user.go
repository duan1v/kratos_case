package service

import (
	"context"
	"time"
	v1 "user/api/user/v1"
	"user/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	// 也可以继承v1.UserServer,自定义UnimplementedUserServer
	v1.UnimplementedUserServer

	uc  *biz.UserUsecase
	log *log.Helper
}

// NewUserService new a greeter service.
func NewUserService(uc *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{uc: uc, log: log.NewHelper(logger)}
}

// CreateUser create a user
func (u *UserService) CreateUser(ctx context.Context, req *v1.CreateUserInfo) (*v1.UserInfoResponse, error) {
	user, err := u.uc.Create(ctx, &biz.User{
		Mobile:   req.Mobile,
		Password: req.Password,
		Nickname: req.Nickname,
	})
	if err != nil {
		return nil, err
	}
	var birthday uint64
	// 需要先判断在赋值,否则会出现空指针异常
	if user.Birthday != nil {
		birthday = uint64(user.Birthday.Unix())
	}
	userInfoRsp := v1.UserInfoResponse{
		Id:       user.ID,
		Mobile:   user.Mobile,
		Password: user.Password,
		Nickname: user.Nickname,
		Gender:   user.Gender,
		Role:     int32(user.Role),
		Birthday: birthday,
	}

	return &userInfoRsp, nil
}

// GetUserList .
func (u *UserService) GetUserList(ctx context.Context, req *v1.PageInfo) (*v1.UserListResponse, error) {
	list, total, err := u.uc.List(ctx, int(req.Pn), int(req.Psize))
	if err != nil {
		return nil, err
	}
	rsp := &v1.UserListResponse{}
	rsp.Total = int32(total)

	for _, user := range list {
		userInfoRsp := UserResponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}

	return rsp, nil
}

func UserResponse(user *biz.User) v1.UserInfoResponse {
	userInfoRsp := v1.UserInfoResponse{
		Id:       user.ID,
		Mobile:   user.Mobile,
		Password: user.Password,
		Nickname: user.Nickname,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoRsp.Birthday = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

// GetUserByMobile .
func (u *UserService) GetUserByMobile(ctx context.Context, req *v1.MobileRequest) (*v1.UserInfoResponse, error) {
	user, err := u.uc.UserByMobile(ctx, req.Mobile)
	if err != nil {
		return nil, err
	}
	rsp := UserResponse(user)
	return &rsp, nil
}

// UpdateUser .
func (u *UserService) UpdateUser(ctx context.Context, req *v1.UpdateUserInfo) (*emptypb.Empty, error) {
	birthDay := time.Unix(int64(req.Birthday), 0)
	user, err := u.uc.UpdateUser(ctx, &biz.User{
		ID:       req.Id,
		Gender:   req.Gender,
		Birthday: &birthDay,
		Nickname: req.Nickname,
	})

	if err != nil {
		return nil, err
	}

	if user == false {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// CheckPassword .
func (u *UserService) CheckPassword(ctx context.Context, req *v1.PasswordCheckInfo) (*v1.CheckResponse, error) {
	check, err := u.uc.CheckPassword(ctx, req.Password, req.EncryptedPassword)
	if err != nil {
		return nil, err
	}
	return &v1.CheckResponse{Success: check}, nil
}

// GetUserById .
func (u *UserService) GetUserById(ctx context.Context, req *v1.IdRequest) (*v1.UserInfoResponse, error) {
	user, err := u.uc.UserById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	rsp := UserResponse(user)
	return &rsp, nil
}
