package middleware

import (
	"log"
	"pop_v1/models"

	"github.com/gofiber/fiber/v2"
)

func Nodecheck(c *fiber.Ctx) error {
	//get the request body
	client := models.Node{}
	if err := c.BodyParser(&client); err != nil {
		log.Println("Error fething data of client :", err)
		return c.Status(400).JSON(fiber.Map{
			"Error": err,
		})
	}
	return c.Next()
}
