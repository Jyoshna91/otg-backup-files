package otg

import (
    "fmt"
    "log"
    "os/exec"
    "testing"
    "time"
    "net/smtp"
)

func TestConnectToRouter11(t *testing.T) {
    const (
        sourceIP      = "10.133.35.158"
        destinationIP = "10.133.35.143"
        duration      = 5 * time.Second
    )

    // Construct nping command
    command := fmt.Sprintf("sudo nping --tcp -c 10 --rate 10 %s %s", sourceIP, destinationIP)

    // Execute nping command
    cmd := exec.Command("sh", "-c", command)

    output, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("Error starting nping command: %v. Output: %s", err, output)
    }

    t.Logf("nping output: %s", string(output))

    // Wait for the specified duration
    time.Sleep(duration)

    // Check if the process has already finished
    if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
        t.Log("nping process has already finished.")
        sendEmail("Test Result", "Nping test case completed successfully.")
        return
    }

    // Stop nping command
    err = cmd.Process.Kill()
    if err != nil {
        t.Fatalf("Error stopping nping command: %v", err)
    }

    t.Log("nping traffic generation completed.")

    sendEmail("Test Result", "Nping test case completed successfully.")
}

func sendEmail(subject, body string) {
    from := "keerthi.narne@tcs.com"
    password := "Keeru@16"
    to := []string{"vijaya.punganuri@tcs.com"}

    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    msg := "From: " + from + "\n" +
        "To: " + to[0] + "\n" +
        "Subject: " + subject + "\n\n" +
        body

    auth := smtp.PlainAuth("", from, password, smtpHost)

    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(msg))
    if err != nil {
        log.Fatalf("Failed to send email: %v", err)
    }
    log.Println("Email sent successfully")
}
