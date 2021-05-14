package core

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

func defaultFileZillaIniFilePath() string {
	usr, err := user.Current()
	if err != nil {
		//log.Fatal(err)
	}
	return usr.HomeDir + "\\AppData\\Roaming\\FileZilla\\sitemanager.xml"
}

func GetFileZillaPasswords() {
	xmlFile, err := os.Open(defaultFileZillaIniFilePath())

	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()

	type Server struct {
		XMLName xml.Name `xml:"Server"`
		Host    string   `xml:"Host"`
		User    string   `xml:"User"`
		Pass    string   `xml:"Pass"`
	}

	type result struct {
		Hostname, Username, Password string
	}

	type Servers struct {
		XMLName xml.Name `xml:"Servers"`
		Servers []Server `xml:"Server"`
	}

	type FileZilla3 struct {
		XMLName xml.Name `xml:"FileZilla3"`
		Servers Servers  `xml:"Servers"`
	}

	// defer the closing of our xmlFile so that we can parse it later on

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var fileZilla3 FileZilla3

	xml.Unmarshal(byteValue, &fileZilla3)
	err = os.Remove("results/fileZilla_password.json")
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.OpenFile("results/fileZilla_password.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	f.WriteString("[\n")
	for i := 0; i < len(fileZilla3.Servers.Servers); i++ {
		sDec, err := base64.StdEncoding.DecodeString(fileZilla3.Servers.Servers[i].Pass)
		if err != nil {
			fmt.Printf("Error decoding string: %s ", err.Error())
			return
		}
		// fmt.Println("HostName: " + fileZilla3.Servers.Servers[i].Host)
		// fmt.Println("Username: " + fileZilla3.Servers.Servers[i].User)
		// fmt.Println("Password: " + string(sDec))

		finalJson := result{
			Hostname: fileZilla3.Servers.Servers[i].Host,
			Username: fileZilla3.Servers.Servers[i].User,
			Password: string(sDec),
		}

		file, _ := json.MarshalIndent(finalJson, "", " ")

		if err != nil {
			panic(err)
		}
		defer f.Close()
		if _, err = f.Write(file); err != nil {
			panic(err)
		}
		if i+1 != len(fileZilla3.Servers.Servers) {
			f.WriteString(",")
		}
		f.WriteString("\n")

	}
	f.WriteString("]")
}
