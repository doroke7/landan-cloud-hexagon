package usecase

import (
	"errors"
	"testing"

	bootstrap "example/bootstrap"
	domain "example/internal/domain"
	helper "example/internal/helper"
	utility "example/internal/utility"
	pkg "example/pkg"
)

type fakeAdminUserRepository struct {
	gotName string
	retUser *domain.AdminUser
	retErr  error
}

func (f *fakeAdminUserRepository) ShowOneByName(name string) (*domain.AdminUser, error) {
	f.gotName = name
	return f.retUser, f.retErr
}

func (f *fakeAdminUserRepository) ShowOneById(id uint) (*domain.AdminUser, error) {
	return f.retUser, f.retErr
}

func (f *fakeAdminUserRepository) AddOne(oAdminUser *domain.AdminUserValue) (bool, error) {
	return true, nil
}

func (f *fakeAdminUserRepository) ShowOnesByWheresWithOrdersLimit(aWheres []*pkg.Where, aOrders []*pkg.Order, oLimit *pkg.Limit) ([]*domain.AdminUser, error) {
	return nil, nil
}

func (f *fakeAdminUserRepository) TotalByWheres(aWheres []*pkg.Where) (uint64, error) {
	return 0, nil
}

func newTestAbstractUsecase() *AbstractUsecase {
	oAbstractHelper := helper.NewAbstractHelper()
	return NewAbstractUsecase(
		helper.NewAesHelper(oAbstractHelper),
		helper.NewJwtHelper(oAbstractHelper),
	)
}

func TestAuthenticatorUsecase_SignIn(t *testing.T) {
	t.Run("name 是空字串時，不應該呼叫 repository，直接回傳錯誤", func(t *testing.T) {
		oRepo := &fakeAdminUserRepository{}
		oUsecase := NewAuthenticatorUsecase(oRepo, newTestAbstractUsecase())

		_, err := oUsecase.SignIn("", "secret", "test-secret")

		if err == nil {
			t.Fatal("預期回傳錯誤，卻沒有錯誤")
		}
		if oRepo.gotName != "" {
			t.Fatalf("不該呼叫到 repository，got gotName=%q", oRepo.gotName)
		}
	})

	t.Run("password 是空字串時，不應該呼叫 repository，直接回傳錯誤", func(t *testing.T) {
		oRepo := &fakeAdminUserRepository{}
		oUsecase := NewAuthenticatorUsecase(oRepo, newTestAbstractUsecase())

		_, err := oUsecase.SignIn("tom", "", "test-secret")

		if err == nil {
			t.Fatal("預期回傳錯誤，卻沒有錯誤")
		}
		if oRepo.gotName != "" {
			t.Fatalf("不該呼叫到 repository，got gotName=%q", oRepo.gotName)
		}
	})

	t.Run("密碼不正確時，回傳錯誤", func(t *testing.T) {
		oExpected := &domain.AdminUser{Id: 1, Name: "tom", Password: "wrong-hash"}
		oRepo := &fakeAdminUserRepository{retUser: oExpected}
		oUsecase := NewAuthenticatorUsecase(oRepo, newTestAbstractUsecase())

		sAuthorization, err := oUsecase.SignIn("tom", "secret", "test-secret")

		if err == nil {
			t.Fatal("預期回傳密碼錯誤，卻沒有錯誤")
		}
		if sAuthorization != "" {
			t.Fatalf("預期回傳空字串，got %q", sAuthorization)
		}
	})

	t.Run("找到帳號且密碼正確時，回傳 JWT", func(t *testing.T) {
		sPassword := "secret"
		sHashed := utility.Md5(sPassword + bootstrap.CONFIG.TABLE.ADMIN_USER.PASSWORD)
		oExpected := &domain.AdminUser{Id: 1, Name: "tom", Password: sHashed}
		oRepo := &fakeAdminUserRepository{retUser: oExpected}
		oUsecase := NewAuthenticatorUsecase(oRepo, newTestAbstractUsecase())

		sAuthorization, err := oUsecase.SignIn("tom", sPassword, "test-secret")

		if err != nil {
			t.Fatalf("預期沒有錯誤，卻收到: %v", err)
		}
		if oRepo.gotName != "tom" {
			t.Fatalf("repository 收到的 name 不對，got %q", oRepo.gotName)
		}
		if sAuthorization == "" {
			t.Fatal("預期回傳非空的 JWT")
		}
	})

	t.Run("repository 沒回傳錯誤、但也沒找到帳號（nil, nil）時，要轉成明確的 not found 錯誤", func(t *testing.T) {
		oRepo := &fakeAdminUserRepository{retUser: nil, retErr: nil}
		oUsecase := NewAuthenticatorUsecase(oRepo, newTestAbstractUsecase())

		sAuthorization, err := oUsecase.SignIn("nobody", "secret", "test-secret")

		if sAuthorization != "" {
			t.Fatalf("預期回傳空字串，got %q", sAuthorization)
		}
		if err == nil {
			t.Fatal("預期回傳 not found 錯誤，卻沒有錯誤")
		}
	})

	t.Run("repository 失敗時，usecase 要原樣把錯誤往上傳", func(t *testing.T) {
		oExpectedErr := errors.New("db down")
		oRepo := &fakeAdminUserRepository{retErr: oExpectedErr}
		oUsecase := NewAuthenticatorUsecase(oRepo, newTestAbstractUsecase())

		_, err := oUsecase.SignIn("tom", "secret", "test-secret")

		if !errors.Is(err, oExpectedErr) {
			t.Fatalf("預期錯誤被原樣傳回，got: %v", err)
		}
	})
}
