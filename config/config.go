// All the configuration is set here only
package config

import "github.com/vviveksharma/auth/db"


func Init() {
	db.ConnectDB()
}