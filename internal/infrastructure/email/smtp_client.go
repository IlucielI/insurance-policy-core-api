package email

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/wneessen/go-mail"
)

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

type SMTPClient struct {
	config *SMTPConfig
}

func NewSMTPClient(config *SMTPConfig) *SMTPClient {
	return &SMTPClient{config: config}
}

type EmailMessage struct {
	To      []string
	Subject string
	Body    string
	IsHTML  bool
}

func (s *SMTPClient) SendEmail(msg *EmailMessage) error {
	// Create new message
	m := mail.NewMsg()
	
	// Set sender
	if err := m.From(s.config.From); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	
	// Set recipients
	if err := m.To(msg.To...); err != nil {
		return fmt.Errorf("failed to set recipients: %w", err)
	}
	
	// Set subject
	m.Subject(msg.Subject)
	
	// Set body
	if msg.IsHTML {
		m.SetBodyString(mail.TypeTextHTML, msg.Body)
	} else {
		m.SetBodyString(mail.TypeTextPlain, msg.Body)
	}
	
	// Create SMTP client
	client, err := mail.NewClient(s.config.Host,
		mail.WithPort(s.config.Port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(s.config.User),
		mail.WithPassword(s.config.Password),
		mail.WithTLSPolicy(mail.TLSMandatory),
		mail.WithTLSConfig(&tls.Config{InsecureSkipVerify: false}),
		mail.WithTimeout(10*time.Second),
	)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	
	// Send email
	if err := client.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	
	return nil
}

// SendWelcomeEmail sends welcome email to newly registered user
func (s *SMTPClient) SendWelcomeEmail(to, fullName string) error {
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; background-color: #f9f9f9; }
		.footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>Welcome to Our Insurance Platform!</h1>
		</div>
		<div class="content">
			<h2>Hi %s,</h2>
			<p>Thank you for registering with us! We're excited to have you on board.</p>
			<p>You can now:</p>
			<ul>
				<li>Browse our insurance products</li>
				<li>Submit applications online</li>
				<li>Track your policies</li>
				<li>File and manage claims</li>
			</ul>
			<p>If you have any questions, feel free to reach out to our support team.</p>
			<p>Best regards,<br>The Insurance Platform Team</p>
		</div>
		<div class="footer">
			<p>&copy; 2026 Insurance Platform. All rights reserved.</p>
		</div>
	</div>
</body>
</html>
	`, fullName)

	return s.SendEmail(&EmailMessage{
		To:      []string{to},
		Subject: "Welcome to Our Insurance Platform",
		Body:    body,
		IsHTML:  true,
	})
}

// SendPolicyIssuedEmail sends notification when policy is approved and issued
func (s *SMTPClient) SendPolicyIssuedEmail(to, fullName, policyNumber, productName string, sumAssured int64) error {
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #2196F3; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; background-color: #f9f9f9; }
		.policy-box { background-color: #fff; border: 2px solid #2196F3; padding: 15px; margin: 20px 0; }
		.footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>🎉 Your Policy Has Been Issued!</h1>
		</div>
		<div class="content">
			<h2>Hi %s,</h2>
			<p>Great news! Your insurance application has been approved and your policy is now active.</p>
			<div class="policy-box">
				<h3>Policy Details:</h3>
				<p><strong>Policy Number:</strong> %s</p>
				<p><strong>Product:</strong> %s</p>
				<p><strong>Sum Assured:</strong> Rp %d</p>
			</div>
			<p>Your policy documents will be available in your dashboard shortly.</p>
			<p>Best regards,<br>The Insurance Platform Team</p>
		</div>
		<div class="footer">
			<p>&copy; 2026 Insurance Platform. All rights reserved.</p>
		</div>
	</div>
</body>
</html>
	`, fullName, policyNumber, productName, sumAssured)

	return s.SendEmail(&EmailMessage{
		To:      []string{to},
		Subject: "Policy Issued - " + policyNumber,
		Body:    body,
		IsHTML:  true,
	})
}

// SendClaimStatusUpdateEmail sends notification when claim status changes
func (s *SMTPClient) SendClaimStatusUpdateEmail(to, fullName, claimNumber, status, notes string) error {
	statusEmoji := "📋"
	statusColor := "#FF9800"
	
	switch status {
	case "approved":
		statusEmoji = "✅"
		statusColor = "#4CAF50"
	case "rejected":
		statusEmoji = "❌"
		statusColor = "#f44336"
	case "under_review":
		statusEmoji = "🔍"
		statusColor = "#2196F3"
	}

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: %s; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; background-color: #f9f9f9; }
		.claim-box { background-color: #fff; border: 2px solid %s; padding: 15px; margin: 20px 0; }
		.footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>%s Claim Status Update</h1>
		</div>
		<div class="content">
			<h2>Hi %s,</h2>
			<p>There's an update on your insurance claim.</p>
			<div class="claim-box">
				<h3>Claim Details:</h3>
				<p><strong>Claim Number:</strong> %s</p>
				<p><strong>Status:</strong> %s</p>
				%s
			</div>
			<p>You can view full details in your dashboard.</p>
			<p>Best regards,<br>The Insurance Platform Team</p>
		</div>
		<div class="footer">
			<p>&copy; 2026 Insurance Platform. All rights reserved.</p>
		</div>
	</div>
</body>
</html>
	`, statusColor, statusColor, statusEmoji, fullName, claimNumber, status, 
	func() string {
		if notes != "" {
			return fmt.Sprintf("<p><strong>Notes:</strong> %s</p>", notes)
		}
		return ""
	}())

	return s.SendEmail(&EmailMessage{
		To:      []string{to},
		Subject: "Claim Update - " + claimNumber,
		Body:    body,
		IsHTML:  true,
	})
}

// SendPasswordResetEmail sends password reset link to user
func (s *SMTPClient) SendPasswordResetEmail(to, fullName, resetToken string) error {
	// Frontend URL should be configured, using localhost for now
	resetURL := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", resetToken)
	
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #FF5722; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; background-color: #f9f9f9; }
		.reset-box { background-color: #fff; border: 2px solid #FF5722; padding: 15px; margin: 20px 0; }
		.button { display: inline-block; padding: 12px 24px; background-color: #FF5722; color: white; text-decoration: none; border-radius: 4px; margin: 10px 0; }
		.footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
		.warning { color: #f44336; font-size: 14px; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>🔐 Password Reset Request</h1>
		</div>
		<div class="content">
			<h2>Hi %s,</h2>
			<p>We received a request to reset your password. Click the button below to create a new password:</p>
			<div class="reset-box">
				<a href="%s" class="button">Reset Password</a>
			</div>
			<p>Or copy this link to your browser:</p>
			<p style="word-break: break-all; color: #666; font-size: 12px;">%s</p>
			<p class="warning"><strong>⚠️ Important:</strong> This link will expire in 1 hour. If you didn't request this, please ignore this email.</p>
			<p>Best regards,<br>The Insurance Platform Team</p>
		</div>
		<div class="footer">
			<p>&copy; 2026 Insurance Platform. All rights reserved.</p>
		</div>
	</div>
</body>
</html>
	`, fullName, resetURL, resetURL)

	return s.SendEmail(&EmailMessage{
		To:      []string{to},
		Subject: "Password Reset Request",
		Body:    body,
		IsHTML:  true,
	})
}

// PreviewEmail returns HTML preview of email template with sample variables
// templateType: "welcome" | "policy" | "claim" | "reset"
func (s *SMTPClient) PreviewEmail(templateType string) (string, string, error) {
	switch templateType {
	case "welcome":
		return "Welcome to Our Insurance Platform", fmt.Sprintf(welcomeTemplate, "John Doe"), nil
	case "policy":
		return "Policy Issued - POL-2026-0001", fmt.Sprintf(policyTemplate, "John Doe", "POL-2026-0001", "Asuransi Jiwa Premium", int64(500000000)), nil
	case "claim":
		status := "approved"
		statusEmoji := "✅"
		statusColor := "#4CAF50"
		return "Claim Update - CLM-2026-001", fmt.Sprintf(claimTemplate, statusColor, statusColor, statusEmoji, "John Doe", "CLM-2026-001", status, "<p><strong>Notes:</strong> Klaim disetujui setelah investigasi</p>"), nil
	case "reset":
		resetURL := "http://localhost:3000/reset-password?token=sample-reset-token-abc123"
		return "Password Reset Request", fmt.Sprintf(resetTemplate, "John Doe", resetURL, resetURL), nil
	default:
		return "", "", fmt.Errorf("unknown template type: %s", templateType)
	}
}

// Template constants for preview
const welcomeTemplate = `
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; background-color: #f9f9f9; }
		.footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>Welcome to Our Insurance Platform!</h1>
		</div>
		<div class="content">
			<h2>Hi %s,</h2>
			<p>Thank you for registering with us! We're excited to have you on board.</p>
			<p>You can now:</p>
			<ul>
				<li>Browse our insurance products</li>
				<li>Submit applications online</li>
				<li>Track your policies</li>
				<li>File and manage claims</li>
			</ul>
			<p>If you have any questions, feel free to reach out to our support team.</p>
			<p>Best regards,<br>The Insurance Platform Team</p>
		</div>
		<div class="footer">
			<p>&copy; 2026 Insurance Platform. All rights reserved.</p>
		</div>
	</div>
</body>
</html>
`

const policyTemplate = `
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #2196F3; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; background-color: #f9f9f9; }
		.policy-box { background-color: #fff; border: 2px solid #2196F3; padding: 15px; margin: 20px 0; }
		.footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>🎉 Your Policy Has Been Issued!</h1>
		</div>
		<div class="content">
			<h2>Hi %s,</h2>
			<p>Great news! Your insurance application has been approved and your policy is now active.</p>
			<div class="policy-box">
				<h3>Policy Details:</h3>
				<p><strong>Policy Number:</strong> %s</p>
				<p><strong>Product:</strong> %s</p>
				<p><strong>Sum Assured:</strong> Rp %d</p>
			</div>
			<p>Your policy documents will be available in your dashboard shortly.</p>
			<p>Best regards,<br>The Insurance Platform Team</p>
		</div>
		<div class="footer">
			<p>&copy; 2026 Insurance Platform. All rights reserved.</p>
		</div>
	</div>
</body>
</html>
`

const claimTemplate = `
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: %s; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; background-color: #f9f9f9; }
		.claim-box { background-color: #fff; border: 2px solid %s; padding: 15px; margin: 20px 0; }
		.footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>%s Claim Status Update</h1>
		</div>
		<div class="content">
			<h2>Hi %s,</h2>
			<p>There's an update on your insurance claim.</p>
			<div class="claim-box">
				<h3>Claim Details:</h3>
				<p><strong>Claim Number:</strong> %s</p>
				<p><strong>Status:</strong> %s</p>
				%s
			</div>
			<p>You can view full details in your dashboard.</p>
			<p>Best regards,<br>The Insurance Platform Team</p>
		</div>
		<div class="footer">
			<p>&copy; 2026 Insurance Platform. All rights reserved.</p>
		</div>
	</div>
</body>
</html>
`

const resetTemplate = `
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #FF5722; color: white; padding: 20px; text-align: center; }
		.content { padding: 20px; background-color: #f9f9f9; }
		.reset-box { background-color: #fff; border: 2px solid #FF5722; padding: 15px; margin: 20px 0; }
		.button { display: inline-block; padding: 12px 24px; background-color: #FF5722; color: white; text-decoration: none; border-radius: 4px; margin: 10px 0; }
		.footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
		.warning { color: #f44336; font-size: 14px; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>🔐 Password Reset Request</h1>
		</div>
		<div class="content">
			<h2>Hi %s,</h2>
			<p>We received a request to reset your password. Click the button below to create a new password:</p>
			<div class="reset-box">
				<a href="%s" class="button">Reset Password</a>
			</div>
			<p>Or copy this link to your browser:</p>
			<p style="word-break: break-all; color: #666; font-size: 12px;">%s</p>
			<p class="warning"><strong>⚠️ Important:</strong> This link will expire in 1 hour. If you didn't request this, please ignore this email.</p>
			<p>Best regards,<br>The Insurance Platform Team</p>
		</div>
		<div class="footer">
			<p>&copy; 2026 Insurance Platform. All rights reserved.</p>
		</div>
	</div>
</body>
</html>
`
