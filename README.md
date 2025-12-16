# Audio Processing Backend

## Overview

Go backend service that receives audio files, transcribes them using OpenAI Whisper, extracts summaries and action items with GPT-4o, and logs results to Google Sheets.

## Features

- Audio upload with validation (mp3, wav, m4a, etc.)
- Transcription via OpenAI Whisper
- Summary and action item extraction using GPT-4o
- Google Sheets integration (optional)
- CORS enabled

## Prerequisites

- Go 1.23+
- OpenAI API key
- Google Cloud service account (optional, for Sheets)

## Setup

### 1. Install Dependencies

```bash
go mod download
```

### 2. Environment Config

Create `.env` file:

```env
PORT=8080
GIN_MODE=debug
OPENAI_API_KEY=your-key-here
MAX_AUDIO_SIZE_MB=10

GOOGLE_CREDENTIALS_FILE=./credentials.json
GOOGLE_SHEET_ID=your-sheet-id
```

### 3. Google Sheets Setup (Optional)

1. Enable Sheets API in Google Cloud
2. Create service account, download credentials.json
3. Share your sheet with the service account email
4. Add sheet ID to .env

### 4. Run

```bash
go run main.go
```

Server runs on `http://localhost:8080`

## API

### Health Check

```bash
GET /health
```

### Upload Audio

```bash
POST /api/audio/upload
```

Form field: `audio`
Formats: mp3, wav, m4a, webm, ogg, flac, aac, mp4, mpeg, mpga
Max size: 10MB

Example:

```bash
curl -X POST http://localhost:8080/api/audio/upload -F "audio=@recording.m4a"
```

Response:

```json
{
  "status": "success",
  "message": "processed",
  "data": {
    "transcription": "...",
    "summary": "...",
    "action_items": ["..."],
    "timestamp": "2024-01-15T10:30:00Z",
    "filename": "recording.m4a"
  }
}
```

## Architecture

Built with Gin framework. Uses OpenAI Whisper for transcription and GPT-4o for extraction. Google Sheets for simple data storage without needing a database.

Code is organized into packages:

- `audio` - handlers, validation, processing
- `llm` - OpenAI client
- `sheets` - Google Sheets integration
- `common` - config and errors
- `middleware` - CORS and logging

## What I had improve with more time

Kept this under 4 hours, but here is what is missing:

- I skipped tests to prioritize features, but normally I had add unit tests for the parsing logic and integration tests for the API. Also, the Google Sheets sync currently fails silently. I had strictly add retry logic so data is not lost if the Google API hiccups.

- Currently, the app loads files into memory and processes them synchronously. For a real-world scenario, I had implement file streaming to handle large audio and move the processing to a background queue (Redis) to avoid blocking the main thread.

- The endpoint is currently open for the sake of the demo. I had implement Auth middleware, sanitize inputs before sending them to the LLM and validate actual file content (magic bytes) rather than just trusting the extension.

- Sheets is fun for a demo, but we had obviously swap that for a real DB (Postgres) and offload audio storage to S3 with proper structured logging and monitoring.
