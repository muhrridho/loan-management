package delivery

import (
	"loan-management/internal/entity"
	"loan-management/internal/usecase"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) RegisterUser(ctx *fiber.Ctx) error {
	var payload entity.CreateUserPayload
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user := &entity.User{
		Email:     payload.Email,
		Name:      payload.Name,
		CreatedAt: time.Now(),
	}

	if err := h.userUsecase.RegisterUser(ctx.Context(), user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(user)
}

func (uh *UserHandler) GetAll(ctx *fiber.Ctx) error {
	users, err := uh.userUsecase.GetAll(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if users == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"data": []entity.User{}})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": users})
}

func (uh *UserHandler) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	user, err := uh.userUsecase.GetByID(ctx.Context(), int64(id))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": user})
}

func (uh *UserHandler) CheckUserDelinquentStatus(ctx *fiber.Ctx) error {
	id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}

	isUserDelinquent, err := uh.userUsecase.IsUserDelinquent(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": isUserDelinquent})

}
