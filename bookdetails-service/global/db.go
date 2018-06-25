package global

import (
	"fmt"
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/openzipkin/zipkin-go"
	"time"
)

type db struct {
	*gorm.DB
}

func newBookDB() *db {
	conf := Conf.DB_BOOK
	entity, err := gorm.Open(
		conf.Driver,
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
			conf.Username,
			conf.Pwd,
			conf.Host,
			conf.Port,
			conf.DBName,
			conf.Charset,
			conf.ParseTime,
			conf.Local,
		),
	)

	if err != nil {
		Logger.Fatalln("mysql conn failed,", err)
	}

	return &db{
		entity,
	}
}

func (this *db) WarpRawScan(ctx context.Context, dest interface{}, sql string, values ...interface{}) (*gorm.DB) {
	//span, _, err := this.zipkin(
	//	ctx,
	//	zipkinOption{
	//		OptionType: ZIPKIN_OPTION_TAG,
	//		zipkinTag:  ZipkinTag{"sql", sql},
	//	},
	//	zipkinOption{
	//		OptionType: ZIPKIN_OPTION_TAG,
	//		zipkinTag:  ZipkinTag{"values", fmt.Sprint(values)},
	//	},
	//)
	//if err == nil {
	//	defer func() {
	//		span.Annotate(time.Now(), "out db")
	//		span.Finish()
	//	}()
	//}

	span, _ := ZPTracer.StartSpanFromContext(
		ctx,
		"fetch data from db",
		zipkin.Tags(map[string]string{"func": "db"}),
	)
	span.Annotate(time.Now(), "start fetch db...")
	defer func() {
		span.Annotate(time.Now(), "ending from db...")
		span.Finish()
	}()

	return this.DB.Raw(sql, values).Scan(dest)
}
