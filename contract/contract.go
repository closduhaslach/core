// Package contract
package contract

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/closduhaslach/core/smoobu"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v3"
)

type Contract struct {
	Template *drive.File
	Booking  *smoobu.Booking
	Folder   *drive.File
}

func (c *Contract) Filename() string {
	return fmt.Sprintf("Contract_%s_%s", c.Booking.Firstname, c.Booking.Lastname)
}

func (c *Contract) DriveFile(folder *drive.File) *drive.File {
	return &drive.File{
		Name:    c.Filename(),
		Parents: []string{folder.Id},
	}
}

func (c *Contract) GeneratePDF(sDocs *docs.Service, sDrive *drive.Service) error {
	// ctx := context.Background()
	// sDocs, _ := docs.NewService(ctx, option.WithHTTPClient(client))
	// sDrive, _ := drive.NewService(ctx, option.WithHTTPClient(client))

	target := c.DriveFile(c.Folder)

	check, err := sDrive.Files.List().Q(fmt.Sprintf("name = '%s' and '%s' in parents and trashed = false", c.Filename(), c.Folder.Id)).Do()
	if err != nil {
		return err
	}

	if len(check.Files) > 1 {
		return fmt.Errorf("multiple files found with name %s in folder %s", target.Name, c.Folder.Id)
	}

	if len(check.Files) == 1 {
		// File already exists, delete it first
		err = sDrive.Files.Delete(check.Files[0].Id).Do()
		if err != nil {
			return err
		}
	}

	newFile, _ := sDrive.Files.Copy(c.Template.Id, target).Do()
	newDocID := newFile.Id

	// 2. Build placeholder replacements
	var requests []*docs.Request
	requests, err = c.TemplateData()
	if err != nil {
		return err
	}

	arrivalTime, _ := time.Parse("2006-01-02", c.Booking.Arrival)
	departureTime, _ := time.Parse("2006-01-02", c.Booking.Departure)
	numberOfNights := departureTime.Sub(arrivalTime).Hours() / 24

	requests = append(requests, &docs.Request{
		ReplaceAllText: &docs.ReplaceAllTextRequest{
			ContainsText: &docs.SubstringMatchCriteria{
				Text:      "{{numberOfNights}}",
				MatchCase: true,
			},
			ReplaceText: fmt.Sprintf("%d", int(numberOfNights)),
		},
	})

	raw, err := json.MarshalIndent(requests, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("Requests: %s\n", string(raw))
	os.WriteFile("raw.json", raw, 0o644)

	_, err = sDocs.Documents.BatchUpdate(newDocID, &docs.BatchUpdateDocumentRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		return err
	}

	// 3. Export as PDF
	resp, _ := sDrive.Files.Export(newDocID, "application/pdf").Download()
	pdfBytes, _ := io.ReadAll(resp.Body)
	_ = os.WriteFile(target.Name, pdfBytes, 0o644)

	return nil
}

func (c *Contract) TemplateData() ([]*docs.Request, error) {
	var data map[string]any
	raw, err := json.Marshal(c.Booking)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(raw, &data)
	if err != nil {
		return nil, err
	}
	return dataToRequests(data, "")
}

func dataToRequests(data map[string]any, prefix string) ([]*docs.Request, error) {
	var requests []*docs.Request

	for key, val := range data {
		nv := ""
		switch v := val.(type) {
		case map[string]any:
			nestedRequests, err := dataToRequests(v, prefix+key+".")
			if err != nil {
				return nil, fmt.Errorf("error processing nested data for key %s: %w", key, err)
			}
			requests = append(requests, nestedRequests...)
		case []any:
			for i, item := range v {
				nestedMap, ok := item.(map[string]any)
				if !ok {
					continue
				}

				idx := fmt.Sprintf("%d", i)
				if key == "priceElements" {
					if t, ok := nestedMap["type"].(string); ok {
						idx = t
					}
				}
				nestedRequests, err := dataToRequests(nestedMap, fmt.Sprintf("%s%s[%s].", prefix, key, idx))
				if err != nil {
					return nil, fmt.Errorf("error processing nested array data for key %s[%s]: %w", key, idx, err)
				}
				requests = append(requests, nestedRequests...)
			}
		case string:
			nv = v
		case float64:
			if v == float64(int(v)) {
				nv = fmt.Sprintf("%d", int(v))
			} else {
				nv = fmt.Sprintf("%.2f", v)
			}
		case int:
			nv = fmt.Sprintf("%d", v)
		case bool:
			nv = fmt.Sprintf("%t", v)
		}
		if nv != "" {
			requests = append(requests, &docs.Request{
				ReplaceAllText: &docs.ReplaceAllTextRequest{
					ContainsText: &docs.SubstringMatchCriteria{
						Text:      "{{" + prefix + key + "}}",
						MatchCase: true,
					},
					ReplaceText: nv,
				},
			})
		}
	}
	return requests, nil
}
