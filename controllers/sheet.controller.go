package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/loopassembly/pentathon-backend/utils"
	"github.com/loopassembly/pentathon-backend/initializers"
	"log"
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
	// Use c.Request().Body() directly, which returns a function that returns []byte
	bodyFn := c.Request().Body

	var soloData SoloData
	if err := json.Unmarshal(bodyFn(), &soloData); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"result": "error",
			"error":  "Error parsing JSON data",
		})
	}

	// Post the soloData to the Google Apps Script endpoint
	err := postToGoogleAppsScript(soloData)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"result": "error",
			"error":  "Error posting to Google Apps Script",
		})
	}

	// Post the soloData to the Discord webhook
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

	// Convert SoloData to JSON
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
	url := "https://discord.com/api/webhooks/your_discord_webhook_url"  // Replace with your actual Discord webhook URL
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

	var requestData utils.RequestData
	if err := c.BodyParser(&requestData); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(http.StatusBadRequest).SendString("Bad Request")
	}

	values := sheets.ValueRange{Values: requestData.Values}
	_, err = sheetsService.Spreadsheets.Values.Append(config.SpreadsheetID, config.Solosheet, &values).ValueInputOption("RAW").Do()
	if err != nil {
		log.Println("Error creating data in Google Sheets:", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.SendStatus(http.StatusCreated)
}

func TeamDataHandler(c *fiber.Ctx) error {
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
	_, err = sheetsService.Spreadsheets.Values.Append(config.SpreadsheetID, config.Teamsheet, &values).ValueInputOption("RAW").Do()
	if err != nil {
		log.Println("Error creating data in Google Sheets:", err)
		return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.SendStatus(http.StatusCreated)
}