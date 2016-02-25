package main

import (
	"errors"
	"fmt"
	"gobackup/Godeps/_workspace/src/github.com/codegangsta/cli"
	"gobackup/Godeps/_workspace/src/github.com/kyawmyintthein/barkup"
	"gobackup/Godeps/_workspace/src/gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Adapter int

const (
	MYSQL Adapter = 1 + iota
	POSTGRE
	MONGODB
)

var adapters = [...]string{
	"mysql",
	"postgre",
	"mongodb",
}

type Config struct {
	Adapter  string   `json:"adapter"`
	Database Database `json:"database"`
	S3       S3       `json:"s3"`
}

type Database struct {
	Name     string `default:"test" json:"name"`
	User     string `default:"root"`
	Password string `required:"true" env:"DBPassword"`
	Port     string `default:"3306"`
	Host     string `default:"localhost"`
}

type S3 struct {
	Bucket       string `required:"true" json:"name"`
	Region       string `required:"true" json:"region"`
	AccessKey    string `required:"true" env:"accessKey" json:"access_key"`
	ClientSecret string `required:"true" env:"accessKey" json:"secret"`
	Path         string `required:"true"  json:"path"`
}

var config Config

func main() {
	app := cli.NewApp()
	app.Name = "gobackup"
	app.Usage = "gobackup -f /path/to/config.yml OR gobackup -adapter=mysql,mongodb -h=localhost -p=3303 -db=test -user=root -password=root -target=s3 -target-path=/path/to/s3_config.yml"

	// Commands
	app.Commands = []cli.Command{
		{
			Name: "export",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "f, file",
					Value: "gobackup.yml",
					Usage: "Specify an alternate gobackup config file (default: gobackup.yml)",
				},
				cli.StringFlag{
					Name:  "adapter",
					Usage: "Specify an alternate database adapter (default: none)",
				},
			},
			Usage: "Export database and upload to target source",
			Action: func(c *cli.Context) {
				log.Println("completed task: ", c.Args().First())
				// if len(c.Args()) > 0 {
				//  name =  c.Args().First()
				// }
				file := c.String("file")
				if file == "" {
					file = c.String("f")
					if file == "" {
						file = "gobackup.yml"
					}
				}

				dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
				if err != nil {
					panic(err)
				}
				fileName := dir + "/" + file
				yamlFile, err := ioutil.ReadFile(fileName)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(string(yamlFile))

				err = yaml.Unmarshal(yamlFile, &config)
				fmt.Println(err)
				if err != nil {
					log.Fatal(err)
				}

				if err := config.validConfig(); err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Adapter: %#v\n", config.Adapter)
				switch config.Adapter {
				case MYSQL.String():
					_, err := config.exportMysqlDatabase()
					fmt.Println(err)
					if err == nil {
						log.Fatal(err)
					}
				case POSTGRE.String():
					config.exportMysqlDatabase()
				case MONGODB.String():
					config.exportMongodbDatabase()
				}

			},
		},
	}
	app.Run(os.Args)
}

func (config *Config) validConfig() error {
	if config.Adapter == "" {
		return errors.New("Adapter not found in config file.")
	}

	// if config.Database{
	// 	return errors.New("math: square root of negative number")
	// }

	// if config.S3 == nil{
	// 	return errors.New("math: square root of negative number")
	// }

	return nil
}

func (config *Config) exportMysqlDatabase() (bool, error) {
	// Configure a MySQL exporter
	mysql := &barkup.MySQL{
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		DB:       config.Database.Name,
		User:     config.Database.User,
		Password: config.Database.Password,
	}

	// Configure a S3 storer
	s3 := &barkup.S3{
		Region:       config.S3.Region,
		Bucket:       config.S3.Bucket,
		AccessKey:    config.S3.AccessKey,
		ClientSecret: config.S3.ClientSecret,
	}

	// Export the database, and send it to the
	// bucket in the `db_backups` folder
	err := mysql.Export().To(config.S3.Path, s3)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (config *Config) exportPostgreDatabase() (bool, error) {
	log.Println("Export Postgre")
	return true, nil
}

func (config *Config) exportMongodbDatabase() (bool, error) {
	log.Println("Export mongodb")
	return true, nil
}

func (adapter Adapter) String() string {
	return adapters[adapter-1]
}
