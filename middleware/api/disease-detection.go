package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AI Service Configuration
var AIServiceURL = getAIServiceURL()

func getAIServiceURL() string {
	url := os.Getenv("AI_SERVICE_URL")
	if url == "" {
		return "http://127.0.0.1:5000" // Default localhost AI service
	}
	return url
}

// Disease prediction structures
type DiseasePrediction struct {
	PredictedDisease string             `json:"predicted_disease"`
	Confidence       float64            `json:"confidence"`
	AllProbabilities map[string]float64 `json:"all_probabilities"`
	IsHealthy        bool               `json:"is_healthy"`
	Recommendation   string             `json:"recommendation"`
}

type PredictionResponse struct {
	Success    bool              `json:"success"`
	Prediction DiseasePrediction `json:"prediction"`
	Timestamp  string            `json:"timestamp"`
	Error      string            `json:"error,omitempty"`
}

// Health check for AI service
func AIHealthCheckHandler(c *fiber.Ctx) error {
	// Validate user authentication
	if err := ValidateCookie(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required"})
	}

	// Check AI service health
	resp, err := http.Get(AIServiceURL + "/health")
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error":   "AI service unavailable",
			"details": err.Error(),
		})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read AI service response"})
	}

	var healthStatus map[string]interface{}
	if err := json.Unmarshal(body, &healthStatus); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse AI service response"})
	}

	return c.JSON(healthStatus)
}

// Disease prediction handler
func PredictDiseaseHandler(c *fiber.Ctx) error {
	// Validate user authentication
	if err := ValidateCookie(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required"})
	}

	// Get uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No image file provided",
			"hint":  "Please upload an image file using the 'image' form field",
		})
	}

	// Validate file type
	if !isValidImageType(file.Header.Get("Content-Type")) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file type",
			"hint":  "Please upload a JPEG, PNG, or WebP image",
		})
	}

	// Validate file size (max 10MB)
	if file.Size > 10*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File too large",
			"hint":  "Please upload an image smaller than 10MB",
		})
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open uploaded file"})
	}
	defer src.Close()

	// Create a buffer to store the file content
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Create form file field
	fw, err := writer.CreateFormFile("image", file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create form data"})
	}

	// Copy file content
	if _, err := io.Copy(fw, src); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to copy file data"})
	}

	// Close the writer to finalize the form
	writer.Close()

	// Create HTTP request to AI service
	req, err := http.NewRequest("POST", AIServiceURL+"/predict", &buf)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create AI service request"})
	}

	// Set content type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Make request with timeout
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error":   "AI service request failed",
			"details": err.Error(),
		})
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read AI service response"})
	}

	// Parse response
	var predictionResp PredictionResponse
	if err := json.Unmarshal(body, &predictionResp); err != nil {
		fmt.Printf("Failed to parse AI response. Status: %d, Body: %s\n", resp.StatusCode, string(body))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to parse AI service response",
			"details": string(body),
		})
	}

	// Return the prediction result
	if resp.StatusCode != http.StatusOK || !predictionResp.Success {
		errorMsg := predictionResp.Error
		if errorMsg == "" {
			errorMsg = "Unknown error from AI service"
		}
		fmt.Printf("AI service error: %s (status: %d)\n", errorMsg, resp.StatusCode)
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"success": false,
			"error":   "AI service returned error",
			"details": errorMsg,
		})
	}

	// Log the prediction for monitoring
	fmt.Printf("✅ Disease prediction: %s (%.2f%% confidence) at %s\n",
		predictionResp.Prediction.PredictedDisease,
		predictionResp.Prediction.Confidence*100,
		predictionResp.Timestamp)

	return c.JSON(predictionResp)
}

// Utility function to validate image types
func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/webp",
	}

	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}

// Get disease information (educational content)
func GetDiseaseInfoHandler(c *fiber.Ctx) error {
	// Validate user authentication
	if err := ValidateCookie(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required"})
	}

	// Disease information database
	diseaseInfo := map[string]interface{}{
		"healthy": map[string]string{
			"name":        "Healthy",
			"description": "Normal, healthy chicken droppings",
			"symptoms":    "Firm, brown droppings with white urates",
			"treatment":   "Continue good hygiene and regular monitoring",
			"prevention":  "Maintain clean environment, fresh water, and proper nutrition",
		},
		"coccidiosis": map[string]string{
			"name":        "Coccidiosis",
			"description": "Parasitic infection affecting the intestinal tract",
			"symptoms":    "Bloody or watery droppings, lethargy, poor growth",
			"treatment":   "Anticoccidial medication, supportive care",
			"prevention":  "Keep environment dry, avoid overcrowding, use medicated feed",
		},
		"salmonella": map[string]string{
			"name":        "Salmonella",
			"description": "Bacterial infection that can affect digestive system",
			"symptoms":    "Watery droppings, decreased appetite, lethargy",
			"treatment":   "Antibiotics (consult veterinarian), supportive care",
			"prevention":  "Good hygiene, clean water, proper feed storage",
		},
		"e_coli": map[string]string{
			"name":        "E. Coli Infection",
			"description": "Bacterial infection often secondary to other conditions",
			"symptoms":    "Watery droppings, depression, poor growth",
			"treatment":   "Antibiotics, improve environmental conditions",
			"prevention":  "Clean environment, good ventilation, stress reduction",
		},
		"newcastle": map[string]string{
			"name":        "Newcastle Disease",
			"description": "Viral disease affecting respiratory and nervous systems",
			"symptoms":    "Greenish droppings, respiratory distress, neurological signs",
			"treatment":   "Supportive care (no specific treatment), isolation",
			"prevention":  "Vaccination, biosecurity measures, quarantine new birds",
		},
	}

	disease := c.Query("disease")
	if disease == "" {
		return c.JSON(fiber.Map{
			"diseases": diseaseInfo,
		})
	}

	info, exists := diseaseInfo[disease]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Disease information not found",
		})
	}

	return c.JSON(fiber.Map{
		"disease": disease,
		"info":    info,
	})
}
