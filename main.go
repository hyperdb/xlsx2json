package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

// JSONファイルを読み込み、mapで返却する関数
func getParams(filename string) (map[string]interface{}, error) {
	// ファイルを開く
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// JSONデコード用の変数
	var data map[string]interface{}

	// JSONをデコード
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

// Excelワークブックから指定ワークシートをJSONに変換する関数
func convertExcelToJSON(filePath, sheetName string) (string, error) {
	// Excelファイルを開く
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// ワークシートの全データを取得
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return "", err
	}

	if len(rows) < 2 {
		return "", fmt.Errorf("データが不足しています (ラベル行 + データ行)")
	}

	// 一行目をラベルとして取得
	labels := rows[0]

	// JSON変換用のスライス
	var jsonData []map[string]string

	// 二行目以降をデータとして処理
	for _, row := range rows[1:] {
		record := make(map[string]string)
		for i, cell := range row {
			if i < len(labels) { // ラベル数を超えたデータは無視
				record[labels[i]] = cell
			}
		}
		jsonData = append(jsonData, record)
	}

	// JSONに変換
	jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func notExists(name string) bool {
	_, err := os.Stat(name)
	return os.IsNotExist(err)
}

// エントリーポイント
func main() {

	// 設定ファイルを指定
	var (
		s = flag.String("s", "config.json", "設定ファイル名(json)")
	)
	flag.Parse()
	if *s == "" {
		os.Exit(1)
	}
	filename := *s

	// 設定ファイルが無ければエラー
	if notExists(filepath.Join("./", filename)) {
		log.Fatalf("設定ファイルが見つかりません: %s", filename)
	}

	// JSONを取得
	params, err := getParams(filename)
	if err != nil {
		log.Fatalf("パラメータファイルの解析に失敗: %v", err)
	}

	// Excelファイルのパスとワークシート名
	xlsxDir, ok1 := params["xlsx_dir"].(string)
	xlsxWb, ok2 := params["xlsx_wb"].(string)
	xlsxWs, ok3 := params["xlsx_ws"].(string)
	distDir, ok4 := params["dist_dir"].(string)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		log.Fatalf("設定ファイルに必要なキーが不足、または型が不正です")
	}

	filePath := filepath.Join(xlsxDir, xlsxWb)
	sheetName := xlsxWs

	// JSON変換を実行
	jsonOutput, err := convertExcelToJSON(filePath, sheetName)
	if err != nil {
		log.Fatalf("Excel->JSON変換エラー: %v", err)
	}

	// 結果をファイルに書き込み
	distPath := filepath.Join(distDir, xlsxWs+".json")
	f, err := os.Create(distPath)
	if err != nil {
		log.Fatalf("ファイル作成失敗: %v", err)
	}

	data := []byte(jsonOutput)
	count, err := f.Write(data)
	if err != nil {
		log.Fatalf("ファイル書き込み失敗: %v", err)
	}

	fmt.Printf("create %s\n", distPath)
	fmt.Printf("write %d bytes\n", count)

	f.Close()
}
