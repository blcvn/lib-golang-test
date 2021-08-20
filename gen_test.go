package main_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func main() {
	// folders := []string{
	// 	"common",
	// 	"config",
	// 	"fabric",
	// 	"generator",
	// 	"grpcclient",
	// 	"handler",
	// }

	// currentPath, err := os.Getwd()
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// for _, folder := range folders {
	// 	folderPath := path.Join(currentPath, folder)
	// 	// fmt.Println(folderPath)
	// 	createTestFolder(folderPath)
	// }

}
func createTestFolder(folderName string) error {
	files, err := FilePathWalkDir(folderName)
	if err != nil {
		log.Fatal(err)
	}
	for _, fileName := range files {
		ext := filepath.Ext(fileName)
		// fmt.Println(fileName,ext)
		if ext == ".go" && !strings.Contains(fileName, "_test") {
			absFileName, err := getAbsolutePath(fileName)
			if err != nil {
				fmt.Println("Error to get abs path: ", err)
				return err
			}
			createTestFile(absFileName)
		}
	}
	return nil
}

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
func createTestFile(fname string) error {
	// 1. read file
	fileName := filepath.Base(fname)
	folderName := filepath.Dir(fname)
	name := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	fileTestName := fmt.Sprintf("%s_test.go", name)
	fileTest := path.Join(folderName, fileTestName)

	// 2. Check test file
	_, err := os.Stat(fileTest)
	if err == nil {
		return fmt.Errorf("File %s already exists.", fileTest)
	}

	//3. Read source code file
	file, err := os.Open(fname)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	srcbuf, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return err
	}
	src := string(srcbuf)

	//4. Parse to test file
	content, err := parseContent(fileName, src)
	if err != nil {
		log.Println("Check file:", fname, " meet Error: ", err)
		return err
	}
	//5. Write test file
	testFile, err := os.Create(fileTest)
	if err != nil {
		log.Println(err)
		return err
	}
	defer testFile.Close()

	_, err = testFile.WriteString(content)
	if err != nil {
		log.Println(err)
		return err
	}
	testFile.Sync()
	return nil
}

func getAbsolutePath(fname string) (string, error) {
	absPath := fname
	if !path.IsAbs(fname) {
		currentPath, err := os.Getwd()
		if err != nil {
			log.Println(err)
			return "", err
		}
		absPath = path.Join(currentPath, fname)
	}
	return absPath, nil
}

func parseContent(fileName string, src string) (string, error) {
	content := ""
	numFunc := 0
	// file set
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, src, 0)
	if err != nil {
		log.Println(err)
		return content, err
	}

	// main inspection
	ast.Inspect(f, func(n ast.Node) bool {
		// fmt.Println(n)
		switch fn := n.(type) {

		// catching all function declarations
		// other intersting things to catch FuncLit and FuncType
		case *ast.FuncDecl:
			// print actual function name
			numFunc += 1
			if fn.Recv != nil {
				className := ""
				for _, f := range fn.Recv.List {
					className = expr(f.Type)
					className = strings.ReplaceAll(className, "*", "")
				}
				content += fmt.Sprintf("func Test_%s_%v (t *testing.T ){  \n", className, fn.Name)
			} else {
				content += fmt.Sprintf("func Test_%v(t *testing.T ){  \n", fn.Name)
			}

			content += fmt.Sprintf("  t.Log(\"Finised Test func:  Test_%v \") \n", fn.Name)
			content += fmt.Sprintf("}\n\n")
		case *ast.File:
			content += fmt.Sprintf("package %v \n\n", fn.Name)
			content += fmt.Sprintf("import (\n")
			content += fmt.Sprintf("//  \"context\"\n")
			content += fmt.Sprintf("  \"testing\"\n")
			// content += fmt.Sprintf( "  \""+ packageName + "\"\n")
			content += fmt.Sprintf(")\n\n")
		}
		return true
	})
	if numFunc > 0 {
		return content, nil
	}
	return "", errors.Errorf("Not find func")
}
func expr(e ast.Expr) (ret string) {
	switch x := e.(type) {
	case *ast.StarExpr:
		return fmt.Sprintf("%s*%v", ret, x.X)
	case *ast.Ident:
		return fmt.Sprintf("%s%v", ret, x.Name)
	case *ast.ArrayType:
		if x.Len != nil {
			log.Println("OH OH looks like homework")
			return "TODO: HOMEWORK"
		}
		res := expr(x.Elt)
		return fmt.Sprintf("%s[]%v", ret, res)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", expr(x.Key), expr(x.Value))
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", expr(x.X), expr(x.Sel))
	default:
		fmt.Printf("\nTODO HOMEWORK: %#v\n", x)
	}
	return
}

func fields(fl ast.FieldList) (ret string) {
	pcomma := ""
	for i, f := range fl.List {
		// get all the names if present
		var names string
		ncomma := ""
		for j, n := range f.Names {
			if j > 0 {
				ncomma = ", "
			}
			names = fmt.Sprintf("%s%s%s ", names, ncomma, n)
		}
		if i > 0 {
			pcomma = ", "
		}
		ret = fmt.Sprintf("%s%s%s%s", ret, pcomma, names, expr(f.Type))
	}
	return ret
}
