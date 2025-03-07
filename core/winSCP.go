package core

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"strconv"

	"github.com/go-ini/ini"
)

const (
	PW_MAGIC = 0xA3
	PW_FLAG  = 0xFF
)

func GetWinSCPPasswords() {
	/*if len(args) != 3 && len(args) != 2 {
		fmt.Println("WinSCP stored password finder\n")
		fmt.Println("Registry:")
		fmt.Println("  Open regedit and navigate to [HKEY_CURRENT_USER\\Software\\Martin Prikryl\\WinSCP 2\\Sessions] to get the hostname, username and encrypted password\n")
		if runtime.GOOS == "windows" {
			fmt.Println("  Usage winscppasswd.exe <host> <username> <encrypted_password>")
		} else {
			fmt.Println("  Usage ./winscppasswd <host> <username> <encrypted_password>")
		}
		fmt.Println("\n\nWinSCP.ini:")
		if runtime.GOOS == "windows" {
			fmt.Println("  Usage winscppasswd.exe ini [<filepath>]")
		} else {
			fmt.Println("  Usage ./winscppasswd ini [<filepath>]")
		}
		fmt.Printf("  Default value <filepath>: %s\n", defaultWinSCPIniFilePath())
		return
	}
	if args[0] == "ini" {
		if len(args) == 2 {
			decryptIni(args[1])
		} else {
			decryptIni(defaultWinSCPIniFilePath())
		}
	} else {
		fmt.Println(decrypt(args[0], args[1], args[2]))
	}*/
	err := os.Remove("results/winscp_password.json")
	if err != nil {
		fmt.Println(err)
	}
	decryptIni(defaultWinSCPIniFilePath())
}

type winSCP struct {
	Hostname, Username, Password string
}

func defaultWinSCPIniFilePath() string {
	usr, err := user.Current()
	if err != nil {
		//log.Fatal(err)
	}
	return usr.HomeDir + "\\AppData\\Roaming\\winSCP.ini"
}

func decryptIni(filepath string) {
	cfg, err := ini.InsensitiveLoad(filepath)
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile("results/winscp_password.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	filecontent := "[\n"
	//f.WriteString("[\n")
	for _, c := range cfg.Sections() {
		if c.HasKey("Password") {
			/*fmt.Printf("%s\n", strings.TrimPrefix(c.Name(), "sessions\\"))
			fmt.Printf("  Hostname: %s\n", c.Key("HostName").Value())
			fmt.Printf("  Username: %s\n", c.Key("UserName").Value())
			fmt.Printf("  Password: %s\n", decryptWinSCP(c.Key("HostName").Value(), c.Key("UserName").Value(), c.Key("Password").Value()))
			fmt.Println("========================")*/
			data := winSCP{
				Hostname: c.Key("HostName").Value(),
				Username: c.Key("UserName").Value(),
				Password: decryptWinSCP(c.Key("HostName").Value(), c.Key("UserName").Value(), c.Key("Password").Value()),
			}

			file, _ := json.MarshalIndent(data, "", " ")
			if err != nil {
				fmt.Println(err)
				return
			}
			filecontent = filecontent + string(file) + ",\n"
			if err != nil {
				panic(err)
			}
			defer f.Close()
			// if _, err = f.Write(file); err != nil {
			// 	panic(err)
			// }

			//f.WriteString(",\n")
		}
	}
	filecontent = filecontent[:len(filecontent)-2]
	filecontent = filecontent + "]"
	if _, err = f.WriteString(filecontent); err != nil {
		panic(err)
	}
	fmt.Println(filecontent)

}

func decryptWinSCP(host, username, password string) string {
	key := username + host
	passbytes := []byte{}
	for i := 0; i < len(password); i++ {
		val, _ := strconv.ParseInt(string(password[i]), 16, 8)
		passbytes = append(passbytes, byte(val))
	}
	var flag byte
	flag, passbytes = dec_next_char(passbytes)
	var length byte = 0
	if flag == PW_FLAG {
		_, passbytes = dec_next_char(passbytes)

		length, passbytes = dec_next_char(passbytes)
	} else {
		length = flag
	}
	toBeDeleted, passbytes := dec_next_char(passbytes)
	passbytes = passbytes[toBeDeleted*2:]

	clearpass := ""
	var (
		i   byte
		val byte
	)
	for i = 0; i < length; i++ {
		val, passbytes = dec_next_char(passbytes)
		clearpass += string(val)
	}

	if flag == PW_FLAG {
		clearpass = clearpass[len(key):]
	}
	return clearpass
}

func dec_next_char(passbytes []byte) (byte, []byte) {
	if len(passbytes) <= 0 {
		return 0, passbytes
	}
	a := passbytes[0]
	b := passbytes[1]
	passbytes = passbytes[2:]
	return ^(((a << 4) + b) ^ PW_MAGIC) & 0xff, passbytes
}
