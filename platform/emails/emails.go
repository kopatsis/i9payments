package emails

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendConfirmation(email, name string) error {
	from := mail.NewEmail("i9 Team", "noreply@i9fit.co")
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
			<h1>Congrats</h1>
			<p>Your i9 Giga Membership has officially begun!</p>
		</div>
	</body>
	</html>
	`

	textContent := "Congrats\nYour i9 Giga Membership has officially begun!"

	message := mail.NewSingleEmail(from, "Confirmation: Your Membership Has Started", to, textContent, htmlContent)
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
	from := mail.NewEmail("i9 Team", "noreply@i9fit.co")
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
			<h1>Aw</h1>
			<p>Your i9 Giga Membership has officially been cancelled. We're sorry to see you go. </p>
			<p>You will be able to use the same features until the end of your billing cycle and will not be charged again. </p>
		</div>
	</body>
	</html>
	`

	textContent := "Aw\nYour i9 Giga Membership has officially been cancelled. We're sorry to see you go.\nYou will be able to use the same features until the end of your billing cycle and will not be charged again."

	message := mail.NewSingleEmail(from, "Confirmation: Your Membership Been Cancelled", to, textContent, htmlContent)
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
	from := mail.NewEmail("i9 Team", "noreply@i9fit.co")
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
			<h1>Nice</h1>
			<p>Your i9 Giga Membership has officially been un-cancelled. We are so back.</p>
		</div>
	</body>
	</html>
	`

	textContent := "Nice\nYour i9 Giga Membership has officially been un-cancelled. We are so back."

	message := mail.NewSingleEmail(from, "Confirmation: Your Membership Been Un-cancelled", to, textContent, htmlContent)
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

func SendOver(email, name string) error {
	from := mail.NewEmail("i9 Team", "noreply@i9fit.co")
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
			<h1>All Done</h1>
			<p>Your i9 Giga Membership has officially ended. Feel free to restart it or reach out to us.</p>
		</div>
	</body>
	</html>
	`

	textContent := "All Done\nYour i9 Giga Membership has officially ended. Feel free to restart it or reach out to us."

	message := mail.NewSingleEmail(from, "Confirmation: Your Membership Ended", to, textContent, htmlContent)
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

func SendFailureNotification(email, name string) error {
	from := mail.NewEmail("i9 Team", "noreply@i9fit.co")
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
			<h1>Uh Oh</h1>
			<p>Your payment has failed for your i9 Giga Membership. In the next 4 days update your information so we can retry, or your benefits will be cancelled.</p>
		</div>
	</body>
	</html>
	`

	textContent := "Uh Oh\nYour payment has failed for your i9 Giga Membership. In the next 4 days update your information so we can retry, or your benefits will be cancelled."

	message := mail.NewSingleEmail(from, "Urgent: i9 Giga Membership Payment Failed", to, textContent, htmlContent)
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
