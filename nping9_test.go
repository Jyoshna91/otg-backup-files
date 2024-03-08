 
import (
    "fmt"
    "os/exec"
    "strings"
    "testing"
    "time"
 
"github.com/tealeg/xlsx"
)
 
func TestConnectToRouter18(t *testing.T) {
    const (
sourceIP = "10.133.35.158"
destinationIP = "10.133.35.143"
        udpDuration   = 10 * time.Second
        npingCommand  = "sudo nping -c 10 --rate 10 %s %s"
    )
 
    // Create a new Excel file
    file := xlsx.NewFile()
    sheet, err := file.AddSheet("Nping Output")
    if err != nil {
        t.Fatalf("Error creating Excel sheet: %v", err)
    }
 
    // Define a style for output cells
    outputStyle := xlsx.NewStyle()
    outputStyle.Fill = *xlsx.NewFill("solid", "FFFFFF00", "FFFFFF00") // Yellow for output
    outputStyle.ApplyFill = true
 
    // Define styles for "Passed" and "Failed" results
    passStyle := xlsx.NewStyle()
    passStyle.Fill = *xlsx.NewFill("solid", "FF00FF00", "FF00FF00") // Green for "Passed"
    passStyle.ApplyFill = true
 
    failStyle := xlsx.NewStyle()
    failStyle.Fill = *xlsx.NewFill("solid", "FFFF0000", "FFFF0000") // Red for "Failed"
    failStyle.ApplyFill = true
 
    // Add headers
    headerRow := sheet.AddRow()
    headerRow.AddCell().Value = "Test Type"
    headerRow.AddCell().Value = "Output"
    headerRow.AddCell().Value = "Result"
 
    // Execute and log UDP nping command
    executeAndLogNpingCommand("UDP Test", fmt.Sprintf(npingCommand, "--udp", sourceIP, destinationIP), sheet, outputStyle, passStyle, failStyle, t)
 
    // Wait for the duration before the next test
    time.Sleep(udpDuration)
 
    // Execute and log TCP nping command
    executeAndLogNpingCommand("TCP Test", fmt.Sprintf(npingCommand, "--tcp", sourceIP, destinationIP), sheet, outputStyle, passStyle, failStyle, t)
 
    // Save the Excel file
    filePath := "/home/tcs/sample/ondatra/otg/nping55_output.xlsx"
    if err := file.Save(filePath); err != nil {
        t.Fatalf("Error saving Excel file: %v", err)
    }
 
    t.Log("Nping traffic generation for UDP and TCP completed.")
}
 
func executeAndLogNpingCommand(testName, command string, sheet *xlsx.Sheet, outputStyle, passStyle, failStyle *xlsx.Style, t *testing.T) {
    output, err := executeNpingCommand(command, t)
    resultStyle := passStyle
    resultText := "Passed"
    if err != nil {
        resultStyle = failStyle
        resultText = "Failed"
    }
 
    // Record the test name, output, and result in the sheet
    row := sheet.AddRow()
    testCell := row.AddCell() // Test name in column A
    testCell.Value = testName
 
    outputCell := row.AddCell() // Output in column B
    outputCell.Value = output
    outputCell.SetStyle(outputStyle)
 
    resultCell := row.AddCell() // Result in column C, styled based on the outcome
    resultCell.Value = resultText
    resultCell.SetStyle(resultStyle)
}
 
func executeNpingCommand(command string, t *testing.T) (string, error) {
    cmd := exec.Command("sh", "-c", command)
    output, err := cmd.CombinedOutput()
    if err != nil {
        t.Logf("Nping command error: %v. Output: %s", err, string(output))
        return string(output), err
    }
    lines := strings.Split(string(output), "\n")
    for _, line := range lines {
        if strings.Contains(line, "Probes Sent") || strings.Contains(line, "Rcvd") {
            t.Log(line)
        }
    }
    return string(output), nil
} 
func extractValue(output, startTag, endTag string) string {
    startIndex := strings.Index(output, startTag)
    if startIndex == -1 {
        return "N/A"
    }
    startIndex += len(startTag)
    endIndex := strings.Index(output[startIndex:], endTag)
    if endIndex == -1 {
        return "N/A"
    }
    return strings.TrimSpace(output[startIndex : startIndex+endIndex])
}
