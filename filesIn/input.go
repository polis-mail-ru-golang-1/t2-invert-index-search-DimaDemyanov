package filesIn

import (
	"fmt"
  "io/ioutil"
)

func ReadData(filename string)(string, error){
  inBytes, err := ioutil.ReadFile(filename)
  if err != nil {
    fmt.Println("Error occured while reading file:")
    fmt.Println(err)
    return "", err
  } else {
    str := string(inBytes)
    return str, err
  }
  //return nil, nil
}
