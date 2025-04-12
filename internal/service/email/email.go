package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/gomail.v2"
)

// EmailService handles email operations
type EmailService struct {
	config Config
}

// Config represents email service configuration
type Config struct {
	Host       string
	Port       int
	Username   string
	Password   string
	FromEmail  string
	FromName   string
	Domain     string // bsu.by
	Prefix     string // rct.
	TemplateDir string
	Debug      bool
}

// NewEmailService creates a new email service
func NewEmailService(config Config) *EmailService {
	// Set sensible defaults if not provided
	if config.Port == 0 {
		config.Port = 587 // Default to TLS port
	}
	if config.Domain == "" {
		config.Domain = "bsu.by"
	}
	if config.Prefix == "" {
		config.Prefix = "rct."
	}
	if config.TemplateDir == "" {
		config.TemplateDir = "templates/email"
	}

	return &EmailService{
		config: config,
	}
}

// FormatBSUEmail formats a BSU email address with the given login
func (s *EmailService) FormatBSUEmail(login string) string {
	// Remove any unsafe characters from login
	login = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '_' || r == '-' {
			return r
		}
		return -1
	}, login)

	return fmt.Sprintf("%s%s@%s", s.config.Prefix, login, s.config.Domain)
}

// SendVerificationEmail sends a verification code to a BSU email
func (s *EmailService) SendVerificationEmail(login, code string, expiresAt time.Time) error {
	to := s.FormatBSUEmail(login)
	subject := "Your Verification Code"
	
	// Data for the email template
	data := map[string]interface{}{
		"Login":       login,
		"Code":        code,
		"ExpiresAt":   expiresAt.Format("2006-01-02 15:04:05 MST"),
		"ExpiresIn":   fmt.Sprintf("%.0f minutes", time.Until(expiresAt).Minutes()),
		"CurrentYear": time.Now().Year(),
	}
	
	err := s.SendTemplatedEmail(to, subject, "verification", data)
	
	if err != nil {
		message := fmt.Sprintf(`Hello %s,asdfasfddsa

Your verification code is: 

📋 %s 📋
(copy and paste this code)

This code will expire in %s.

Best regards,
The Quiz Team

This email was sent to %s`, login, code, data["ExpiresIn"], to)
		
		return s.SendPlainTextEmail(to, subject, message)
	}
	
	return nil
}
// SendPlainTextEmail sends a plain text email
func (s *EmailService) SendPlainTextEmail(to, subject, body string) error {
	// Create message
	m := gomail.NewMessage()
	m.SetHeader("From", s.formatFrom())
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	
	if s.config.Debug {
		log.Printf("Debug mode: Would send email to %s with subject: %s", to, subject)
		log.Printf("Email body: %s", body)
		return nil
	}
	
	return s.send(m)
}

// SendHTMLEmail sends an HTML email with an optional plain text alternative
func (s *EmailService) SendHTMLEmail(to, subject, htmlBody, textBody string) error {
	// Create message
	m := gomail.NewMessage()
	m.SetHeader("From", s.formatFrom())
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	
	if textBody != "" {
		m.SetBody("text/plain", textBody)
		m.AddAlternative("text/html", htmlBody)
	} else {
		m.SetBody("text/html", htmlBody)
	}
	
	if s.config.Debug {
		log.Printf("Debug mode: Would send HTML email to %s with subject: %s", to, subject)
		log.Printf("HTML Body: %s", htmlBody)
		return nil
	}
	
	return s.send(m)
}

// SendTemplatedEmail sends an email using an HTML template
func (s *EmailService) SendTemplatedEmail(to, subject, templateName string, data map[string]interface{}) error {
	// Locate the template files
	htmlFile := filepath.Join(s.config.TemplateDir, templateName+".html")
	textFile := filepath.Join(s.config.TemplateDir, templateName+".txt")
	
	// Read and parse the HTML template
	htmlTmpl, err := template.ParseFiles(htmlFile)
	if err != nil {
		return fmt.Errorf("failed to parse HTML template: %w", err)
	}
	
	// Execute the HTML template
	var htmlBuf bytes.Buffer
	if err := htmlTmpl.Execute(&htmlBuf, data); err != nil {
		return fmt.Errorf("failed to execute HTML template: %w", err)
	}
	
	// Try to read and parse the text template if it exists
	var textBuf bytes.Buffer
	if _, err := os.Stat(textFile); err == nil {
		textTmpl, err := template.ParseFiles(textFile)
		if err == nil {
			if err := textTmpl.Execute(&textBuf, data); err == nil {
				// If we successfully generated both HTML and text, send both
				return s.SendHTMLEmail(to, subject, htmlBuf.String(), textBuf.String())
			}
		}
	}
	
	// If we don't have a text template or couldn't parse it, just send the HTML
	return s.SendHTMLEmail(to, subject, htmlBuf.String(), "")
}

// SendBulkEmail sends the same email to multiple recipients
func (s *EmailService) SendBulkEmail(logins []string, subject, body string, isHTML bool) error {
	for _, login := range logins {
		to := s.FormatBSUEmail(login)
		var err error
		
		if isHTML {
			err = s.SendHTMLEmail(to, subject, body, "")
		} else {
			err = s.SendPlainTextEmail(to, subject, body)
		}
		
		if err != nil {
			log.Printf("Failed to send email to %s: %v", to, err)
			// Continue with other recipients
		}
	}
	
	return nil
}

// formatFrom formats the From header
func (s *EmailService) formatFrom() string {
	if s.config.FromName != "" {
		return fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	}
	return s.config.FromEmail
}

// send delivers the email using the configured SMTP server
func (s *EmailService) send(m *gomail.Message) error {
	dialer := gomail.NewDialer(s.config.Host, s.config.Port, s.config.Username, s.config.Password)
	
	// Configure TLS
	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.config.Host,
	}
	
	return dialer.DialAndSend(m)
}


// SendMarkdownEmail sends an email using a Markdown template
func (s *EmailService) SendMarkdownEmail(to, subject, markdownTemplateName string, data map[string]interface{}) error {
	// Locate the template file
	mdFile := filepath.Join(s.config.TemplateDir, markdownTemplateName+".md")
	
	// Read the markdown template
	mdContent, err := ioutil.ReadFile(mdFile)
	if err != nil {
		return fmt.Errorf("failed to read markdown template: %w", err)
	}
	
	// Parse the markdown template with text/template first to fill in variables
	tmpl, err := template.New(markdownTemplateName).Parse(string(mdContent))
	if err != nil {
		return fmt.Errorf("failed to parse markdown template: %w", err)
	}
	
	// Execute the template with the provided data
	var mdBuf bytes.Buffer
	if err := tmpl.Execute(&mdBuf, data); err != nil {
		return fmt.Errorf("failed to execute markdown template: %w", err)
	}
	
	// Convert markdown to HTML
	markdown := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	
	var htmlBuf bytes.Buffer
	if err := markdown.Convert(mdBuf.Bytes(), &htmlBuf); err != nil {
		return fmt.Errorf("failed to convert markdown to HTML: %w", err)
	}
	
	// Wrap the HTML content in a proper HTML document with styling
	htmlContent := wrapHTMLContent(htmlBuf.String(), data)
	
	// Create a plain text version from the Markdown
	plainText := createPlainTextFromMarkdown(mdBuf.String())
	
	// Send the email with both HTML and plain text versions
	return s.SendHTMLEmail(to, subject, htmlContent, plainText)
}

// wrapHTMLContent wraps the converted HTML content in a full HTML document with styling
func wrapHTMLContent(htmlContent string, data map[string]interface{}) string {
	currentYear, _ := data["CurrentYear"].(int)
	if currentYear == 0 {
		currentYear = 2025 // Default if not provided
	}
	
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Email Notification</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .container {
            border: 1px solid #ddd;
            border-radius: 5px;
            padding: 20px;
            background-color: #f9f9f9;
        }
        h1, h2, h3, h4, h5, h6 {
            color: #444;
            margin-top: 1.2em;
            margin-bottom: 0.8em;
        }
        a {
            color: #4caf50;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
        pre, code {
            background-color: #f0f0f0;
            border-radius: 3px;
            padding: 2px 4px;
            font-family: monospace;
        }
        pre {
            padding: 10px;
            overflow-x: auto;
        }
        blockquote {
            border-left: 4px solid #ddd;
            padding-left: 10px;
            color: #666;
            margin-left: 0;
        }
        table {
            border-collapse: collapse;
            width: 100%%;
        }
        table, th, td {
            border: 1px solid #ddd;
            padding: 8px;
        }
        tr:nth-child(even) {
            background-color: #f2f2f2;
        }
        th {
            padding-top: 12px;
            padding-bottom: 12px;
            text-align: left;
            background-color: #4caf50;
            color: white;
        }
        .footer {
            margin-top: 20px;
            font-size: 12px;
            color: #777;
            text-align: center;
        }
    </style>
</head>
<body>
    <div class="container">
        %s
    </div>
    <div class="footer">
        <p>&copy; %d Quiz Platform. Все права защищены.</p>
    </div>
</body>
</html>
`, htmlContent, currentYear)
}

// createPlainTextFromMarkdown creates a simple plain text version from Markdown
// This is a basic implementation that handles common markdown elements
func createPlainTextFromMarkdown(markdown string) string {
	// Replace headers
	for i := 6; i >= 1; i-- {
		prefix := strings.Repeat("#", i) + " "
		markdown = strings.ReplaceAll(markdown, prefix, "")
	}
	
	// Replace links [text](url) with text (url)
	// This is a simplistic approach and might not handle all edge cases
	for {
		start := strings.Index(markdown, "[")
		if start == -1 {
			break
		}
		
		mid := strings.Index(markdown[start:], "](")
		if mid == -1 {
			break
		}
		mid += start
		
		end := strings.Index(markdown[mid:], ")")
		if end == -1 {
			break
		}
		end += mid
		
		text := markdown[start+1 : mid]
		url := markdown[mid+2 : end]
		
		replacement := text
		if url != text && !strings.Contains(text, url) {
			replacement = text + " (" + url + ")"
		}
		
		markdown = markdown[:start] + replacement + markdown[end+1:]
	}
	
	// Remove image markers
	markdown = strings.ReplaceAll(markdown, "![", "[")
	
	// Replace bold/italic markers
	markdown = strings.ReplaceAll(markdown, "**", "")
	markdown = strings.ReplaceAll(markdown, "__", "")
	markdown = strings.ReplaceAll(markdown, "*", "")
	markdown = strings.ReplaceAll(markdown, "_", "")
	
	// Replace horizontal rules
	markdown = strings.ReplaceAll(markdown, "---", "----------------------------")
	markdown = strings.ReplaceAll(markdown, "***", "----------------------------")
	
	// Replace code blocks with simple indentation
	// This is a simplistic approach
	markdown = strings.ReplaceAll(markdown, "```", "")
	markdown = strings.ReplaceAll(markdown, "`", "")
	
	// Handle blockquotes (simplistic)
	lines := strings.Split(markdown, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "> ") {
			lines[i] = "  " + strings.TrimPrefix(line, "> ")
		}
	}
	
	return strings.Join(lines, "\n")
}