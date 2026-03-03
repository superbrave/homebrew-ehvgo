package aws

import (
    "bufio"
    "errors"
    "os"
    "path/filepath"
    "sort"
    "strings"

    "github.com/manifoldco/promptui"
)

// ReadProfiles reads AWS profiles from ~/.aws/config.
func ReadProfiles() ([]string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }

    configPath := filepath.Join(home, ".aws", "config")
    file, err := os.Open(configPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    profiles := make(map[string]struct{})
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if !strings.HasPrefix(line, "[") || !strings.HasSuffix(line, "]") {
            continue
        }
        section := strings.TrimSpace(line[1 : len(line)-1])
        if section == "default" {
            profiles["default"] = struct{}{}
            continue
        }
        if strings.HasPrefix(section, "profile ") {
            name := strings.TrimSpace(strings.TrimPrefix(section, "profile "))
            if name != "" {
                profiles[name] = struct{}{}
            }
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    if len(profiles) == 0 {
        return nil, errors.New("no profiles found in ~/.aws/config")
    }

    list := make([]string, 0, len(profiles))
    for name := range profiles {
        list = append(list, name)
    }
    sort.Strings(list)

    return list, nil
}

// PromptForProfile prompts for a profile selection.
func PromptForProfile(profiles []string) (string, error) {
    selectPrompt := promptui.Select{
        Label:  "Select AWS profile",
        Items:  profiles,
        Size:   10,
        Stdout: bellSkipper{},
        Templates: &promptui.SelectTemplates{
            Active:   "  \x1b[4m{{ . }}\x1b[0m",
            Inactive: "  {{ . }}",
            Selected: "  {{ . }}",
        },
    }

    _, result, err := selectPrompt.Run()
    if err != nil {
        return "", err
    }

    return result, nil
}
