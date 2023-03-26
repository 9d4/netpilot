package board

import (
	"github.com/9d4/netpilot/util"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	jww "github.com/spf13/jwalterweatherman"
	"gorm.io/gorm"
	"strconv"
)

type Handler struct {
	db        *gorm.DB
	validator *validator.Validate
}

func NewBoardHandler(db *gorm.DB) *Handler {
	db.AutoMigrate(&Board{})

	return &Handler{
		db:        db,
		validator: validator.New(),
	}
}

func (h *Handler) SetupRoutes(router fiber.Router) {
	// Create a sub-router with the "/boards" prefix
	boardRouter := router.Group("/boards")

	// Map routes to handlers
	boardRouter.Post("", h.CreateBoard)
	boardRouter.Get("", h.GetAllBoards)
	boardRouter.Get("/:uuid", h.GetBoardByUUID)
	boardRouter.Put("/:uuid", h.UpdateBoard)
	boardRouter.Delete("/:uuid", h.DeleteBoard)
}

func (h *Handler) CreateBoard(c *fiber.Ctx) error {
	// Parse input
	var input CreateBoardRequest
	if err := c.BodyParser(&input); err != nil {
		return err
	}

	// Generate UUID
	uuid := uuid.New().String()

	// Validate input using validator
	if err := h.validator.Struct(input); err != nil {
		return h.validationErrorResponse(c, err.(validator.ValidationErrors))
	}

	// Create board in database
	board := &Board{
		UUID:               uuid,
		Host:               input.Host,
		Port:               input.Port,
		InsecureSkipVerify: input.InsecureSkipVerify,
		User:               input.User,
		Password:           input.Password,
	}
	if err := h.db.Create(board).Error; err != nil {
		jww.TRACE.Println(err)
		return fiber.ErrInternalServerError
	}

	out := CreateBoardResponse{
		UUID: board.UUID,
	}

	c.Status(fiber.StatusCreated)
	return c.JSON(out)
}

func (h *Handler) GetAllBoards(c *fiber.Ctx) error {
	// Get query parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	// Calculate offset and limit based on query parameters
	offset := (page - 1) * pageSize
	limit := pageSize

	// Get 20 most recent boards from database
	var boards []*Board
	if err := h.db.Order("created_at desc").Offset(offset).Limit(limit).Find(&boards).Error; err != nil {
		return err
	}

	response := BoardsResponse{}
	for _, board := range boards {
		response.Boards = append(response.Boards, BoardsEachResponse{
			UUID: board.UUID,
			Name: board.Name,
			Host: board.Host,
			Port: board.Port,
		})
	}

	return c.JSON(response)
}

func (h *Handler) GetBoardByUUID(c *fiber.Ctx) error {
	// Get board UUID from path parameter
	uuid := c.Params("uuid")

	// Find board in database
	board := &Board{}
	if err := h.db.Where("uuid = ?", uuid).First(&board).Error; err != nil {
		return err
	}

	return c.JSON(board)
}

func (h *Handler) UpdateBoard(c *fiber.Ctx) error {
	// Get board UUID from path parameter
	uuidParam := c.Params("uuid")

	// Get board from db
	board := &Board{}
	if err := h.db.Where("uuid = ?", uuidParam).First(board).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.ErrNotFound
		}

		return fiber.ErrInternalServerError
	}

	// Parse input from request body
	req := &CreateBoardRequest{}
	if err := c.BodyParser(req); err != nil {
		return fiber.ErrBadRequest
	}

	// Validate input
	if err := h.validator.Struct(req); err != nil {
		return h.validationErrorResponse(c, err.(validator.ValidationErrors))
	}

	// Save changes to database
	updateBoard(board, req)
	if err := h.db.Save(&board).Error; err != nil {
		return err
	}

	return c.JSON(board)
}

func (h *Handler) DeleteBoard(c *fiber.Ctx) error {
	// Get board UUID from path parameter
	uuidParam := c.Params("uuid")

	// Find board in database
	board := &Board{}
	if err := h.db.Where("uuid = ?", uuidParam).First(&board).Error; err != nil {
		return fiber.ErrNotFound
	}

	// Delete board from database
	if err := h.db.Delete(&board).Error; err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) validationErrorResponse(c *fiber.Ctx, validationErrors validator.ValidationErrors) error {
	var errors []map[string]string
	for _, err := range validationErrors {
		errorMap := map[string]string{
			"error":   "validation",
			"message": "Validation error",
			"field":   util.ToSnakeCase(err.Field()),
			"tag":     err.ActualTag(),
		}
		errors = append(errors, errorMap)
	}
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"errors": errors,
	})
}

func (h *Handler) errorResponse(c *fiber.Ctx, statusCode int, error string, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"error":   error,
		"message": message,
	})
}
