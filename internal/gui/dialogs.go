package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ShowMigrationDialog shows a dialog asking the user if they want to encrypt their database
func ShowMigrationDialog(win fyne.Window, onMigrate func(), onSkip func()) {
	title := "Database Encryption Available"
	message := "Your clipboard database is currently unencrypted.\n\n" +
		"Would you like to encrypt it now for better security?\n\n" +
		"• Your data will be encrypted with AES-256\n" +
		"• Encryption key stored securely in your system keychain\n" +
		"• A backup will be created automatically\n\n" +
		"You can continue without encryption if you prefer."

	dialog.ShowConfirm(title, message, func(encrypt bool) {
		if encrypt {
			onMigrate()
		} else {
			onSkip()
		}
	}, win)
}

// ShowMigrationProgressDialog shows a progress dialog during migration
func ShowMigrationProgressDialog(win fyne.Window) dialog.Dialog {
	progress := widget.NewProgressBarInfinite()
	content := container.NewVBox(
		widget.NewLabel("Encrypting your clipboard database..."),
		widget.NewLabel("This may take a moment. Please wait."),
		progress,
	)

	d := dialog.NewCustomWithoutButtons("Migration in Progress", content, win)
	d.Show()
	return d
}

// ShowMigrationSuccessDialog shows a success message after migration
func ShowMigrationSuccessDialog(win fyne.Window, onOK func()) {
	message := "Your clipboard database has been successfully encrypted!\n\n" +
		"• All data migrated safely\n" +
		"• Backup created\n" +
		"• Encryption key stored in system keychain\n\n" +
		"The app will now restart to use the encrypted database."

	d := dialog.NewInformation("Encryption Complete", message, win)
	if onOK != nil {
		d.SetOnClosed(onOK)
	}
	d.Show()
}

// ShowMigrationErrorDialog shows an error message if migration fails
func ShowMigrationErrorDialog(win fyne.Window, err error) {
	message := fmt.Sprintf("Failed to encrypt the database:\n\n%v\n\n"+
		"Your original database is safe and unchanged.\n"+
		"The app will continue using the unencrypted database.", err)

	dialog.ShowError(fmt.Errorf(message), win)
}
