package biz_test

import (
	"time"
	"user/internal/biz"
	"user/internal/data"
	"user/internal/mock/mrepo"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _ = Describe("UserUsecase", func() {
	var userCase *biz.UserUsecase
	var mUserRepo *mrepo.MockUserRepo

	BeforeEach(func() {
		mUserRepo = mrepo.NewMockUserRepo(ctl)
		userCase = biz.NewUserUsecase(mUserRepo, nil)
	})

	It("Create", func() {
		birthday := time.Unix(int64(693629981), 0)
		info := &biz.User{
			ID:       1,
			Mobile:   "13803881388",
			Password: "admin123456",
			Nickname: "aliliin",
			Role:     1,
			Birthday: &birthday,
		}
		mUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(info, nil)
		l, err := userCase.Create(ctx, info)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(err).ToNot(HaveOccurred())
		Ω(l.ID).To(Equal(int64(1)))
		Ω(l.Mobile).To(Equal("13803881388"))
	})
})

var _ = Describe("User", func() {
	var ro biz.UserRepo
	var uD *biz.User
	BeforeEach(func() {
		mysqlSource := "dywily:q123456we@tcp(127.0.0.1:3306)/shop_user?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := gorm.Open(mysql.Open(mysqlSource), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		dataData, _, _ := data.NewData(nil, nil, db, nil)
		ro = data.NewUserRepo(dataData, nil)
		// 这里你可以不引入外部组装好的数据，可以在这里直接写
		uD = &biz.User{
			ID:       1,
			Mobile:   "13509876789",
			Password: "admin123",
			Nickname: "aliliin",
			Role:     1,
			Birthday: &time.Time{},
		}
	})
	// 设置 It 块来添加单个规格
	// It("CreateUser", func() {
	// 	u, err := ro.CreateUser(ctx, uD)
	// 	Ω(err).ShouldNot(HaveOccurred())
	// 	// 组装的数据 mobile 为 13509876789
	// 	Ω(u.Mobile).Should(Equal("13509876789")) // 手机号应该为创建的时候写入的手机号
	// })
	// 设置 It 块来添加单个规格
	It("ListUser", func() {
		user, total, err := ro.ListUser(ctx, 1, 10)
		Ω(err).ShouldNot(HaveOccurred()) // 获取列表不应该出现错误
		Ω(user).ShouldNot(BeEmpty())     // 结果不应该为空
		Ω(total).Should(Equal(1))        // 总数应该为 1，因为上面只创建了一条
		Ω(len(user)).Should(Equal(1))
		Ω(user[0].Mobile).Should(Equal("13509876789"))
	})
	// 设置 It 块来添加单个规格
	It("UpdateUser", func() {
		birthDay := time.Unix(int64(693646426), 0)
		uD.Nickname = "gyl"
		uD.Birthday = &birthDay
		uD.Gender = "female"
		user, err := ro.UpdateUser(ctx, uD)
		Ω(err).ShouldNot(HaveOccurred()) // 更新不应该出现错误
		Ω(user).Should(BeTrue())         // 结果应该为 true
	})

	It("CheckPassword", func() {
		p1 := "admin123"
		encryptedPassword := "$2a$10$ARy./L.TGyshthF6.VnmD.5fPbyANFzqVB6OKaQAoE68LU6s13wqi"
		password, err := ro.CheckPassword(ctx, p1, encryptedPassword)
		Ω(err).ShouldNot(HaveOccurred()) // 密码验证通过
		Ω(password).Should(BeTrue())     // 结果应该为true

		encryptedPassword1 := "$2a$10$ARy./L.TGyshthF6.VnmD.5fPbyANFzqVB6OKafdoE68LU6s13wqi"
		password1, err := ro.CheckPassword(ctx, p1, encryptedPassword1)
		if err != nil {
			return
		}
		Ω(err).ShouldNot(HaveOccurred())
		Ω(password1).Should(BeFalse()) // 密码验证不通过
	})
})
