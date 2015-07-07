package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fzzy/radix/redis"
	// "io/ioutil"
	"log"
	// "net/http"
	"os"
	// "os/exec"
	"strconv"
	"strings"
	// "sync"
	"time"
	//"github.com/davecgh/go-spew/spew"
)

type TmRedisStat struct {
	UnixTime int64 `json:"UnixTime"`
	Value    int64 `json:"Value"`
}
type TmRedisBackup struct {
	KeyName string        `json:"KeyName"`
	Stats   []TmRedisStat `json:"Stats"`
}

func main() {
	var fileName = flag.String("file", "", "The backup file")
	var doRestore = flag.Bool("restore", false, "Restore")
	var redisCnString = flag.String("redis", "", "The redis connection string")
	var doForce = flag.Bool("force", false, "Overwrite existing keys")
	var doAppend = flag.Bool("append", false, "Append to existing keys")
	var startFrom = flag.Int("start", -1, "Leave data before this unix time alone")
	flag.Parse()
	fmt.Println("File is ", *fileName, "redis is", *redisCnString, "restore is", *doRestore)

	if *doRestore {
		restore(fileName, redisCnString, *doForce, *doAppend, *startFrom)
	} else {
		backup(fileName, redisCnString)
	}
}

func restore(fileName *string, redisCnString *string, doForce bool, doAppend bool, startFrom int) {
	redisClient, err := redis.DialTimeout("tcp", *redisCnString, time.Duration(10)*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	var backupList []TmRedisBackup
	bckFile, err := os.Open(*fileName)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(bckFile)
	err = decoder.Decode(&backupList)
	if err != nil {
		log.Fatal(err)
	}

	for _, bck := range backupList {
		fmt.Println("restoring", bck.KeyName, "...")
		if ! doForce {
		    currentList, err := redisClient.Cmd("lrange", bck.KeyName, -1000000, -1).List()
		    if err != nil {
			    log.Fatal(err)
		    }
		    if len(currentList) > 0 && !doAppend {
			    log.Fatal("This key already exists, and has ", len(currentList), " values - no can do")
		    }
		} else {
			err := redisClient.Cmd("del", bck.KeyName )
		    log.Println(err)
		}
		for _, stat := range bck.Stats {
			if (int(stat.UnixTime) < startFrom) {
				log.Println("skipping ", stat.UnixTime, " for ", bck.KeyName)
				continue;
			}
			r := redisClient.Cmd("rpush", bck.KeyName, fmt.Sprintf("%d", stat.UnixTime)+":"+fmt.Sprintf("%d", stat.Value))
			if r.Err != nil {
				log.Fatal("HHHH", err)
			}
		}
	}
}

func backup(fileName *string, redisCnString *string) {
	redisClient, err := redis.DialTimeout("tcp", *redisCnString, time.Duration(10)*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	var backupList []TmRedisBackup
	keyList, err := redisClient.Cmd("keys", "*:daily_*").List()
	if err != nil {
		log.Fatal(err)
	}
	for _, keyName := range keyList {
		list, err := redisClient.Cmd("lrange", keyName, -1000000, -1).List()
		if err != nil {
			log.Fatal(err)
		}

		var bck TmRedisBackup
		bck.KeyName = keyName
		for _, val := range list {
			line := strings.Split(val, ":")
			var stat TmRedisStat
			stat.UnixTime, err = strconv.ParseInt(line[0], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			stat.Value, err = strconv.ParseInt(line[1], 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			bck.Stats = append(bck.Stats, stat)
		}
		backupList = append(backupList, bck)
	}
	jBytes, err := json.MarshalIndent(backupList, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(string(jBytes))
	outputFile, err := os.Create(*fileName)
	if err != nil {
		log.Fatal(err)
	}
	_, err = outputFile.Write(jBytes)
	if err != nil {
		log.Fatal(err)
	}
	outputFile.Close()
}
