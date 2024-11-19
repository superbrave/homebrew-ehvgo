package infisical

import (
	"strings"
)

func getFullApplicationName(n string) string {
  switch strings.ToLower(n) {
  case "dok", "dokteronline": 
    return "Dokteronline"
  case "smnp", "seemenopause":
    return "SeeMeNoPause"
  case "ehvg", "ehealthventuresgroup":
    return "eHealth Ventures Group"
  default:
    return ""
  }
}

func CompletionListSites() []string {
  return []string{"dokteronline", "dok", "seemenopause", "smnp", "ehvg", "ehealthventuresgroup"}
}
