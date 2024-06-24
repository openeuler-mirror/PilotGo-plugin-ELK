package global

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ReadJSONResponse(r *http.Response, v any) error {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("response status: %d, body: %s", r.StatusCode, string(b))
	}

	err = json.Unmarshal(b, v)
	if err != nil {
		return fmt.Errorf("unmarshalling response json: %w", err)
	}
	return nil
}

func ReadFileJSON(path string, v any) error {
	bytes, err := FileReadBytes(path)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	err = json.Unmarshal(bytes, v)
	if err != nil {
		return fmt.Errorf("unmarshalling file json: %w", err)
	}
	return nil
}
