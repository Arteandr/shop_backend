package service

import (
	"context"
	"mime/multipart"
	"shop_backend/internal/models"
	"shop_backend/internal/repository"
	"shop_backend/pkg/auth"
	"shop_backend/pkg/hash"
	"shop_backend/pkg/mail"
	"time"
)

type Images interface {
	Upload(ctx context.Context, images []*multipart.FileHeader) error
	GetAll(ctx context.Context) ([]models.Image, error)
	Exist(ctx context.Context, imageId int) (bool, error)
	Delete(ctx context.Context, imagesId []int) error
}

type Colors interface {
	GetById(ctx context.Context, colorId int) (models.Color, error)
	GetAll(ctx context.Context) ([]models.Color, error)
	Create(ctx context.Context, name, hex string, price float64) (int, error)
	Update(ctx context.Context, id int, name, hex string, price float64) error
	Delete(ctx context.Context, colorsId []int) error
	DeleteFromItems(ctx context.Context, colorId int) error
	AddToItems(ctx context.Context, colorId int) error
	Exist(ctx context.Context, colorId int) (bool, error)
}

type Categories interface {
	GetAll(ctx context.Context) ([]models.Category, error)
	GetById(ctx context.Context, categoryId int) (models.Category, error)
	Create(ctx context.Context, name string, imageId int) (int, error)
	Delete(ctx context.Context, categoryId int) error
	Update(ctx context.Context, categoryId int, name string, imageId int) error
	Exist(ctx context.Context, colorId int) (bool, error)
}

type Items interface {
	Create(ctx context.Context, item models.Item) (models.Item, error)
	Update(ctx context.Context, item models.Item) error
	LinkColors(ctx context.Context, itemId int, colors []models.Color) error
	LinkTags(ctx context.Context, itemId int, tags []models.Tag) error
	LinkImages(ctx context.Context, itemId int, images []models.Image) error
	GetAll(ctx context.Context, sortOptions models.SortOptions) ([]models.Item, error)
	GetNew(ctx context.Context) ([]models.Item, error)
	GetById(ctx context.Context, itemId int) (models.Item, error)
	GetBySku(ctx context.Context, sku string) (models.Item, error)
	GetByCategory(ctx context.Context, categoryId int) ([]models.Item, error)
	GetByTag(ctx context.Context, tag string) ([]models.Item, error)
	Delete(ctx context.Context, itemsId []int) error
	Exist(ctx context.Context, itemId int) (bool, error)
}

type Users interface {
	SignUp(ctx context.Context, email, login, password string) (models.User, error)
	SignIn(ctx context.Context, findBy, login, password string) (models.Tokens, error)
	Logout(ctx context.Context, userId int) error
	GetMe(ctx context.Context, userId int) (models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	RefreshTokens(ctx context.Context, refreshToken string) (models.Tokens, error)
	UpdateEmail(ctx context.Context, userId int, email string) error
	UpdatePassword(ctx context.Context, userId int, oldPassword, newPassword string) error
	UpdateInfo(ctx context.Context, userId int, login, firstName, lastName, companyName, phoneCode, phoneNumber string) error
	UpdateAddress(ctx context.Context, userId int, different bool, invoiceAddress models.Address, shippingAddress models.Address) error
	DeleteMe(ctx context.Context, userId int) error
	IsCompleted(ctx context.Context, userId int) (bool, error)
	CompleteVerify(ctx context.Context, token string) error
	SendVerify(ctx context.Context, userId int) error
}

type Delivery interface {
	Create(ctx context.Context, delivery models.Delivery) (int, error)
	GetById(ctx context.Context, deliveryId int) (models.Delivery, error)
	GetAll(ctx context.Context) ([]models.Delivery, error)
	Update(ctx context.Context, delivery models.Delivery) error
	Delete(ctx context.Context, deliveryId int) error
	Exist(ctx context.Context, deliveryId int) (bool, error)
}

type Orders interface {
	Create(ctx context.Context, order models.Order) (int, error)
	Exist(ctx context.Context, orderId int) (bool, error)
	GetAll(ctx context.Context) ([]models.ServiceOrder, error)
	GetAllByUserId(ctx context.Context, userId int) ([]models.ServiceOrder, error)
	GetAllStatuses(ctx context.Context) ([]models.OrderStatus, error)
	UpdateStatus(ctx context.Context, orderId, statusId int) error
	Delete(ctx context.Context, orderId int) error
}

type Mails interface {
	CreateVerify(ctx context.Context, userId int, login, email string) error
	CompleteVerify(ctx context.Context, token string) (int, error)
}

type Services struct {
	Users      Users
	Items      Items
	Categories Categories
	Colors     Colors
	Images     Images
	Delivery   Delivery
	Orders     Orders
	Mails      Mails
}

type ServicesDeps struct {
	Repos           *repository.Repositories
	Hasher          hash.PasswordHasher
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	MailSender      mail.Sender
}

func NewServices(deps ServicesDeps) *Services {
	images := NewImagesService(deps.Repos.Images)
	categories := NewCategoriesService(deps.Repos.Categories, images)
	colors := NewColorsService(deps.Repos.Colors)
	mails := NewMailsService(deps.Repos.Mails, deps.MailSender)
	users := NewUsersService(UsersServiceDeps{
		repo:            deps.Repos.Users,
		hasher:          deps.Hasher,
		tokenManager:    deps.TokenManager,
		accessTokenTTL:  deps.AccessTokenTTL,
		refreshTokenTTL: deps.RefreshTokenTTL,
		mailsService:    mails,
	})
	delivery := NewDeliveryService(deps.Repos.Delivery)
	items := NewItemsService(deps.Repos.Items, categories, colors, images)

	return &Services{
		Items:      items,
		Categories: categories,
		Colors:     colors,
		Images:     images,
		Orders:     NewOrdersService(deps.Repos.Orders, users, delivery, items, colors),
		Delivery:   delivery,
		Users:      users,
		Mails:      mails,
	}
}
