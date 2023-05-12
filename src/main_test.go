package main

import (
	"irfanabduhu/bookstore/utils"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	utils.InitDB()
	code := m.Run()
	utils.TearDown()
    os.Exit(code)
}