package command

import (
    "encoding/csv"
    "fmt"
    "os"

    "github.com/xuri/excelize/v2"
)

// writeXLSXFromCSV creates an .xlsx file with the same contents as the provided CSV.
func writeXLSXFromCSV(csvPath, xlsxPath string) error {
    f, err := os.Open(csvPath)
    if err != nil {
        return err
    }
    defer f.Close()

    r := csv.NewReader(f)
    rows, err := r.ReadAll()
    if err != nil {
        return err
    }

    xf := excelize.NewFile()
    sheet := xf.GetSheetName(0)
    for i, row := range rows {
        axis := fmt.Sprintf("A%d", i+1)
        if err := xf.SetSheetRow(sheet, axis, &row); err != nil {
            return err
        }
    }
    return xf.SaveAs(xlsxPath)
}
