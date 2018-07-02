package global

import (
	"os"
	"log"

	"github.com/joho/godotenv"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/davecgh/go-spew/spew"
	"github.com/openzipkin/zipkin-go"
	"bookinfo/bookdetails-service/models"
	"google.golang.org/grpc"
)

type conf struct {
	ServiceName string    `yaml:"service_name"`
	Redis       redisConf
	Zipkin      zipkinConf
	DB_BOOK     mysqlConf `yaml:"db_book"`
	HttpServer struct {
		Addr string
	} `yaml:"http_server"`
	GrpcServer struct {
		Addr string
	} `yaml:"grpc_server"`
	DebugServer struct {
		Addr string
	} `yaml:"debug_server"`
	MetricsServer struct {
		Addr string
	} `yaml:"metrics_server"`
	Servers struct {
		BookComments server `yaml:"book_comments"`
	}
}

type redisConf struct {
	Addr string
	Pwd  string
	DB   int
}
type mysqlConf struct {
	Username        string
	Pwd             string
	Host            string
	Port            int
	DBName          string `yaml:"db_name"`
	Driver          string
	Charset         string
	ParseTime       string `yaml:"parse_time"`
	Local           string
	ConnMaxLifeTime int    `yaml:"conn_max_life_time"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
}

type zipkinConf struct {
	Addr        string
	ServiceName string `yaml:"service_name"`
	Reporter struct {
		Timeout       int
		BatchSize     int `yaml:"batch_size"`
		BatchInterval int `yaml:"batch_interval"`
	}
}

type server struct {
	Grpc string
	Http string
}

var ProjectPath = os.Getenv("GOPATH") + "/src/bookinfo/bookdetails-service"
var RuntimePath = ProjectPath + "/runtime"
var LogPath = RuntimePath + "/logs"

const (
	RUN_MODE_LOCAL     = "local"
	RUN_MODE_CONTAINER = "container"
)

var Conf conf

var Logger logger

var BOOK_DB *db

var ZPTracer *zipkin.Tracer

var GrpcOpts []grpc.ServerOption

func init() {
	loadConf()
	Logger = newLogger()
	BOOK_DB = newBookDB()
	ZPTracer = newZPTracer()
	GrpcOpts = loadGrpcOpts()
	models.Migrate(BOOK_DB.DB)
}

func loadConf() {

	os.MkdirAll(LogPath, os.ModePerm)

	log.Println(ProjectPath)

	if err := godotenv.Load(ProjectPath + "/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	runMode := os.Getenv("RUN_MODE")
	log.Println("run mode:", runMode)

	var confFile string
	switch runMode {
	case RUN_MODE_LOCAL:
		confFile = ProjectPath + "/conf/local.yaml"
	case RUN_MODE_CONTAINER:
		confFile = ProjectPath + "/conf/container.yaml"
	default:
		log.Fatalln("unsuppoer run mode! supports:[local,container]")
	}

	conf, _ := ioutil.ReadFile(confFile)
	if err := yaml.Unmarshal(conf, &Conf); err != nil {
		log.Fatalln("conf load failed", err)
	}

	spew.Dump(Conf)

}
