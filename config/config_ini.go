package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

type ConfigIni struct {
	Configs    map[string](map[string]string)
	configFile string
}

func (o *ConfigIni) LoadConfigFromFile(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	//fmt.Println(string(data))
	err = o.parseIniFileData(string(data))
	if err != nil {
		return err
	}
	o.configFile = fname
	return nil
}

// if value not exist, log fatal and exit program
func (o *ConfigIni) GetConfigValueFatal(session, key string) (value string) {
	value, err := o.GetConfigValue(session, key)
	if err != nil {
		log.Fatalln("GetConfigValueFatal: " + err.Error())
	}
	return
}

// if value not exist, log fatal and exit program
func (o *ConfigIni) GetAndSetConfigValueFatal(session string, key string, opChangeValue func(v *string)) (value string) {
	value, err := o.GetConfigValue(session, key)
	if err != nil {
		log.Fatalln("GetConfigValueFatal: " + err.Error())
		return
	}
	opChangeValue(&value)
	// log.Println("Debug: ", value)
	// set value after opChangeValue
	o.Configs[session][key] = value
	return
}

func (o *ConfigIni) GetConfigSessionFatal(session string) (m map[string]string) {
	m, err := o.GetConfigSession(session)
	if err != nil {
		log.Fatalln("GetConfigSessionFatal: " + err.Error())
	}
	return
}

// return current loaded config file name
func (o *ConfigIni) GetConfigFileFullName() string {
	return o.configFile
}

// clear all configs
func (o *ConfigIni) Clear() {
	o.configFile = ""
	o.Configs = nil
}

// check if the config has load successful
func (o *ConfigIni) IsConfigLoaded() bool {
	return o.configFile != "" // NOTE：设计约束：如果文件名存在，则配置已加载
}

// if value not exist, return err only
func (o *ConfigIni) GetConfigValue(session, key string) (value string, err error) {
	if o.Configs == nil {
		return "", fmt.Errorf("%s", "config empty!")
	}

	if o.Configs[session] == nil {
		return "", fmt.Errorf("session %s not exist", session)
	}

	value, exist := o.Configs[session][key]
	if !exist {
		return "", fmt.Errorf("key %s not exist", key)
	}
	return value, nil
}

// get a session config map
func (o *ConfigIni) GetConfigSession(session string) (m map[string]string, err error) {
	if o.Configs == nil {
		return nil, fmt.Errorf("%s", "config empty!")
	}
	if o.Configs[session] == nil {
		return nil, fmt.Errorf("session %s not exist", session)
	}
	return o.Configs[session], nil
}

func (o *ConfigIni) parseIniFileData(data string) error {
	if data == "" {
		return fmt.Errorf("%s", "config set is empty")
	}

	o.Configs = map[string](map[string]string){}

	lines := strings.Split(data, "\n")
	validCount := 0 // valid line count
	currSessionName := ""

	currline := 0

	// prepare the regexp obj
	regexp_session := regexp.MustCompile(`^\[\s*?([^\[\]]+?)\s*?\]$`)
	regexp_keypair := regexp.MustCompile(`^(\S+)\s*?=\s*?([^#\s]+)`)
	for _, line := range lines {
		currline++
		//fmt.Println("debug: @@@", line, "@@@")

		// strip the comment content from line
		e := strings.Index(line, "#")
		if e > -1 {
			//fmt.Println("debug: line ", currline, "strip line is:", e, line[0:e])
			line = line[0:e]
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue // passthrough the empty line and comment line
		}

		// is a session name ?
		{
			ss := regexp_session.FindStringSubmatch(line)
			if len(ss) >= 2 {
				currSessionName = ss[1]
				//fmt.Println("debug: currsessionname is :" + currSessionName)
				o.Configs[currSessionName] = map[string]string{}
				validCount++
				continue
			}
		}

		// is a key pair?
		{
			ss := regexp_keypair.FindStringSubmatch(line)
			if len(ss) >= 3 {
				//				fmt.Println("debug: ", line, "-->", ss)
				key := ss[1]
				value := ss[2]
				if o.Configs[currSessionName][key] != "" {
					return fmt.Errorf("line %d: config key(%s) confict value(%s != %s)",
						currline, key, o.Configs[currSessionName][key], value)
				}
				o.Configs[currSessionName][key] = value
				validCount++
				continue
			}

		}

		// else is a bug line
		return fmt.Errorf("line %d : error ini config content", currline)

	}

	if validCount == 0 {
		return fmt.Errorf("%s", "config set is empty")
	}

	//fmt.Println("debug: validCount is ", validCount)
	return nil
}

// print the configs loaded
func (o *ConfigIni) PrintConfigs() error {
	log.Println("@Print Config:")
	log.Println("--------------------------------------------------")
	log.Printf("# Config file: %s\n", o.configFile)
	log.Printf("# Config set: {\n")
	if o.Configs == nil {
		return fmt.Errorf("config not loaded!")
	}

	for kn, vn := range o.Configs {
		log.Printf("+ [%s]\n", kn)

		if vn == nil {
			continue
		}

		for k, v := range vn {
			log.Printf("\t- %s=%s\n", k, v)
		}
	}

	log.Printf("}\n")
	log.Println("--------------------------------------------------")
	return nil
}
