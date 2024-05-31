package emails

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendConfirmation(email, name string) error {
	from := mail.NewEmail("i9 Team", "main@i9fit.co")
	to := mail.NewEmail(name, email)

	htmlContent := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<style>
			body { font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px; }
			.container { max-width: 600px; margin: auto; background: #ffffff; padding: 20px; border-radius: 8px; }
			h1 { color: #333333; }
			p { color: #666666; }
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Hello, World!</h1>
			<p>This is a sample email with CSS styling.</p>
		</div>
	</body>
	</html>
	`

	message := mail.NewSingleEmail(from, "Confirmation: Your Membership Has Started", to, "", htmlContent)
	client := sendgrid.NewSendClient("SENDGRID_KEY")
	response, err := client.Send(message)
	if err != nil {
		return err
	}
	if response.StatusCode >= 400 {
		return fmt.Errorf("failed to send email: %v", response.StatusCode)
	}
	return nil
}

func SendCancelled(email, name string) error {
	from := mail.NewEmail("i9 Team", "main@i9fit.co")
	to := mail.NewEmail(name, email)

	htmlContent := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<style>
			body { font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px; }
			.container { max-width: 600px; margin: auto; background: #ffffff; padding: 20px; border-radius: 8px; }
			h1 { color: #333333; }
			p { color: #666666; }
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Hello, World!</h1>
			<p>This is a sample email with CSS styling.</p>
		</div>
	</body>
	</html>
	`

	message := mail.NewSingleEmail(from, "Confirmation: Your Membership Has Started", to, "", htmlContent)
	client := sendgrid.NewSendClient("SENDGRID_KEY")
	response, err := client.Send(message)
	if err != nil {
		return err
	}
	if response.StatusCode >= 400 {
		return fmt.Errorf("failed to send email: %v", response.StatusCode)
	}
	return nil
}

func SendUnCancelled(email, name string) error {
	from := mail.NewEmail("i9 Team", "main@i9fit.co")
	to := mail.NewEmail(name, email)

	htmlContent := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<style>
			body { font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px; }
			.container { max-width: 600px; margin: auto; background: #ffffff; padding: 20px; border-radius: 8px; }
			h1 { color: #333333; }
			p { color: #666666; }
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Hello, World!</h1>
			<p>This is a sample email with CSS styling.</p>
		</div>
	</body>
	</html>
	`

	message := mail.NewSingleEmail(from, "Confirmation: Your Membership Has Started", to, "", htmlContent)
	client := sendgrid.NewSendClient("SENDGRID_KEY")
	response, err := client.Send(message)
	if err != nil {
		return err
	}
	if response.StatusCode >= 400 {
		return fmt.Errorf("failed to send email: %v", response.StatusCode)
	}
	return nil
}
