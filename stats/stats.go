package stats

import (
    "fmt"
    "os"
    "encoding/json"
    "time"
)

type DayStats struct {
    WorkMs        float32
    BreakMs       float32
}

func (s *DayStats) IsEmpty() bool {
    return s.BreakMs == 0 && s.WorkMs == 0
}

type Stats map[string]DayStats

func LoadStats(path string) Stats {
    fmt.Printf("Stats :: Loading stats file: \"%s\"\n", path)
    stats := Stats{}

    fileContent, err := os.ReadFile(path)
    if err != nil {
        fmt.Println("Stats :: Failed to load stats file:", err)
        return Stats{};
    } else {
        err = json.Unmarshal(fileContent, &stats)
        if err != nil {
            fmt.Println("Stats :: Failed to parse JSON:", err)
            return Stats{};
        } else {
            fmt.Printf("Stats :: Loaded stats: %+v\n", stats)
        }
    }
    return stats
}

func WriteStats(path string, stats *Stats) {
    fmt.Printf("Stats :: Writing stats file: \"%s\"\n", path)

    jsonData, err := json.MarshalIndent(&stats, "", "    ")
    if err != nil {
        fmt.Println("Stats :: Failed to build JSON:", err)
    } else {
        err = os.WriteFile(path, jsonData, 0o644)
        if err != nil {
            fmt.Println("Stats :: Failed to write stats file:", err)
        } else {
            fmt.Printf("Stats :: Wrote stats: %+v\n", *stats)
        }
    }
}

func (s *Stats) GetDay(date *string) DayStats {
    val, exists := (*s)[*date]
    if exists {
        return val
    } else {
        return DayStats{}
    }
}

func GetCurrentDate() string {
    return time.Now().Format("2006-01-02")
}
