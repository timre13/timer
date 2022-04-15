package confreader

import (
    "timer/common"
    "fmt"
    "os"
    "encoding/json"
)

// ---------- Config defaults ----------
var CONF_DEF_WORK_SESS_LEN_MIN      = 25
var CONF_DEF_BREAK_SESS_LEN_MIN     = 5
var CONF_DEF_AUTO_START_WORK_SESS   = true
var CONF_DEF_AUTO_START_BREAK_SESS  = true
var CONF_DEF_SESS_END_SHOW_NOTIF    = true
var CONF_DEF_SHOW_REM_TIME_IN_TITLE = true

type Config struct {
    WorkSessDurMin          int
    AutoStartWorkSess       bool

    BreakSessDurMin         int
    AutoStartBreakSess      bool

    SessEndShowNotif        bool
    ShowRemTimeInWinTitle   bool
}

func LoadConf(path string) Config {
    fmt.Printf("ConfReader :: Loading config file: \"%s\"\n", path)

    conf := Config{}
    conf.AutoStartWorkSess      = CONF_DEF_AUTO_START_WORK_SESS
    conf.AutoStartBreakSess     = CONF_DEF_AUTO_START_BREAK_SESS
    conf.SessEndShowNotif       = CONF_DEF_SESS_END_SHOW_NOTIF
    conf.ShowRemTimeInWinTitle  = CONF_DEF_SHOW_REM_TIME_IN_TITLE

    fileContent, err := os.ReadFile(path)
    if err != nil {
        fmt.Println("ConfReader :: Failed to load config file:", err)
    } else {
        err = json.Unmarshal(fileContent, &conf)
        if err != nil {
            fmt.Println("ConfReader :: Failed to parse JSON:", err)
        } else {
            fmt.Printf("ConfReader :: Loaded config: %+v\n", conf)
        }
    }

    /*
        Correct invalid or missing values.
    */
    if conf.WorkSessDurMin < 1 {
        conf.WorkSessDurMin = CONF_DEF_WORK_SESS_LEN_MIN
        fmt.Println("ConfReader :: Using default value for `WorkSessDurMin`:", conf.WorkSessDurMin)
    }
    if conf.BreakSessDurMin < 1 {
        conf.BreakSessDurMin = CONF_DEF_BREAK_SESS_LEN_MIN
        fmt.Println("ConfReader :: Using default value for `BreakSessDurMin`:", conf.BreakSessDurMin)
    }

    return conf
}

func WriteConf(path string, conf *Config) {
    fmt.Printf("ConfReader :: Writing config file: \"%s\"\n", path)

    jsonData, err := json.MarshalIndent(&conf, "", "    ")
    if err != nil {
        fmt.Println("ConfReader :: Failed to build JSON:", err)
    } else {
        err = os.WriteFile(path, jsonData, 0o644)
        if err != nil {
            fmt.Println("ConfReader :: Failed to write config file:", err)
        } else {
            fmt.Printf("ConfReader :: Wrote config: %+v\n", *conf)
        }
    }
}

func (c *Config) GetSessLenMs(typ common.SessionType) int {
    switch typ {
    case common.SESSION_TYPE_WORK:      return common.MinsToMillisecs(c.WorkSessDurMin)
    case common.SESSION_TYPE_BREAK:     return common.MinsToMillisecs(c.BreakSessDurMin)
    }
    panic(typ)
}
