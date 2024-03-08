package otg
 
import (
   // "fmt"
    "os/exec"
    "testing"
    "time"
)
 
func TestConnectToRouter11(t *testing.T) {
    const (
        duration = 5 * time.Second
    )
 
    // Construct nping command
command := "sudo nping --udp -c 10 --rate 10 10.133.35.158 10.133.35.143"
 
    // Execute nping command
    output, err := executeNpingCommand(command)
    if err != nil {
        t.Fatalf("Error starting nping command: %v. Output: %s", err, output)
    }
 
    t.Logf("nping output: %s", output)
 
    // Wait for the specified duration
    time.Sleep(duration)
 
    t.Log("nping traffic generation completed.")
}
 
func executeNpingCommand(command string) (string, error) {
    cmd := exec.Command("sh", "-c", command)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return string(output), err
    }
    return string(output), nil
}
