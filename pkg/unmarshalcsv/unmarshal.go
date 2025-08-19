package unmarshalcsv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type UnmarshalCsv struct {
	reader      io.Reader
	headerStart int
}

type UnmarshalledData struct {
	// Original direction-specific fields preserved for CSV/XLSX compatibility
	Direction            string `csv:"direction" ommitempty:"true"`
	SourceSpecifier      string `csv:"source_specifier" ommitempty:"true"`
	DestinationNamespace string `csv:"destination_namespace" ommitempty:"true"`
	DestinationSelector  string `csv:"destination_selector" ommitempty:"true"`
	DestinationProtocol  string `csv:"destination_protocol" ommitempty:"true"`
	DestinationPorts     string `csv:"destination_ports" ommitempty:"true"`
	SourceNamespace      string `csv:"source_namespace" ommitempty:"true"`
	SourceSelector       string `csv:"source_selector" ommitempty:"true"`
	NodeRole             string `csv:"node_role" ommitempty:"true"`
	DestinationSpecifier string `csv:"destination_specifier" ommitempty:"true"`
	Comment              string `csv:"comment" ommitempty:"true"`
	NetworkPolicyName    string `csv:"network_policy_name" ommitempty:"true"`

	// Generic aliases (not bound to CSV headers) populated via Normalize()
	// These allow downstream code to be direction-agnostic.
	PolicyName       string `csv:"-"` // alias for NetworkPolicyName
	SubjectNamespace string `csv:"-"` // namespace of the selected pods (source for egress, destination for ingress)
	SubjectSelector  string `csv:"-"` // selector of the selected pods (source for egress, destination for ingress)
	PeerSpecifier    string `csv:"-"` // cidr(s) of the opposite side (destination for egress, source for ingress)
	Protocols        string `csv:"-"` // alias for DestinationProtocol
	Ports            string `csv:"-"` // alias for DestinationPorts
	Role             string `csv:"-"` // alias for NodeRole
}

// NewUnmarshalCsv keeps backward compatibility for CSV files
func NewUnmarshalCsv(fileName string, headerStart int) (*UnmarshalCsv, error) {
	reader, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New("failed to open file")
	}
	return &UnmarshalCsv{reader: reader, headerStart: headerStart}, nil
}

// Unmarshal provides a generic entry point to unmarshal CSV or XLSX by file extension
func Unmarshal(out interface{}, fileName string, headerStart int) error {
	ext := filepath.Ext(fileName)
	switch ext {
	case ".csv":
		u, err := NewUnmarshalCsv(fileName, headerStart)
		if err != nil {
			return err
		}
		return u.UnmarshalCsv(out)
	case ".xlsx":
		return unmarshalXlsx(out, fileName, headerStart)
	default:
		return fmt.Errorf("unsupported file extension: %s", ext)
	}
}

func (u *UnmarshalCsv) UnmarshalCsv(out interface{}) error {
	r := csv.NewReader(u.reader)
	records, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("unmarshalcsv, failed to read csv data: %w", err)
	}
	return unmarshalRecords(out, records, u.headerStart)
}

// unmarshalXlsx reads the first sheet of an .xlsx file and maps rows to the struct slice
func unmarshalXlsx(out interface{}, fileName string, headerStart int) error {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return fmt.Errorf("unmarshalxlsx, failed to open xlsx: %w", err)
	}
	defer func() { _ = f.Close() }()
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return fmt.Errorf("unmarshalxlsx, no sheets found")
	}
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return fmt.Errorf("unmarshalxlsx, failed to read rows: %w", err)
	}
	if len(rows) == 0 {
		return fmt.Errorf("unmarshalxlsx, file has no data")
	}
	return unmarshalRecords(out, rows, headerStart)
}

// unmarshalRecords maps a matrix of strings (records) to the provided slice of structs using `csv` tags
func unmarshalRecords(out interface{}, records [][]string, headerStart int) error {
	outValue := reflect.ValueOf(out)
	if outValue.Kind() != reflect.Ptr || outValue.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("unmarshalcsv: out must be a pointer to a slice")
	}
	if len(records) < headerStart+1 {
		return fmt.Errorf("unmarshalcsv, not enough rows to contain header at index %d", headerStart)
	}
	header := records[headerStart]
	sliceElementType := outValue.Elem().Type().Elem()
	if sliceElementType.Kind() != reflect.Struct {
		return fmt.Errorf("unmarshalcsv, expected a struct, got %s", sliceElementType.Kind())
	}
	headerMap := make(map[int]int)
	for i, colName := range header {
		for j := 0; j < sliceElementType.NumField(); j++ {
			field := sliceElementType.Field(j)
			if tag := field.Tag.Get("csv"); tag == colName {
				headerMap[i] = j
				break
			}
		}
	}
	dataRows := [][]string{}
	if len(records) > headerStart+1 {
		dataRows = records[headerStart+1:]
	}
	slice := reflect.MakeSlice(outValue.Elem().Type(), len(dataRows), len(dataRows))
	for i, row := range dataRows {
		structInstance := slice.Index(i)
		for csvIndex, csvValue := range row {
			if structFieldIndex, ok := headerMap[csvIndex]; ok {
				field := structInstance.Field(structFieldIndex)
				if err := setField(field, csvValue); err != nil {
					return fmt.Errorf("unmarshalcsv, failed to set field on row %d, column %d: %w", i+1, csvIndex, err)
				}
			}
		}
	}
	outValue.Elem().Set(slice)
	return nil
}

// Normalize populates generic alias fields according to the Direction and original fields
func (ud *UnmarshalledData) Normalize() {
	ud.PolicyName = ud.NetworkPolicyName
	ud.Protocols = ud.DestinationProtocol
	ud.Ports = ud.DestinationPorts
	ud.Role = ud.NodeRole
	if strings.EqualFold(ud.Direction, "egress") {
		ud.SubjectNamespace = ud.SourceNamespace
		ud.SubjectSelector = ud.SourceSelector
		ud.PeerSpecifier = ud.DestinationSpecifier
	} else if strings.EqualFold(ud.Direction, "ingress") {
		ud.SubjectNamespace = ud.DestinationNamespace
		ud.SubjectSelector = ud.DestinationSelector
		ud.PeerSpecifier = ud.SourceSpecifier
	}
}

// NormalizeAll populates generic alias fields for a slice of UnmarshalledData
func NormalizeAll(rows []UnmarshalledData) {
	for i := range rows {
		rows[i].Normalize()
	}
}

func setField(field reflect.Value, value string) error {
	if !field.CanSet() {
		return fmt.Errorf("unmarshalcsv: cannot set field %s", field.Kind())
	}
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == "" {
			field.SetInt(0)
			return nil
		}
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("unmarshalcsv, failed to parse %s as integer: %w", field.Type(), err)
		}
		field.SetInt(intValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value == "" {
			field.SetUint(0)
			return nil
		}
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("unmarshalcsv, failed to parse %s as integer: %w", field.Type(), err)
		}
		field.SetUint(uintValue)
	case reflect.Float32, reflect.Float64:
		if value == "" {
			field.SetFloat(0)
			return nil
		}
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("unmarshalcsv, failed to parse %s as float: %w", field.Type(), err)
		}
		field.SetFloat(floatValue)
	case reflect.Bool:
		if value == "" {
			field.SetBool(false)
			return nil
		}
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("unmarshalcsv, failed to parse %s as boolean: %w", field.Type(), err)
		}
		field.SetBool(boolValue)
	default:
		return fmt.Errorf("unmarshalcsv: unsupported type %s", field.Type())
	}
	return nil
}
