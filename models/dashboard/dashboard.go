package dashboard

import (
	"gin-blog/models/blog"
	"gin-blog/models/resource"
	"gin-blog/models/user"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// DashboardStats 面板统计数据结构体
type DashboardStats struct {
	ArticleCount    int     `json:"article_count"`
	PhotoCount      int     `json:"photo_count"`
	MusicCount      int     `json:"music_count"`
	UserCount       int     `json:"user_count"`
	CpuUsage        float64 `json:"cpu_usage"`
	MemoryUsage     float64 `json:"memory_usage"`
	ProcessCount    int     `json:"process_count"`
	ServiceMemoryMB float64 `json:"service_memory_mb"`
}

type ProcessStat struct {
	PID       string `json:"pid"`
	Name      string `json:"name"`
	MemoryMB  string `json:"memory_mb"`
	MemoryPct string `json:"memory_pct"`
}

// GetDashboardStats 获取面板统计数据
func GetDashboardStats() DashboardStats {
	stats := DashboardStats{}

	// 获取文章数量 (is_delete = 0)
	dbBlog := blog.GetDB()
	var articleCount int
	dbBlog.Table("t_go_article").Where("is_delete = ?", 0).Count(&articleCount)
	stats.ArticleCount = articleCount

	// 获取图片数量 (is_delete = 1)
	dbResource := resource.GetDB()
	var photoCount int
	dbResource.Table("t_go_photos").Where("is_delete = ?", 1).Count(&photoCount)
	stats.PhotoCount = photoCount

	// 获取音乐数量 (is_delete = 1)
	var musicCount int
	dbResource.Table("t_go_music").Where("is_delete = ?", 1).Count(&musicCount)
	stats.MusicCount = musicCount

	// 获取用户数量
	dbUser := user.GetDB()
	var userCount int
	dbUser.Table("ci_admin_user").Count(&userCount)
	stats.UserCount = userCount
	stats.CpuUsage = getCPUUsage()
	stats.MemoryUsage = getMemoryUsage()
	stats.ProcessCount = getProcessCount()
	stats.ServiceMemoryMB = getServiceMemoryMB()

	return stats
}

func GetTopProcessStats() []ProcessStat {
	return getTopProcesses()
}

func getTopProcesses() []ProcessStat {
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("sh", "-c", "ps -axo pid=,comm=,rss=,%mem= -r | head -n 10")
	} else {
		cmd = exec.Command("sh", "-c", "ps -eo pid=,comm=,rss=,%mem= --sort=-rss | head -n 10")
	}

	out, err := cmd.Output()
	if err != nil {
		return []ProcessStat{}
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	processes := make([]ProcessStat, 0, len(lines))
	for _, line := range lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) < 4 {
			continue
		}
		rssKB, err := strconv.ParseFloat(fields[len(fields)-2], 64)
		if err != nil {
			rssKB = 0
		}
		processes = append(processes, ProcessStat{
			PID:       fields[0],
			Name:      strings.Join(fields[1:len(fields)-2], " "),
			MemoryMB:  strconv.FormatFloat(roundFloat(rssKB/1024, 2), 'f', 2, 64),
			MemoryPct: fields[len(fields)-1],
		})
	}

	return processes
}

func getCPUUsage() float64 {
	if runtime.GOOS == "darwin" {
		out, err := exec.Command("top", "-l", "1", "-n", "0").Output()
		if err != nil {
			return 0
		}
		re := regexp.MustCompile(`([0-9.]+)% idle`)
		matches := re.FindStringSubmatch(string(out))
		if len(matches) < 2 {
			return 0
		}
		idle, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0
		}
		return roundFloat(100-idle, 2)
	}

	if runtime.GOOS == "linux" {
		out, err := exec.Command("sh", "-c", "top -bn1 | grep 'Cpu(s)'").Output()
		if err != nil {
			return 0
		}
		re := regexp.MustCompile(`([0-9.]+)\s*id`)
		matches := re.FindStringSubmatch(string(out))
		if len(matches) < 2 {
			return 0
		}
		idle, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0
		}
		return roundFloat(100-idle, 2)
	}

	return 0
}

func getMemoryUsage() float64 {
	if runtime.GOOS == "darwin" {
		vmOut, err := exec.Command("vm_stat").Output()
		if err != nil {
			return 0
		}
		memOut, err := exec.Command("sysctl", "-n", "hw.memsize").Output()
		if err != nil {
			return 0
		}

		totalBytes, err := strconv.ParseFloat(strings.TrimSpace(string(memOut)), 64)
		if err != nil || totalBytes <= 0 {
			return 0
		}

		pageSize := 4096.0
		pageRe := regexp.MustCompile(`page size of ([0-9]+) bytes`)
		if matches := pageRe.FindStringSubmatch(string(vmOut)); len(matches) > 1 {
			if parsed, parseErr := strconv.ParseFloat(matches[1], 64); parseErr == nil {
				pageSize = parsed
			}
		}

		usedPages := extractVMStatPages(string(vmOut), "Pages active") +
			extractVMStatPages(string(vmOut), "Pages wired down") +
			extractVMStatPages(string(vmOut), "Pages occupied by compressor")
		if usedPages <= 0 {
			return 0
		}

		usedBytes := usedPages * pageSize
		return roundFloat((usedBytes/totalBytes)*100, 2)
	}

	if runtime.GOOS == "linux" {
		out, err := exec.Command("sh", "-c", "grep -E 'MemTotal|MemAvailable' /proc/meminfo").Output()
		if err != nil {
			return 0
		}
		total := extractMemInfoValue(string(out), "MemTotal")
		available := extractMemInfoValue(string(out), "MemAvailable")
		if total <= 0 {
			return 0
		}
		return roundFloat((float64(total-available)/float64(total))*100, 2)
	}

	return 0
}

func getProcessCount() int {
	out, err := exec.Command("sh", "-c", "ps -A -o pid=").Output()
	if err != nil {
		return 0
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	count := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			count++
		}
	}
	return count
}

func getServiceMemoryMB() float64 {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return roundFloat(float64(mem.Alloc)/1024/1024, 2)
}

func extractVMStatPages(input string, label string) float64 {
	re := regexp.MustCompile(label + `:\s+([0-9.]+)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return 0
	}
	value := strings.ReplaceAll(matches[1], ".", "")
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return parsed
}

func extractMemInfoValue(input string, key string) int64 {
	re := regexp.MustCompile(key + `:\s+([0-9]+)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return 0
	}
	parsed, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0
	}
	return parsed
}

func roundFloat(value float64, precision int) float64 {
	ratio := mathPow10(precision)
	return float64(int(value*ratio+0.5)) / ratio
}

func mathPow10(n int) float64 {
	result := 1.0
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}
