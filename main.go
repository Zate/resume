package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {

	app := fiber.New(fiber.Config{
		ServerHeader:          "CombatWombat",
		ProxyHeader:           "X-Forwarded-*",
		DisableStartupMessage: false,
	})
	app.Use(logger.New(logger.Config{
		Format: "{\"pid\": \"${pid}\",\"time\": \"${time}\",\"referer\": \"${referer}\",\"protocol\": \"${protocol}\",\"ip\": \"${ip}\",\"ips\": \"${ips}\",\"host\": \"${host}\",\"method\": \"${method}\",\"path\": \"${path}\",\"url\": \"${url}\",\"ua\": \"${ua}\",\"latency\": \"${latency}\",\"status\": \"${status}\",\"bytesSent\": \"${bytesSent}\",\"bytesReceived\": \"${bytesReceived}\"}\n",
		// Format:     "${time}\n",
		TimeFormat: "2006-01-02T15:04:05.9999Z07:00",
		TimeZone:   "UTC",
		Output:     os.Stdout,
	}))
	app.Use(recover.New())

	app.Static("/", "site")

	log.Fatal(app.Listen(":3000"))

}

// PDFGenerator should generate a PDF of the site.
func PDFGenerator() {

	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	// Set global options
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set("Portrait")
	pdfg.Grayscale.Set(false)

	// Create a new input page from an URL
	page := wkhtmltopdf.NewPage("https://themes.3rdwavemedia.com/demo/devresume/")

	// Set options for this page
	page.FooterRight.Set("[page]")
	page.FooterFontSize.Set(10)
	page.Zoom.Set(0.95)

	// Add to document
	pdfg.AddPage(page)

	// Create PDF document in internal buffer
	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	// Write buffer contents to file on disk
	err = pdfg.WriteFile("./simplesample.pdf")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
	// Output: Done
}
