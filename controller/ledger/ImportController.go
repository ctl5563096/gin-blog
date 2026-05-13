package ledger

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"gin-blog/models/ledger"
	"gin-blog/pkg/app"
	"io"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type xlsxWorksheet struct {
	Rows []xlsxRow `xml:"sheetData>row"`
}

type xlsxRow struct {
	Index int        `xml:"r,attr"`
	Cells []xlsxCell `xml:"c"`
}

type xlsxCell struct {
	Ref          string               `xml:"r,attr"`
	Type         string               `xml:"t,attr"`
	Value        string               `xml:"v"`
	InlineString xlsxSharedStringItem `xml:"is"`
}

type xlsxSharedStrings struct {
	Items []xlsxSharedStringItem `xml:"si"`
}

type xlsxSharedStringItem struct {
	Texts     []string `xml:"t"`
	RichTexts []string `xml:"r>t"`
}

type ledgerWorkbookImport struct {
	Records       []ledger.CreateTransaction
	Skipped       int
	FinalBalance  float64
	BalanceSynced bool
}

func ImportTransactions(c *gin.Context) {
	ownerId := getLedgerOwnerId(c)
	if ownerId == "" {
		app.MissToken(c)
		return
	}
	account, ok := getLedgerAccount(c, ownerId)
	if !ok {
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		app.FailWithParameter("请选择 Excel 文件", c)
		return
	}
	if strings.ToLower(filepath.Ext(fileHeader.Filename)) != ".xlsx" {
		app.FailWithParameter("仅支持 .xlsx 文件", c)
		return
	}

	importData, err := parseLedgerWorkbook(fileHeader, account.Id, ownerId)
	if err != nil {
		app.FailWithMessage(err.Error(), 1, c)
		return
	}
	if len(importData.Records) == 0 {
		app.FailWithMessage("没有可导入的记账记录", 1, c)
		return
	}

	success, err := ledger.CreateRecords(importData.Records)
	if err != nil {
		app.FailWithMessage("导入记账记录失败", 1, c)
		return
	}
	if importData.BalanceSynced {
		if err = ledger.SyncSavingCurrentAmount(account, ownerId, importData.FinalBalance, "Excel余额同步"); err != nil {
			app.FailWithMessage("导入成功，余额同步失败", 1, c)
			return
		}
	}
	app.OkWithData(ledger.ImportTransactionResult{
		Total:         len(importData.Records) + importData.Skipped,
		Success:       success,
		Skipped:       importData.Skipped,
		BalanceSynced: importData.BalanceSynced,
		FinalBalance:  importData.FinalBalance,
	}, c)
}

func parseLedgerWorkbook(fileHeader *multipart.FileHeader, accountId int, ownerId string) (ledgerWorkbookImport, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return ledgerWorkbookImport{}, fmt.Errorf("Excel 文件打开失败")
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return ledgerWorkbookImport{}, fmt.Errorf("Excel 文件读取失败")
	}
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return ledgerWorkbookImport{}, fmt.Errorf("Excel 文件格式错误")
	}

	sharedStrings, err := readSharedStrings(reader)
	if err != nil {
		return ledgerWorkbookImport{}, err
	}
	sheetData, err := readZipFile(reader, "xl/worksheets/sheet1.xml")
	if err != nil {
		return ledgerWorkbookImport{}, fmt.Errorf("未找到总账工作表")
	}

	var sheet xlsxWorksheet
	if err = xml.Unmarshal(sheetData, &sheet); err != nil {
		return ledgerWorkbookImport{}, fmt.Errorf("总账工作表解析失败")
	}

	records := make([]ledger.CreateTransaction, 0)
	skipped := 0
	finalBalance := 0.0
	balanceSynced := false
	for _, row := range sheet.Rows {
		if row.Index <= 3 {
			continue
		}
		cellMap := map[string]string{}
		for _, cell := range row.Cells {
			cellMap[cellColumn(cell.Ref)] = cellText(cell, sharedStrings)
		}

		dateText := strings.TrimSpace(cellMap["B"])
		source := strings.TrimSpace(cellMap["C"])
		income := parseAmount(cellMap["D"])
		expense := parseAmount(cellMap["E"])
		balanceText := strings.TrimSpace(cellMap["F"])
		balance := parseAmount(balanceText)
		note := strings.TrimSpace(cellMap["G"])
		importFlag := strings.TrimSpace(cellMap["H"])
		if note == "" {
			note = "Excel导入"
		}
		if dateText == "" || source == "" || (income <= 0 && expense <= 0) {
			skipped++
			continue
		}

		recordType := "expense"
		amount := expense
		if income > 0 {
			recordType = "income"
			amount = income
		}
		if importFlag == "2" {
			if income <= 0 {
				skipped++
				continue
			}
			recordType = "income"
			amount = income
		}
		if importFlag == "3" {
			if expense <= 0 {
				skipped++
				continue
			}
			recordType = "expense"
			amount = expense
		}

		date, ok := normalizeExcelDate(dateText)
		if !ok {
			skipped++
			continue
		}

		records = append(records, ledger.CreateTransaction{
			AccountId:                 accountId,
			OwnerId:                   ownerId,
			Type:                      recordType,
			Amount:                    amount,
			Category:                  limitRunes(source, 32),
			Description:               limitRunes(source, 80),
			Note:                      limitRunes(note, 80),
			Date:                      date,
			PaymentMethod:             "导入",
			MemberName:                "我",
			SaveAsCategoryConfig:      importFlag == "1",
			SaveAsMonthlySavingConfig: importFlag == "2" || importFlag == "3",
			ImportedSavingBalance:     balance,
			HasImportedSavingBalance:  balanceText != "",
		})
		if balanceText != "" {
			finalBalance = balance
			balanceSynced = true
		}
	}
	return ledgerWorkbookImport{
		Records:       records,
		Skipped:       skipped,
		FinalBalance:  finalBalance,
		BalanceSynced: balanceSynced,
	}, nil
}

func readSharedStrings(reader *zip.Reader) ([]string, error) {
	data, err := readZipFile(reader, "xl/sharedStrings.xml")
	if err != nil {
		return []string{}, nil
	}
	var shared xlsxSharedStrings
	if err = xml.Unmarshal(data, &shared); err != nil {
		return nil, fmt.Errorf("Excel 共享文本解析失败")
	}
	result := make([]string, 0, len(shared.Items))
	for _, item := range shared.Items {
		result = append(result, item.String())
	}
	return result, nil
}

func readZipFile(reader *zip.Reader, name string) ([]byte, error) {
	for _, file := range reader.File {
		if file.Name != name {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		return io.ReadAll(rc)
	}
	return nil, fmt.Errorf("file not found")
}

func cellText(cell xlsxCell, sharedStrings []string) string {
	value := strings.TrimSpace(cell.Value)
	if cell.Type == "inlineStr" {
		return cell.InlineString.String()
	}
	if cell.Type != "s" {
		return value
	}
	index, err := strconv.Atoi(value)
	if err != nil || index < 0 || index >= len(sharedStrings) {
		return ""
	}
	return sharedStrings[index]
}

func (item xlsxSharedStringItem) String() string {
	if len(item.Texts) > 0 {
		return strings.Join(item.Texts, "")
	}
	return strings.Join(item.RichTexts, "")
}

func cellColumn(ref string) string {
	var builder strings.Builder
	for _, r := range ref {
		if r >= 'A' && r <= 'Z' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func parseAmount(value string) float64 {
	value = strings.TrimSpace(strings.ReplaceAll(value, ",", ""))
	if value == "" || value == "/" || value == "-" {
		return 0
	}
	amount, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return amount
}

func normalizeExcelDate(value string) (string, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", false
	}
	if serial, err := strconv.ParseFloat(value, 64); err == nil {
		date := time.Date(1899, 12, 30, 0, 0, 0, 0, time.Local).AddDate(0, 0, int(serial))
		return date.Format("2006-01-02"), true
	}
	layouts := []string{"2006-01-02", "2006/1/2", "2006/01/02"}
	for _, layout := range layouts {
		if date, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return date.Format("2006-01-02"), true
		}
	}
	return "", false
}

func limitRunes(value string, max int) string {
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	return string(runes[:max])
}
