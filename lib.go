package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	dice "git.cmcode.dev/cmcode/go-dicewarelib"
	"github.com/adrg/xdg"
)

// Used for the config file directory and other things.
const APP_NAME = "go-fltk-diceware"

// Name of the config file.
const CONFIG_FILE = "config.json"

// Wraps around log.Println() as well as adding activity to the
// activity text buffer. Always adds a newline to the activity buffer.
func Log(v ...any) {
	// activityText = fmt.Sprintf("%v<p>%v</p>", activityText, fmt.Sprintf("%v", v...))
	// activity.SetValue(activityText)
	log.Println(v...)
	// activity.SetTopLine(MAX_TOP_LINE)
	// activity.SetTopLine(activity.TopLine() - activity.H()) // scroll to the bottom
}

// Wraps around log.Printf() as well as adding activity to the
// activity text buffer. Always adds a newline to the activity buffer.
func Logf(format string, v ...any) {
	format = fmt.Sprintf("%v\n", format)
	// activityText = fmt.Sprintf("%v<p>%v</p>", activityText, fmt.Sprintf(format, v...))
	// activity.SetValue(activityText)
	log.Printf(format, v...)
	// activity.SetTopLine(MAX_TOP_LINE)
	// activity.SetTopLine(activity.TopLine() - activity.H()) // scroll to the bottom
}

/*
// Returns a minimum value of 0 if the provided integer is less than 0.
// Otherwise, returns the int itself.
func floorz(i int) int {
	if i < 0 {
		return 0
	}
	return i
}

// Returns the lesser of the two numbers, with a floor of zero.
func minz(a, b int) int {
	if a < b {
		return floorz(a)
	}
	return floorz(b)
}

// Replaces all occurences of any keys in the secrets map with their masked
// values.
func obscure(s string, secrets map[string]string) string {
	r := s
	for k, v := range secrets {
		if k == "" {
			continue
		}

		if v == "" {
			v = strings.Repeat("*", len(k))
		}

		r = strings.ReplaceAll(r, k, v)
	}

	return r
} */

/*
func encr(s, key string) (string, error) {
	keyb := []byte(key)
	block, err := aes.NewCipher(keyb)
	if err != nil {
		return "", err
	}

	// GCM mode requires a nonce (number used once)
	nonce := make([]byte, 12) // GCM standard nonce size is 12 bytes
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Create a GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Encrypt the plaintext
	ciphertext := gcm.Seal(nonce, nonce, []byte(s), nil)

	// Return the base64 encoded ciphertext
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decr(s, key string) (string, error) {
	// Decode the base64 encoded ciphertext
	ciphertextBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher
	keyb := []byte(key)
	block, err := aes.NewCipher(keyb)
	if err != nil {
		return "", err
	}

	// GCM mode requires a nonce (number used once)
	nonce, ciphertextBytes := ciphertextBytes[:12], ciphertextBytes[12:]

	// Create a GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Decrypt the ciphertext
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// Retrieves the encrypted keys (secrets) from the provided map `s`, decrypts
// each of them, and puts their unencrypted values into the resulting map.
func getDecryptedSecrets(s map[string]string, key string) (map[string]string, error) {
	r := make(map[string]string)
	if s == nil {
		return r, fmt.Errorf("received nil map")
	}

	for k, v := range s {
		decrypted, err := decr(k, key)
		if err != nil {
			return r, fmt.Errorf("failed to decrypt: %v", err.Error())
		}

		r[decrypted] = v
	}

	return r, nil
}
*/

func (app *App) loadConfig() {
	var err error
	if app.configFilePath == "" {
		app.configFilePath, err = xdg.SearchConfigFile(path.Join(APP_NAME, CONFIG_FILE))
		if err != nil {
			log.Printf("failed to get xdg config dir: %v", err.Error())
		}
	}

	if app.configFilePath != "" {
		bac, err := os.ReadFile(app.configFilePath)
		if err != nil {
			log.Printf("config file not readable at %v", app.configFilePath)
		}

		err = json.Unmarshal(bac, app.conf)
		if err != nil {
			log.Printf("config file %v failed to parse: %v", app.configFilePath, err.Error())
		}

		log.Printf("loaded config from %v", app.configFilePath)
	} else {
		if xdg.ConfigHome != "" {
			app.configFilePath = path.Join(xdg.ConfigHome, APP_NAME, CONFIG_FILE)
			log.Printf("using %v for config file path", app.configFilePath)
		} else {
			log.Println("unable to automatically identify any suitable config dirs; configuration will not be saved")
		}
	}
}

func (app *App) saveConfig() error {
	if app.configFilePath == "" {
		return fmt.Errorf("received empty config filename")
	}

	if app.conf == nil {
		return fmt.Errorf("config was nil")
	}

	b, err := json.Marshal(app.conf)
	if err != nil {
		return fmt.Errorf("failed to marshal app config to yaml: %v", err.Error())
	} else {
		dir, _ := filepath.Split(app.configFilePath)
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			return fmt.Errorf("failed to create app config parent dir %v: %v", dir, err.Error())
		}
		err = os.WriteFile(app.configFilePath, b, 0o644)
		if err != nil {
			return fmt.Errorf("failed to save app config to %v: %v", app.configFilePath, err.Error())
		}
	}

	return nil
}

// Initializes the diceware library with the stored word lists. Can be executed
// repeatedly.
func (app *App) initDice() {
	simple, scount := dice.GetWords(content, "words-simple.txt")
	app.words.Simple = &simple
	app.words.SimpleCount = scount
	if app.conf.Extra {
		complex, ccount := dice.GetWords(content, "words-complex.txt")
		app.words.Complex = &complex
		app.words.ComplexCount = ccount
	} else {
		app.words.Complex = &map[int]string{} // zero out the ram usage
		app.words.ComplexCount = 0
	}
	Logf("loaded %v simple words and %v complex words", app.words.SimpleCount, app.words.ComplexCount)
	// log.Println(dice.GeneratePassword(&app.words, 3, " ", 64, 20, false))
}
