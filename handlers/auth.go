package handlers

import (
	"auth/models"
	"auth/utils"
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Register allows a user to register
func Register(c *fiber.Ctx) error {
	collection := utils.DB.Collection("users")
	var newUser models.User

	// Parse the request body
	if err := c.BodyParser(&newUser); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Validate required fields
	if newUser.Email == "" || newUser.Name == "" || newUser.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "All fields are required. Please provide email, name, and password.",
		})
	}

	// Check if user already exists
	var existingUser models.User
	err := collection.FindOne(context.Background(), bson.M{"email": newUser.Email}).Decode(&existingUser)
	if err == nil {
		// User already exists
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "User already exists"})
	} else if err != mongo.ErrNoDocuments {
		log.Printf("Error checking for existing user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
	}
	newUser.Password = string(hashedPassword)

	// Insert the user into the database
	_, err = collection.InsertOne(context.Background(), newUser)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not register user", "details": err.Error()})
	}

	log.Println("User registered successfully: ", newUser.Email) // Log successful registration
	return c.JSON(fiber.Map{"message": "User registered successfully"})
}


// Login allows a user to log in
func Login(c *fiber.Ctx) error {
	collection := utils.DB.Collection("users")
	var credentials models.User

	// Parse the request body
	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Find the user by email
	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"email": credentials.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// User not found
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
		}
		log.Printf("Error finding user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// Generate a JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	log.Println("User logged in successfully: ", user.Email) // Log successful login
	return c.JSON(fiber.Map{"token": tokenString})
}

