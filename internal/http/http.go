package http_int

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func PrintResponse(resp *http.Response) {
	if resp == nil {
		fmt.Println("Response is nil")
		return
	}

	fmt.Println("═══════════════════════════════════════")
	fmt.Println("  HTTP RESPONSE")
	fmt.Println("═══════════════════════════════════════")

	// Status
	fmt.Printf("Status:        %s\n", resp.Status)
	fmt.Printf("Status Code:   %d\n", resp.StatusCode)
	fmt.Printf("Proto:         %s\n", resp.Proto)
	fmt.Printf("ProtoMajor:    %d\n", resp.ProtoMajor)
	fmt.Printf("ProtoMinor:    %d\n", resp.ProtoMinor)
	fmt.Printf("ContentLength: %d\n", resp.ContentLength)
	fmt.Printf("Uncompressed:  %t\n", resp.Uncompressed)

	// Transfer-Encoding
	if len(resp.TransferEncoding) > 0 {
		fmt.Printf("TransferEnc:   %s\n", strings.Join(resp.TransferEncoding, ", "))
	}

	// Headers
	fmt.Println("\n── Headers ────────────────────────────")
	for key, values := range resp.Header {
		fmt.Printf("  %-30s %s\n", key+":", strings.Join(values, ", "))

	}

	// Cookies
	if cookies := resp.Cookies(); len(cookies) > 0 {
		fmt.Println("\n── Cookies ────────────────────────────")
		for _, c := range cookies {
			fmt.Printf("  %-20s = %s\n", c.Name, c.Value)
			if c.Domain != "" {
				fmt.Printf("    Domain:   %s\n", c.Domain)
			}
			if c.Path != "" {
				fmt.Printf("    Path:     %s\n", c.Path)
			}
			if !c.Expires.IsZero() {
				fmt.Printf("    Expires:  %s\n", c.Expires)
			}
			fmt.Printf("    Secure:   %t\n", c.Secure)
			fmt.Printf("    HttpOnly: %t\n", c.HttpOnly)
		}
	}

	// Trailer headers
	if len(resp.Trailer) > 0 {
		fmt.Println("\n── Trailer Headers ────────────────────")
		for key, values := range resp.Trailer {
			fmt.Printf("  %-30s %s\n", key+":", strings.Join(values, ", "))
		}
	}

	// // TLS info
	// if resp.TLS != nil {
	// 	fmt.Println("\n── TLS ────────────────────────────────")
	// 	fmt.Printf("  Version:     0x%04x\n", resp.TLS.Version)
	// 	fmt.Printf("  CipherSuite: 0x%04x\n", resp.TLS.CipherSuite)
	// 	fmt.Printf("  ServerName:  %s\n", resp.TLS.ServerName)
	// }

	// Request info (the request that generated this response)
	// if resp.Request != nil {
	// 	fmt.Println("\n── Original Request ───────────────────")
	// 	fmt.Printf("  Method: %s\n", resp.Request.Method)
	// 	fmt.Printf("  URL:    %s\n", resp.Request.URL)
	// }

	// Body
	if resp.Body != nil {
		fmt.Println("\n── Body ───────────────────────────────")
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("  Error reading body: %v\n", err)
		} else {
			const maxDisplay = 1024
			text := string(body)
			if len(text) > maxDisplay {
				fmt.Printf("%s\n  ... (%d bytes truncated)\n", text[:maxDisplay], len(text)-maxDisplay)
			} else {
				fmt.Println(text)
			}
			resp.Body = io.NopCloser(bytes.NewReader(body))
		}

	}

	fmt.Println("═══════════════════════════════════════")
}
