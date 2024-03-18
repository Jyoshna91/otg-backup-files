# scp -r tcs@10.133.35.133:/home/tcs/149ondatra/sample/ondatra/otg/otg3/nping67_output.xlsx Desktop

package otg3
 
import (
"fmt"
"os/exec"
"strings"
"testing"
"time"
 
"github.com/tealeg/xlsx"
"otg3/otg_lib"
)
 
func TestConnectToRouter67(t *testing.T) {
// Load test parameters from library
params := otg_lib.DefaultTestParams()
 
// Create a new Excel file
file := xlsx.NewFile()
sheet, err := file.AddSheet("Nping Output")
if err != nil {
t.Fatalf("Error creating Excel sheet: %v", err)
}
 
// Define styles for "Passed" and "Failed" results
outputStyle := xlsx.NewStyle()
outputStyle.Fill = *xlsx.NewFill("solid", "FFFFFF00", "FFFFFF00") // Yellow for output
outputStyle.ApplyFill = true
 
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
executeAndLogNpingCommand("Sending UDP Traffic using Nping otg ", fmt.Sprintf("sudo nping --udp -c 10 --rate 10 %s %s", params.SourceIP, params.DestinationIP), sheet, outputStyle, passStyle, failStyle, t)
 
// Wait for the duration before the next test
time.Sleep(params.UDPDuration)
 
// Execute and log TCP nping command
executeAndLogNpingCommand("Sending TCP Traffic using Nping otg", fmt.Sprintf("sudo nping --tcp -c 10 --rate 10 %s %s", params.SourceIP, params.DestinationIP), sheet, outputStyle, passStyle, failStyle, t)
 
// Save the Excel file
filePath := "/home/tcs/149ondatra/sample/ondatra/otg/otg3/nping67_output.xlsx"
if err := file.Save(filePath); err != nil {
t.Fatalf("Error saving Excel file: %v", err)
}
 
// Print the paths at the end of the output
t.Log("Nping traffic generation for UDP and TCP completed.\n")
t.Logf("path for testcase : /home/tcs/149ondatra/sample/ondatra/otg/otg3/nping87_test.go")
t.Logf("path for excel:%v",filePath)
}
 

func executeAndLogNpingCommand(testName, command string, sheet *xlsx.Sheet, outputStyle, passStyle, failStyle *xlsx.Style, t *testing.T) {
    output, _ := executeNpingCommand(command, t)
    resultStyle := passStyle
    resultText := "Passed"

    // Check for packet loss
    if strings.Contains(output, "Lost: 0 (0.00%)") {
        resultStyle = passStyle
        resultText = "Passed"
    } else {
        resultStyle = failStyle
        resultText = "Failed"
    }

    // Extract relevant information from the output
    startIndex := strings.Index(output, "Statistics for host")
    endIndex := strings.Index(output, "Nping done")
    if startIndex != -1 && endIndex != -1 {
        output = output[startIndex:endIndex]
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
t.Logf("Nping output: %s", string(output))
return string(output), nil
}
