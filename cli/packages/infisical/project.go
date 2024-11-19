package infisical

import (
	"errors"
	"fmt"
	"strings"
)


func GetProjectIdByName(n string, a string) (string, error) {
  nUpper := strings.ToUpper(n)
  switch strings.ToUpper(a){
  case "DOK", "DOKTERONLINE":
    return getProjectID(strings.Replace(fmt.Sprintf("DOK_%v_PROJECT_ID", nUpper),"-", "_", -1)), nil
  case "SMNP", "SEEMENOPAUSE":
    return getProjectID(strings.Replace(fmt.Sprintf("SMNP_%v_PROJECT_ID", nUpper),"-", "_", -1)), nil
  case "EHVG", "EHEALTHVENTURESGROUP":
    return getProjectID(strings.Replace(fmt.Sprintf("SMNP_%v_PROJECT_ID", nUpper),"-", "_", -1)), nil
  default:
    return "", errors.New("no valid application found")
  }
}