package nft

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/skip2/go-qrcode"
)

// TicketImageGenerator generates ticket images with QR codes
type TicketImageGenerator struct {
	templateWidth  int
	templateHeight int
	bgColor        color.Color
	primaryColor   color.Color
	font           *truetype.Font
}

// NewTicketImageGenerator creates a new ticket image generator
func NewTicketImageGenerator(fontData []byte) (*TicketImageGenerator, error) {
	f, err := freetype.ParseFont(fontData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}

	return &TicketImageGenerator{
		templateWidth:  1200,
		templateHeight: 600,
		bgColor:        color.RGBA{R: 245, G: 247, B: 250, A: 255}, // Light gray background
		primaryColor:   color.RGBA{R: 34, G: 139, B: 34, A: 255},   // Forest green
		font:           f,
	}, nil
}

// TicketData contains all information needed to generate a ticket image
type TicketData struct {
	TokenID            uint64
	BookingID          string
	ExperienceTitle    string
	Location           string
	Date               time.Time
	ValidUntil         time.Time
	ParticipantName    string
	ExperienceType     string
	CarbonOffset       int
	ConservationImpact string
}

// GenerateTicketImage creates a ticket image with QR code
func (g *TicketImageGenerator) GenerateTicketImage(data TicketData) ([]byte, error) {
	// Create base image
	img := image.NewRGBA(image.Rect(0, 0, g.templateWidth, g.templateHeight))

	// Fill background
	draw.Draw(img, img.Bounds(), &image.Uniform{g.bgColor}, image.Point{}, draw.Src)

	// Draw border
	g.drawBorder(img)

	// Add header section
	g.drawHeader(img, data)

	// Generate and add QR code
	qrData := fmt.Sprintf("bogowi://redeem/%d/%s", data.TokenID, data.BookingID)
	qrCode, err := qrcode.New(qrData, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}
	qrCode.BackgroundColor = color.White
	qrCode.ForegroundColor = color.Black
	qrImage := qrCode.Image(200)

	// Position QR code on the right side
	qrPosition := image.Point{X: 950, Y: 200}
	draw.Draw(img, image.Rect(qrPosition.X, qrPosition.Y, qrPosition.X+200, qrPosition.Y+200),
		qrImage, image.Point{}, draw.Over)

	// Add ticket details on the left side
	g.drawTicketDetails(img, data)

	// Add footer with conservation impact
	g.drawFooter(img, data)

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return buf.Bytes(), nil
}

// drawBorder adds a decorative border to the ticket
func (g *TicketImageGenerator) drawBorder(img *image.RGBA) {
	// Draw thick green border
	borderColor := g.primaryColor
	borderWidth := 10

	// Top border
	for y := 0; y < borderWidth; y++ {
		for x := 0; x < g.templateWidth; x++ {
			img.Set(x, y, borderColor)
		}
	}

	// Bottom border
	for y := g.templateHeight - borderWidth; y < g.templateHeight; y++ {
		for x := 0; x < g.templateWidth; x++ {
			img.Set(x, y, borderColor)
		}
	}

	// Left border
	for x := 0; x < borderWidth; x++ {
		for y := 0; y < g.templateHeight; y++ {
			img.Set(x, y, borderColor)
		}
	}

	// Right border
	for x := g.templateWidth - borderWidth; x < g.templateWidth; x++ {
		for y := 0; y < g.templateHeight; y++ {
			img.Set(x, y, borderColor)
		}
	}
}

// drawHeader adds the header section with BOGOWI branding
func (g *TicketImageGenerator) drawHeader(img *image.RGBA, data TicketData) {
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(g.font)
	c.SetFontSize(48)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(&image.Uniform{g.primaryColor})

	// Draw BOGOWI title
	pt := freetype.Pt(50, 80)
	c.DrawString("BOGOWI", pt)

	// Draw subtitle
	c.SetFontSize(24)
	pt = freetype.Pt(50, 120)
	c.DrawString("Eco-Experience Ticket", pt)

	// Draw ticket number
	c.SetFontSize(18)
	c.SetSrc(&image.Uniform{color.RGBA{R: 100, G: 100, B: 100, A: 255}})
	ticketNumber := fmt.Sprintf("Ticket #%d", data.TokenID)
	pt = freetype.Pt(50, 150)
	c.DrawString(ticketNumber, pt)
}

// drawTicketDetails adds the main ticket information
func (g *TicketImageGenerator) drawTicketDetails(img *image.RGBA, data TicketData) {
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(g.font)
	c.SetClip(img.Bounds())
	c.SetDst(img)

	// Experience title
	c.SetFontSize(32)
	c.SetSrc(&image.Uniform{color.Black})
	pt := freetype.Pt(50, 220)
	c.DrawString(data.ExperienceTitle, pt)

	// Details
	c.SetFontSize(20)
	c.SetSrc(&image.Uniform{color.RGBA{R: 60, G: 60, B: 60, A: 255}})

	details := []string{
		fmt.Sprintf("Location: %s", data.Location),
		fmt.Sprintf("Date: %s", data.Date.Format("January 2, 2006")),
		fmt.Sprintf("Valid Until: %s", data.ValidUntil.Format("January 2, 2006")),
		fmt.Sprintf("Experience Type: %s", data.ExperienceType),
	}

	y := 260
	for _, detail := range details {
		pt = freetype.Pt(50, y)
		c.DrawString(detail, pt)
		y += 35
	}

	// Booking reference
	c.SetFontSize(16)
	bookingRef := "Booking Reference: "
	if len(data.BookingID) > 16 {
		bookingRef += data.BookingID[:16] + "..."
	} else {
		bookingRef += data.BookingID
	}
	pt = freetype.Pt(50, 420)
	c.DrawString(bookingRef, pt)
}

// drawFooter adds conservation impact information
func (g *TicketImageGenerator) drawFooter(img *image.RGBA, data TicketData) {
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(g.font)
	c.SetClip(img.Bounds())
	c.SetDst(img)

	// Conservation impact section
	c.SetFontSize(18)
	c.SetSrc(&image.Uniform{g.primaryColor})
	pt := freetype.Pt(50, 500)
	c.DrawString("Conservation Impact:", pt)

	c.SetFontSize(16)
	c.SetSrc(&image.Uniform{color.RGBA{R: 60, G: 60, B: 60, A: 255}})

	// Carbon offset
	carbonText := fmt.Sprintf("• Carbon Offset: %d kg CO2", data.CarbonOffset)
	pt = freetype.Pt(50, 530)
	c.DrawString(carbonText, pt)

	// Conservation message
	pt = freetype.Pt(50, 555)
	c.DrawString(fmt.Sprintf("• %s", data.ConservationImpact), pt)
}

// Helper function to draw text with word wrapping
func (g *TicketImageGenerator) drawWrappedText(
	c *freetype.Context,
	text string,
	x, y int,
	maxWidth int,
) {
	// Simple word wrapping implementation
	words := bytes.Fields([]byte(text))
	var line []byte
	currentY := y

	for _, word := range words {
		testLine := append(append(line, ' '), word...)

		// Measure line width (simplified - would need proper implementation)
		if len(testLine) > maxWidth/10 { // Rough estimation
			// Draw current line
			pt := freetype.Pt(x, currentY)
			c.DrawString(string(line), pt)

			// Start new line
			line = word
			currentY += 25
		} else {
			line = testLine
		}
	}

	// Draw remaining text
	if len(line) > 0 {
		pt := freetype.Pt(x, currentY)
		c.DrawString(string(line), pt)
	}
}
