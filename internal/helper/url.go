package helper

import "strings"

func NormalizeURL(raw string) string {
    if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
        return "https://" + raw
    }
    return raw
}
