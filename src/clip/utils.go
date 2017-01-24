package clip
import (
    "fmt"
    "os"
    "runtime"

    "github.com/joho/godotenv"
    "github.com/Sirupsen/logrus"
)

type Conf struct {
    Port string //Clip Running Port
}

var Logr = logrus.New()

var Config Conf

func Init() {
    //logrus need to init coz ben berwarna
    err := godotenv.Load("./.env")
    if err != nil {
        fmt.Print("Error loading .env file")
    }
    Logr.Formatter = &logrus.TextFormatter{
        ForceColors: true,
        FullTimestamp:true,
    }
    Logr.Level = logrus.DebugLevel

    //get config
    init_config()
}

//kalau mau ngelog file apa dan line berapa yang terjadi error atau ndebug, panggil fungsi iki wae
func LogLine() {
    _, file, line, _ := runtime.Caller(1)
    Logr.Info(fmt.Sprintf("Line logger : %s:%d", file, line))
}

//load config from .env file
func init_config() {
    appPort := os.Getenv("PORT")

    if appPort == "" {
        Config.Port = ":3001"
    }else{
        Config.Port = ":"+appPort
    }

}

func StringInSlice(needle string, haystack []string) bool {
    for _, b := range haystack {
        if b == needle {
            return true
        }
    }
    return false
}
