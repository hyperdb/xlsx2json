# xlsx2json

This is a command-line tool that outputs any worksheet from an Excel workbook in JSON format.
It was created to generate JSON for use with Svelte.

## Usage

### After Build

```bash
xlsx2json -s [config json]
```

### Interpreter

```bash
go run main.go -s [config json]
```

## Configuration

Set the parameters required for processing.
By switching this file, you can generate multiple JSON files.

```json
{
  "xlsx_dir": "(Directory name of the Excel workbook)",
  "xlsx_wb": "(Workbook name)",
  "xlsx_ws": "(Worksheet name)",
  "dist_dir": "(Directory name for JSON output)"
}
```

## Additional Notes

- The first row of the Excel worksheet is used as the label name for JSON output. Please input the actual data starting from the second row.
- The attached "review.md" is the result of a code review by ChatGPT4.1. The review results have already been reflected in the source code.

## License

The source code of this project is licensed under the MIT License. Texts and other documents except for the source code are licensed under CC BY 4.0. Please refer to the "LICENSE" file for details.
