package data

import (
	"context"
	"time"
	"user/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

// 定义数据表结构体
type User struct {
	ID          int64     `gorm:"primarykey"`
	Mobile      string    `gorm:"index:idx_mobile;unique;type:varchar(11) comment '手机号码，用户唯一标识';not null"`
	Password    string    `gorm:"type:varchar(100);not null "` // 用户密码的保存需要注意是否加密
	Nickname    string    `gorm:"type:varchar(25) comment '用户昵称'"`
	Birthday    uint64    `gorm:"type:datetime comment '出生日日期'"`
	Gender      string    `gorm:"column:gender;default:male;type:varchar(16) comment 'female:女,male:男'"`
	Role        int       `gorm:"column:role;default:1;type:int comment '1:普通用户，2:管理员'"`
	CreatedAt   time.Time `gorm:"column:created_time"`
	UpdatedAt   time.Time `gorm:"column:updated_time"`
	DeletedAt   gorm.DeletedAt
	IsDeletedAt bool
}
type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo . 这里需要注意，上面 data 文件 wire 注入的是此方法，方法名不要写错了
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// CreateUser .
func (r *userRepo) CreateUser(ctx context.Context, u *biz.User) (*biz.User, error) {
	var user User
	// 验证是否已经创建
	result := r.data.db.Where(&biz.User{Mobile: u.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	user.Mobile = u.Mobile
	user.Nickname = u.Nickname
	user.Password = encrypt(u.Password) // 密码加密
	res := r.data.db.Create(&user)
	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, res.Error.Error())
	}

	return &biz.User{
		ID:       user.ID,
		Mobile:   user.Mobile,
		Password: user.Password,
		Nickname: user.Nickname,
		Gender:   user.Gender,
		Role:     user.Role,
	}, nil
}

// Password encryption
func encrypt(pwd string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		panic("密码生成失败!")
	}
	return string(bytes)
}

// ListUser .
func (r *userRepo) ListUser(ctx context.Context, pageNum, pageSize int) ([]*biz.User, int, error) {
	var users []User
	result := r.data.db.Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	total := int(result.RowsAffected)
	r.data.db.Scopes(paginate(pageNum, pageSize)).Find(&users)
	rv := make([]*biz.User, 0)
	for _, u := range users {
		birthDay := time.Unix(int64(u.Birthday), 0)
		rv = append(rv, &biz.User{
			ID:       u.ID,
			Mobile:   u.Mobile,
			Password: u.Password,
			Nickname: u.Nickname,
			Gender:   u.Gender,
			Role:     u.Role,
			Birthday: &birthDay,
		})
	}
	return rv, total, nil
}

// paginate 分页
func paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// UserByMobile .
func (r *userRepo) UserByMobile(ctx context.Context, mobile string) (*biz.User, error) {
	var user User
	result := r.data.db.Where(&User{Mobile: mobile}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	re := modelToResponse(user)
	return &re, nil
}

func modelToResponse(u User) biz.User {
	birthDay := time.Unix(int64(u.Birthday), 0)
	return biz.User{
		ID:       u.ID,
		Mobile:   u.Mobile,
		Password: u.Password,
		Nickname: u.Nickname,
		Gender:   u.Gender,
		Role:     u.Role,
		Birthday: &birthDay,
	}
}

// UpdateUser .
func (r *userRepo) UpdateUser(ctx context.Context, user *biz.User) (bool, error) {
	var userInfo User
	result := r.data.db.Where(&User{ID: user.ID}).First(&userInfo)
	if result.RowsAffected == 0 {
		return false, status.Errorf(codes.NotFound, "用户不存在")
	}

	userInfo.Nickname = user.Nickname
	userInfo.Birthday = uint64(user.Birthday.Unix())
	userInfo.Gender = user.Gender

	res := r.data.db.Save(&userInfo)
	if res.Error != nil {
		return false, status.Errorf(codes.Internal, res.Error.Error())
	}

	return true, nil
}

// CheckPassword .
func (r *userRepo) CheckPassword(ctx context.Context, pwd, encryptedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(pwd))
	return err == nil, err
}

// GetUserById .
func (r *userRepo) GetUserById(ctx context.Context, Id int64) (*biz.User, error) {
	var user User
	result := r.data.db.Where(&User{ID: Id}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	re := modelToResponse(user)
	return &re, nil
}
