package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// 定义返回数据结构体
type User struct {
	ID       int64
	Mobile   string
	Password string
	Nickname string
	Birthday *time.Time
	Gender   string
	Role     int
}

// 注意这一行新增的 mock 数据的命令 "user/internal/biz" 必须写全否则生成的测试文件引入包报错
//go:generate mockgen -destination=../mock/mrepo/user.go -package=mrepo "user/internal/biz" UserRepo
type UserRepo interface {
	CreateUser(context.Context, *User) (*User, error)
	ListUser(ctx context.Context, pageNum, pageSize int) ([]*User, int, error)
	UserByMobile(ctx context.Context, mobile string) (*User, error)
	GetUserById(ctx context.Context, id int64) (*User, error)
	UpdateUser(context.Context, *User) (bool, error)
	CheckPassword(ctx context.Context, password, encryptedPassword string) (bool, error)
}

type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *UserUsecase) Create(ctx context.Context, u *User) (*User, error) {
	return uc.repo.CreateUser(ctx, u)
}

func (uc *UserUsecase) List(ctx context.Context, pageNum, pageSize int) ([]*User, int, error) {
	return uc.repo.ListUser(ctx, pageNum, pageSize)
}

func (uc *UserUsecase) UserByMobile(ctx context.Context, mobile string) (*User, error) {
	return uc.repo.UserByMobile(ctx, mobile)
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, user *User) (bool, error) {
	return uc.repo.UpdateUser(ctx, user)
}

func (uc *UserUsecase) CheckPassword(ctx context.Context, password, encryptedPassword string) (bool, error) {
	return uc.repo.CheckPassword(ctx, password, encryptedPassword)
}

func (uc *UserUsecase) UserById(ctx context.Context, id int64) (*User, error) {
	return uc.repo.GetUserById(ctx, id)
}
