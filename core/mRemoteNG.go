package core

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

func defaultmRemoteNGIniFilePath() string {
	usr, err := user.Current()
	if err != nil {
		//log.Fatal(err)
	}
	return usr.HomeDir + "\\AppData\\Roaming\\mRemoteNG\\confCons.xml"
}

func GetmRemoteNGPasswords() {
	xmlFile, err := os.Open(defaultmRemoteNGIniFilePath())
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()

	type Node struct {
		XMLName  xml.Name `xml:"Node"`
		Password string   `xml:"Password,attr"`
	}

	type Mrng struct {
		XMLName xml.Name `xml:"Connections"`
		Nodes   []Node   `xml:"Node"`
	}

	// defer the closing of our xmlFile so that we can parse it later on

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var mrng Mrng

	xml.Unmarshal(byteValue, &mrng)
	for i := 0; i < len(mrng.Nodes); i++ {
		fmt.Println("Password: " + mrng.Nodes[i].Password)
		// a, err := base64.StdEncoding.DecodeString(mrng.Nodes[i].Password)
		// if err != nil {
		// 	panic(err)
		// }
		// ciphertext := a
		// saltlen := 16
		// salt := ciphertext[:saltlen]
		// iv := ciphertext[saltlen : aes.BlockSize+saltlen]
		// key := pbkdf2.Key([]byte("mR3m"), salt, 1000, 32, sha1.New)

		// block, err := aes.NewCipher(key)
		// if err != nil {
		// 	panic(err)
		// }

		// // if len(ciphertext) < aes.BlockSize {
		// // 	return ""
		// // }

		// decrypted := ciphertext[saltlen+aes.BlockSize:]
		// stream := cipher.NewCFBDecrypter(block, iv)
		// stream.XORKeyStream(decrypted, decrypted)
		// fmt.Println(string(decrypted))

		// pass, err := base64.StdEncoding.DecodeString(mrng.Nodes[i].Password)
		// if err != nil {
		// 	fmt.Printf("Error decoding string: %s ", err.Error())
		// 	return
		// }
		// fmt.Printf(string(pass))
		// salt := pass[0:16]
		// associated_data := pass[0:16]
		// nonce := pass[16:32]
		// ciphertext := pass[32 : len(pass)-16]
		// tag = pass[len(pass)-16:]
		// dk := pbkdf2.Key([]byte("mR3m"), []byte(salt), 1000, 32, sha1.New)
		// block, err := aes.NewCipher(dk)
		// if err != nil {
		// 	panic(err.Error())
		// }

		// aesgcm, err := cipher.NewGCM(block)
		// if err != nil {
		// 	panic(err.Error())
		// }

		// plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
		// if err != nil {
		// 	panic(err.Error())
		// }
		// fmt.Printf("%s\n", plaintext)
		//key = hashlib.pbkdf2_hmac("sha1", args.password.encode(), salt, 1000, dklen=32)
	}
}
