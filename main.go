package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/hashicorp/hcl"
)

func main() {
	Title := figure.NewFigure("TF AUDIT", "", true)
	Title.Print()

	FilePath := flag.String("d", "/var/foo", "Base project directory")
	Verbose := flag.Bool("v", false, "Display each hibernate file location inside the table -default:false")
	Tfversion := flag.Bool("tf", false, "Search .terraform-version and versions.tf file in project -default:false")
	flag.Parse()

	//check mandatory -d flag
	required := []string{"d"}
	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			// or possibly use `log.Fatalf` instead of:
			fmt.Fprintf(os.Stderr, "missing required -%s argument\n", req)
			os.Exit(2) // the same exit code flag.Parse uses
		}
	}

	LoadConfigs(*FilePath,*Verbose,*Tfversion) 
}

func LoadConfigs(root string,verbose bool,tfversion bool) ( err error) {
	files, _ := WalkMatch(root, "*.tf")
	for _, file := range files {



		//ignore /.terraform/ directory
		if !strings.Contains(file, "/.terraform/"){





			//full path
			parent := filepath.Dir(file)
			//containing directory
			dir := filepath.Base(parent)
			
			if strings.Contains(file, "main.tf") {
				fmt.Println("---------------------------------------------------------")
				//Print current project
				color.Green(dir)
				if verbose {
					//print full path
					color.Yellow(parent)
				}
	
				if tfversion {
					if _, err := os.Stat(filepath.Join(parent, ".terraform-version")); err == nil {
						// path/to/whatever exists
						FileContent, err := ioutil.ReadFile(filepath.Join(parent, ".terraform-version")) // just pass the file name
						if err != nil {
							fmt.Print(err)
						}
						fmt.Println(".terraform-version " , string(FileContent))
					  }

	
				}


				cmd := exec.Command("terraform", "--version")
				cmd.Dir = parent
				out, err_tv := cmd.Output()

				  if err_tv != nil{
					fmt.Print(err_tv)
				  }
				 
				  outString := string(out)
				  temp := strings.Split(outString,"\n")
	
				  fmt.Println(temp[0])
	
			} 

			if tfversion {
				//version constraints
				if strings.Contains(file, "versions.tf") {// true ""{	
					ReadVersion(file)
				} 
			}
		}  
		

	}
	return
}

func ReadVersion(path string){
	FileContent, err := ioutil.ReadFile(path) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	out := make(map[string]map[string]interface{})
	FileContentString := string(FileContent) // convert content to a 'string'
	err = hcl.Decode(&out, FileContentString)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("versions.tf " , out["terraform"]["required_version"])
}


//WalkMatch get all files on a given path for specific file extention
func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

