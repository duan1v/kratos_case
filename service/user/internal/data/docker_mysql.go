package data

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ory/dockertest/v3"
)

func CommonMysql() string {
	conStr := fmt.Sprintf("dywily:q123456we@tcp(localhost:%s)/mysql?parseTime=true", "3306")
	return conStr

}

func DockerMysql(img, version string) (string, func()) {
	return innerDockerMysql(img, version)
}

// 初始化 Docker mysql 容器
func innerDockerMysql(img, version string) (string, func()) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	pool.MaxWait = time.Minute * 2
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run(img, version, []string{"MYSQL_ROOT_PASSWORD=secret", "MYSQL_ROOT_HOST=%"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	conStr := fmt.Sprintf("root:secret@(localhost:%s)/mysql?parseTime=true", resource.GetPort("3306/tcp"))

	if err := pool.Retry(func() error {
		var err error
		db, err := sql.Open("mysql", conStr)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// 回调函数关闭容器
	return conStr, func() {
		if err = pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}
}
