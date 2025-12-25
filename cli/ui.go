package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Tony-ArtZ/goxfer/server"
	"github.com/Tony-ArtZ/goxfer/utils"
	"github.com/pterm/pterm"
)

func StartCLI(addr string, port string, tm *server.TransferManager) {
	print("\033[H\033[2J")

	introSpinner, _ := pterm.DefaultSpinner.
		WithStyle(pterm.NewStyle(pterm.FgLightGreen)).
		WithSequence("⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏").
		Start("Initializing Goxfer...")
	time.Sleep(time.Millisecond * 800)
	introSpinner.Stop()

	print("\033[H\033[2J")

	pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("Go", pterm.NewStyle(pterm.FgLightWhite)),
		pterm.NewLettersFromStringWithStyle("xfer", pterm.NewStyle(pterm.FgLightGreen)),
	).Render()

	pterm.Println()

	greenBox := pterm.DefaultBox.
		WithBoxStyle(pterm.NewStyle(pterm.FgLightGreen)).
		WithTitle("🌐 CONNECTION INFO").
		WithTitleTopCenter().
		WithRightPadding(3).
		WithLeftPadding(3).
		WithTopPadding(1).
		WithBottomPadding(1)

	infoText := fmt.Sprintf("%s\n\n%s\n%s\n%s",
		pterm.LightWhite("Scan the QR code below on your device:"),
		pterm.White("• Access URL:  ")+pterm.LightGreen(addr),
		pterm.White("• Server Port: ")+pterm.LightGreen(port),
		pterm.White("• Save Dir:    ")+pterm.LightGreen(tm.SaveDir))

	greenBox.Println(infoText)
	pterm.Println()

	utils.PrintQR(addr)
	pterm.Println()

	pterm.Info.WithPrefix(pterm.Prefix{
		Text:  "READY",
		Style: pterm.NewStyle(pterm.BgLightGreen, pterm.FgBlack),
	}).Println("Server is running. Waiting for connections...")
	pterm.Println()

	scanner := bufio.NewScanner(os.Stdin)
	transferCount := 0

	for {
		pterm.Print(pterm.LightGreen("┃ "))
		fmt.Print(pterm.White("📁 Enter file path (or 'q' to quit): "))

		if scanner.Scan() {
			input := strings.TrimSpace(scanner.Text())

			if input == "q" || input == "quit" || input == "exit" {
				pterm.Println()

				if transferCount > 0 {
					pterm.DefaultBox.WithBoxStyle(pterm.NewStyle(pterm.FgLightGreen)).
						WithHorizontalString("─").
						Printfln("📊 Session Stats: %d file(s) queued", transferCount)
					pterm.Println()
				}

				farewellSpinner, _ := pterm.DefaultSpinner.
					WithStyle(pterm.NewStyle(pterm.FgLightGreen)).
					Start("Shutting down...")
				time.Sleep(time.Millisecond * 500)
				farewellSpinner.Stop()

				pterm.DefaultBox.
					WithBoxStyle(pterm.NewStyle(pterm.FgLightGreen)).
					WithHorizontalString("═").
					WithTopPadding(0).
					WithBottomPadding(0).
					Println(pterm.LightWhite("✨ Thank you for using Goxfer! Goodbye! ✨"))
				break
			}

			// Handle empty input
			if input == "" {
				continue
			}

			// Handle help command
			if input == "help" || input == "h" || input == "?" {
				pterm.Println()
				helpBox := pterm.DefaultBox.
					WithBoxStyle(pterm.NewStyle(pterm.FgLightGreen)).
					WithTitle("💡 HELP").
					WithTitleTopLeft()
				helpBox.Println(
					pterm.White("Commands:\n") +
						pterm.LightGreen("  • Enter file path") + pterm.White(" - Queue file for transfer\n") +
						pterm.LightGreen("  • q/quit/exit") + pterm.White("    - Exit application\n") +
						pterm.LightGreen("  • help/?/h") + pterm.White("        - Show this help"))
				pterm.Println()
				continue
			}

			// Validate file exists
			fileInfo, err := os.Stat(input)
			if os.IsNotExist(err) {
				pterm.Error.WithPrefix(pterm.Prefix{
					Text:  "NOT FOUND",
					Style: pterm.NewStyle(pterm.BgRed, pterm.FgBlack),
				}).Printfln("File or directory does not exist: %s", input)
				continue
			}

			// Show file info
			fileName := filepath.Base(input)
			fileType := "File"
			if fileInfo.IsDir() {
				fileType = "Directory"
			}

			pterm.Println()
			pterm.Info.WithPrefix(pterm.Prefix{
				Text:  "INFO",
				Style: pterm.NewStyle(pterm.BgLightCyan, pterm.FgBlack),
			}).Printfln("%s: %s (Size: %s)", fileType, fileName, formatSize(fileInfo.Size()))

			// Queue with spinner
			spinner, _ := pterm.DefaultSpinner.
				WithStyle(pterm.NewStyle(pterm.FgLightGreen)).
				WithSequence("⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏").
				Start("Queueing file for transfer...")

			time.Sleep(time.Millisecond * 600)
			tm.QueueFile(input)
			transferCount++

			spinner.Success(pterm.Sprintf("%s %s",
				pterm.LightGreen("✅ File queued successfully!"),
				pterm.White("Check your device.")))
			pterm.Println()
		}
	}
}

// Helper function to format file size
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
