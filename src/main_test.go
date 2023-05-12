package main

import (
	"irfanabduhu/bookstore/utils"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	utils.Init()
	code := m.Run()
	utils.TearDown()
    os.Exit(code)
}