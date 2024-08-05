package forecast

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/mattismoel/tmpr/internal/model"
)

type Printer interface {
	Print(model.Forecast) error
}

type JSONPrinter struct {
	w io.Writer
}

type StdOutPrinter struct{}

func (p JSONPrinter) Print(f model.Forecast) error {
	err := json.NewEncoder(p.w).Encode(f)
	if err != nil {
		return fmt.Errorf("could not encode JSON to writer: %v", err)
	}

	return nil
}

func (p StdOutPrinter) Print(f model.Forecast) error {
	// At <location> it is currently <temp>, <description>.
	s := fmt.Sprintf("At %s it is currently %.1f°, %s\n", f.Location.String(), f.Weather.Temperature, f.Weather.Description)
	os.Stdout.WriteString(s)
	return nil
}
