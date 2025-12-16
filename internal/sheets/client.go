package sheets

import (
	"context"
	"fmt"
	"os"

	"voiceline-audio-backend/internal/common"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Client struct {
	service *sheets.Service
	sheetID string
}

func NewClient(credentialsFile, sheetID string) (*Client, error) {
	ctx := context.Background()

	credentials, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, common.NewInternalServerError(
			fmt.Sprintf("Unable to read credentials file: %s", credentialsFile),
			err,
		)
	}

	config, err := google.JWTConfigFromJSON(credentials, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, common.NewInternalServerError("Unable to parse credentials", err)
	}

	service, err := sheets.NewService(ctx, option.WithHTTPClient(config.Client(ctx)))
	if err != nil {
		return nil, common.NewInternalServerError("Unable to create Sheets service", err)
	}

	return &Client{
		service: service,
		sheetID: sheetID,
	}, nil
}

func (c *Client) AppendRow(ctx context.Context, values []any) error {
	valueRange := &sheets.ValueRange{
		Values: [][]any{values},
	}

	_, err := c.service.Spreadsheets.Values.Append(
		c.sheetID,
		"Sheet1!A:F",
		valueRange,
	).ValueInputOption("RAW").Context(ctx).Do()

	if err != nil {
		return common.NewInternalServerError("sheets append failed", err)
	}

	return nil
}

func (c *Client) CreateHeaderRow(ctx context.Context) error {
	resp, err := c.service.Spreadsheets.Values.Get(c.sheetID, "Sheet1!A1:F1").Context(ctx).Do()
	if err != nil {
		return common.NewInternalServerError("couldn't check headers", err)
	}

	if len(resp.Values) == 0 {
		headers := []any{
			"Timestamp",
			"Filename",
			"Transcription",
			"Summary",
			"Action Items",
			"Status",
		}

		valueRange := &sheets.ValueRange{
			Values: [][]any{headers},
		}

		_, err := c.service.Spreadsheets.Values.Update(
			c.sheetID,
			"Sheet1!A1:F1",
			valueRange,
		).ValueInputOption("RAW").Context(ctx).Do()

		if err != nil {
			return common.NewInternalServerError("header creation failed", err)
		}
	}

	return nil
}
