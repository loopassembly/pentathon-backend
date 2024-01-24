package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/loopassembly/pentathon-backend/initializers"
	"github.com/loopassembly/pentathon-backend/models"
	"github.com/loopassembly/pentathon-backend/utils"
	"google.golang.org/api/sheets/v4"
)

type SoloData struct {
	Name            string `json:"Name"`
	WhatsAppNo      string `json:"WhatsApp No"`
	SRMISTEmail     string `json:"SRMIST Email"`
	RegistrationNo  string `json:"Registration No"`
	YearOfStudy     string `json:"Year of Study"`
	Department      string `json:"Department"`
	FaName          string `json:"FA Name"`
	Section         string `json:"Section"`
	Time            string `json:"time"`
	SRMISTEmailForm string `json:"SRMIST e-mail"`
}

func SoloController(c *fiber.Ctx) error {

	bodyFn := c.Request().Body

	var soloData SoloData
	if err := json.Unmarshal(bodyFn(), &soloData); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"result": "error",
			"error":  "Error parsing JSON data",
		})
	}

	err := postToGoogleAppsScript(soloData)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"result": "error",
			"error":  "Error posting to Google Apps Script",
		})
	}

	err = postToDiscordWebhook(soloData)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"result": "error",
			"error":  "Error posting to Discord webhook",
		})
	}

	return c.JSON(fiber.Map{
		"result":  "success",
		"message": fmt.Sprintf("Solo registration successful for %s", soloData.Name),
	})
}

func postToGoogleAppsScript(data SoloData) error {
	url := "https-://script.google.com/macros/s/AKfycbzVYt0n-KBCrL9kN_d9LQNcu4kkgiCMsd4vPjSJLHVNZ9zDaWGISmb30-zh0sgWlS_FCw/exec"
	method := "POST"

	payloadJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	payload := strings.NewReader(string(payloadJSON))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	return err
}

func postToDiscordWebhook(data SoloData) error {
	url := "https://discord.com/api/webhooks/your_discord_webhook_url"
	method := "POST"

	payload := strings.NewReader(`{"content": "Solo registration\nName: ` + data.Name + `\nEmail: ` + data.SRMISTEmail + `"}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	return err
}

func ReadDataHandler(c *fiber.Ctx) error {
	config, _ := initializers.LoadConfig(".")
	sheetsService, err := utils.GetSheetsService()
	if err != nil {
		log.Println("Error getting Google Sheets service:", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	resp, err := sheetsService.Spreadsheets.Values.Get(config.SpreadsheetID, config.Testsheet).Context(c.Context()).Do()
	if err != nil {
		log.Println("Error reading data from Google Sheets:", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	data, _ := json.Marshal(resp.Values)
	return c.Status(http.StatusOK).JSON(data)
}

func CreateDataHandler(c *fiber.Ctx) error {
	config, _ := initializers.LoadConfig(".")
	sheetsService, err := utils.GetSheetsService()
	if err != nil {
		log.Println("Error getting Google Sheets service:", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	var requestData utils.RequestData
	if err := c.BodyParser(&requestData); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(http.StatusBadRequest).SendString("Bad Request")
	}

	values := sheets.ValueRange{Values: requestData.Values}
	_, err = sheetsService.Spreadsheets.Values.Append(config.SpreadsheetID, config.Testsheet, &values).ValueInputOption("RAW").Do()
	if err != nil {
		log.Println("Error creating data in Google Sheets:", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.SendStatus(http.StatusCreated)
}
func SoloDataHandler(c *fiber.Ctx) error {
	config, _ := initializers.LoadConfig(".")
	sheetsService, err := utils.GetSheetsService()
	if err != nil {
		log.Println("Error getting Google Sheets service:", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	var requestData *utils.RequestData
	if err := c.BodyParser(&requestData); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(http.StatusBadRequest).SendString("Bad Request")
	}

	if requestData == nil || len(requestData.Values) == 0 {
		return c.Status(http.StatusBadRequest).SendString("Invalid or empty request data")
	}

	
	currentTime := time.Now().Format(time.RFC3339)

	
	for i := range requestData.Values {
		requestData.Values[i] = append(requestData.Values[i], currentTime)
	}

	userEmail := models.Email{
		Email: requestData.Values[0][2].(string),
	}

	if initializers.DB == nil {
		return c.Status(http.StatusInternalServerError).SendString("Database is not initialized")
	}

	log.Println("Before GORM Create")
	result := initializers.DB.Create(&userEmail)
	log.Println("After GORM Create")

	if result == nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Result is nil"})
	}

	if result.Error != nil {
		log.Println("Error creating data in the database")

		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "User with that email already exists"})
		}

		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}

	log.Println("Before Google Sheets API Call")
	values := sheets.ValueRange{Values: requestData.Values}
	_, err = sheetsService.Spreadsheets.Values.Append(config.SpreadsheetID, config.Solosheet, &values).ValueInputOption("RAW").Do()
	log.Println("After Google Sheets API Call")

	if err != nil {
		log.Println("Error creating data in Google Sheets:", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	// return c.SendStatus(http.StatusCreated)
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Solo registration successful",
	})
}

func TeamDataHandler(c *fiber.Ctx) error {
	config, _ := initializers.LoadConfig(".")
	sheetsService, err := utils.GetSheetsService()
	if err != nil {
		log.Println("Error getting Google Sheets service:", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	var requestData *utils.RequestData
	if err := c.BodyParser(&requestData); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(http.StatusBadRequest).SendString("Bad Request")
	}

	if requestData == nil || len(requestData.Values) == 0 {
		return c.Status(http.StatusBadRequest).SendString("Invalid or empty request data")
	}

	// Get the current timestamp
	currentTime := time.Now().Format(time.RFC3339)

	// Append the timestamp to each user's data
	for i := range requestData.Values {
		requestData.Values[i] = append(requestData.Values[i], currentTime)
	}

	if initializers.DB == nil {
		return c.Status(http.StatusInternalServerError).SendString("Database is not initialized")
	}

	emailIndices := [4]int{3, 11, 19, 27}

	
	for i := range requestData.Values {
		
		if len(requestData.Values[i]) > emailIndices[0] {
			
			for j := range emailIndices {
				
				if len(requestData.Values[i]) > emailIndices[j] {
					userEmail := models.Email{
						Email: requestData.Values[i][emailIndices[j]].(string),
					}
					
					fmt.Println("Email:", userEmail.Email)

					
					result := initializers.DB.Create(&userEmail)
					

					if result == nil {
						return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Result is nil"})
					}

					if result.Error != nil {
						log.Println("Error creating data in the database")

						if strings.Contains(result.Error.Error(), "UNIQUE constraint failed") {
							return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "User with that email already exists"})
						}

						return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
					}
				}
			}
		}
	}

	log.Println("Before Google Sheets API Call")
	values := sheets.ValueRange{Values: requestData.Values}
	_, err = sheetsService.Spreadsheets.Values.Append(config.SpreadsheetID, config.Teamsheet, &values).ValueInputOption("RAW").Do()
	log.Println("After Google Sheets API Call")

	if err != nil {
		log.Println("Error creating data in Google Sheets:", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	// return c.SendStatus(http.StatusCreated).json(fiber.Map{"status": "success", "message": "Team registration successful"})
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Team registration successful",
	})
	
}

