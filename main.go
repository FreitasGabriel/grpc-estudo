package main

import (
	"bitbucket.org/bemobidev/grpc-estudo/user"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"github.com/golang/protobuf/proto"
)

//go:generate protoc --go_out=. ./user/user.proto

var endianness = binary.LittleEndian

const (
	dbPath = "user.pb"
)

type length uint16

func add(id int64, name, email string) (err error) {
	u := &user.User{
		ID: id,
		Name: name,
		Email: email,
	}

	//Fazendo Marshal da struct User (vai passar a ser binário)
	b, err := proto.Marshal(u)
	if err != nil {
		return fmt.Errorf("could not encode task: %v", err)
	}

	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil{
		return fmt.Errorf("could not open %s: %v", dbPath, err)
	}

	if err = binary.Write(f, endianness, length(len(b))); err != nil {
		return fmt.Errorf("could not encode length of message: %v", err)
	}

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("could not write task to file: %v", err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("could not close file %s: %v", dbPath, err)
	}

	return nil
}

func list() (err error){
	f, err := os.Open(dbPath)
	if err != nil {
		return fmt.Errorf("could not open file %s: %v", dbPath, err)
	}
	defer func() {
		e := f.Close()
		if e != nil {
			fmt.Println(e)
		}
	}()

	for {
		//load record file

		var l length
		err = binary.Read(f, endianness, &l)
		if err != nil {
			if err == io.EOF {
				err = nil
				return
			}
			return fmt.Errorf("could not read file %s: %v", dbPath, err)
		}

		//load record
		bs := make([]byte, l)
		_, err = io.ReadFull(f, bs)
		if err != nil {
			return fmt.Errorf("could not read file %s: %v", dbPath, err)
		}

		//Unmarshal
		var u user.User
		err = proto.Unmarshal(bs, &u)
		if err != nil {
			return fmt.Errorf("could not read user: %v", err)
		}

		//Print
		fmt.Println("id:", u.GetID())
		fmt.Println("name:", u.GetName())
		fmt.Println("email:", u.GetEmail())
		fmt.Println("----------------------")
	}


}

func main() {
	err := add(2, "Teste2", "teste2@grpc.com.br")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = list()
	if err != nil {
		fmt.Println(err)
		return
	}
}