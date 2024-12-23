package ifsc

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
)

var baseURL = "https://ifsc.razorpay.com/"

// IFSC represents the IFSC service
type IFSC struct {
    client *http.Client
}

// New returns a new instance of IFSC service
func New() *IFSC {
    return &IFSC{
        client: &http.Client{
            Timeout: 5 * time.Second,
        },
    }
}

// Query handles the IFSC code lookup
func (i *IFSC) Query(q string) ([]string, error) {
    // Clean up the query string and remove the trailing dot if present
    ifscCode := strings.TrimSuffix(q, ".")
    ifscCode = strings.TrimSuffix(ifscCode, ".ifsc")
    ifscCode = strings.ToUpper(ifscCode)

    if len(ifscCode) < 11 { // IFSC codes are 11 characters
        return nil, fmt.Errorf("invalid IFSC code length: %d", len(ifscCode))
    }

    url := baseURL + ifscCode

    // Make HTTP request to Razorpay API
    resp, err := i.client.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch IFSC details: %w", err)
    }
    defer resp.Body.Close()

    // Read response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }

    // Handle non-200 responses
    if resp.StatusCode != http.StatusOK {
        return []string{fmt.Sprintf(`%s.ifsc. 1 IN TXT "IFSC code %s not found"`, ifscCode, ifscCode)}, nil
    }

    // Parse JSON response
    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }

    // Format response for DNS
    output := []string{
        fmt.Sprintf(`%s.ifsc. 1 IN TXT "BANK: %s"`, ifscCode, result["BANK"]),
        fmt.Sprintf(`%s.ifsc. 1 IN TXT "BRANCH: %s"`, ifscCode, result["BRANCH"]),
        //fmt.Sprintf(`%s.ifsc. 1 IN TXT "ADDRESS: %s"`, ifscCode, result["ADDRESS"]),
        fmt.Sprintf(`%s.ifsc. 1 IN TXT "CITY: %s"`, ifscCode, result["CITY"]),
        fmt.Sprintf(`%s.ifsc. 1 IN TXT "STATE: %s"`, ifscCode, result["STATE"]),
    }

    return output, nil
}

// Dump is not implemented for this service
func (i *IFSC) Dump() ([]byte, error) {
    return nil, nil
}