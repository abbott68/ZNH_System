package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// 数据结构，用于展示报表
type Report struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// 数据库连接信息
const (
	DBUsername = "root"
	DBPassword = "Root@123"
	DBHost     = "81.70.196.2"
	DBPort     = "3306"
	DBName     = "test_db"
)

// 数据库连接
var db *sql.DB

type Host struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IP       string `json:"ip"`
	Location string `json:"location"`
	Status   string `json:"status"`
	ScanTime string `json:"scan_time"`
}

//// 构造主机列表
//var hosts []Host

func init() {
	// 初始化主机信息
	//hosts = []Host{
	//	{ID: 1, Name: "Host 1", IP: "192.168.0.1", Location: "Data Center A"},
	//	{ID: 2, Name: "Host 2", IP: "192.168.0.2", Location: "Data Center B"},
	//	{ID: 3, Name: "Host 3", IP: "192.168.0.3", Location: "Data Center C"},
	//}
	// 连接数据库
	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBUsername, DBPassword, DBHost, DBPort, DBName)
	var err error
	db, err = sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// 确保数据库连接正常
	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// 创建主机表（如果不存在）
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS hosts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		ip VARCHAR(255) NOT NULL,
		location VARCHAR(255) NOT NULL
	)`)
	if err != nil {
		log.Fatal("Error creating hosts table:", err)
	}
}

// 处理主机列表请求
func hostListHandler(w http.ResponseWriter, r *http.Request) {

	// 从数据库中查询主机信息
	rows, err := db.Query("SELECT id, name, ip, status, scan_time FROM hosts")
	if err != nil {
		http.Error(w, "Error retrieving hosts", http.StatusInternalServerError)
		log.Println("Error retrieving hosts:", err)
		return
	}
	defer rows.Close()
	// 构造主机列表
	var hostList []Host
	for rows.Next() {
		var host Host
		err := rows.Scan(&host.ID, &host.Name, &host.IP, &host.Status, &host.ScanTime)
		if err != nil {
			http.Error(w, "Error retrieving hosts", http.StatusInternalServerError)
			log.Println("Error retrieving hosts:", err)
			return
		}
		hostList = append(hostList, host)
	}
	// 输出主机列表
	for _, host := range hostList {
		fmt.Fprintf(w, "ID: %d\n", host.ID)
		fmt.Fprintf(w, "Name: %s\n", host.Name)
		fmt.Fprintf(w, "IP: %s\n", host.IP)
		fmt.Fprintf(w, "Status: %s\n", host.Status)
		fmt.Fprintf(w, "Scan Time: %s\n", host.ScanTime)
		fmt.Fprintln(w, "-------------------------")
	}

	defer rows.Close()
	// 构造主机列表
	var hosts []Host
	for rows.Next() {
		var host Host
		err := rows.Scan(&host.ID, &host.Name, &host.IP, &host.Location)
		if err != nil {
			log.Println("Error scanning host:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		hosts = append(hosts, host)
	}
	if err = rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// 将主机信息转换为JSON格式
	jsonData, err := json.Marshal(hosts)
	if err != nil {
		log.Println("Error marshaling host data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 发送响应数据
	_, err = w.Write(jsonData)
	if err != nil {
		log.Println("Error writing response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 刷新响应，确保完整发送给客户端
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

// 处理删除主机请求

// 处理添加主机请求
func addHostHandler(w http.ResponseWriter, r *http.Request) {
	// 解析请求参数
	name := r.FormValue("name")
	ip := r.FormValue("ip")
	location := r.FormValue("location")
	// 插入新主机记录到数据库
	result, err := db.Exec("INSERT INTO hosts (name, ip, location) VALUES (?, ?, ?)", name, ip, location)
	if err != nil {
		log.Println("Error inserting host:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// 获取插入的主机ID
	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// 构造响应数据
	response := map[string]interface{}{
		"id":       id,
		"name":     name,
		"ip":       ip,
		"location": location,
	}
	// 将响应数据转换为JSON格式
	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshaling response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// 发送响应数据
	_, err = w.Write(jsonData)
	if err != nil {
		log.Println("Error writing response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// 处理报表请求
func reportHandler(w http.ResponseWriter, r *http.Request) {
	// 获取报表数据
	report := fetchReportData()

	// 将报表数据转换为JSON格式
	jsonData, err := json.Marshal(report)
	if err != nil {
		log.Println("Error marshaling report data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 发送报表数据作为响应
	_, err = w.Write(jsonData)
	if err != nil {
		log.Println("Error writing response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

//获取报表数据（示例函数）
func fetchReportData() *Report {
	// 从数据库或其他数据源获取报表数据
	// 这里只是一个简单的示例，返回固定的报表数据
	return &Report{
		Title:   "运维报表",
		Content: "这是一个示例报表",
	}
}

func main() {
	// 设置路由
	http.HandleFunc("/report", reportHandler)
	http.HandleFunc("/hosts", hostListHandler)
	//http.HandleFunc("/deletehost", deleteHostHandler)
	http.HandleFunc("/addhost", addHostHandler)
	http.HandleFunc("/scanhosts", scanHostsHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// 设置静态文件目录
	//fs := http.FileServer(http.Dir("static"))
	//http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 设置模板渲染
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			log.Println("Error parsing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println("Error executing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	// 启动Web服务器
	fmt.Println("Server listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}

//扫描主机
func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
func scanHostsHandler(w http.ResponseWriter, r *http.Request) {
	// 获取本地IP地址
	localIP, err := getLocalIP()
	if err != nil {
		log.Println("Error getting local IP:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 扫描网络上的主机
	hosts := scanHosts(localIP)

	// 将扫描结果转换为JSON格式
	jsonBytes, err := json.Marshal(hosts)
	if err != nil {
		log.Println("Error marshaling scan results:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 发送扫描结果作为响应
	_, err = w.Write(jsonBytes)
	if err != nil {
		log.Println("Error writing response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
func getLocalIP() (string, error) {
	// 获取主机名
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	// 解析主机名对应的IP地址
	addrs, err := net.LookupIP(hostname)
	if err != nil {
		return "", err
	}

	// 找到第一个非回环地址
	for _, addr := range addrs {
		if !addr.IsLoopback() {
			return addr.String(), nil
		}
	}

	return "", fmt.Errorf("No non-loopback address found for host")
}

func scanHosts(localIP string) []Host {
	// 定义存储主机扫描结果的切片
	var hosts []Host

	// 扫描主机
	for i := 1; i <= 255; i++ {
		ip := fmt.Sprintf("192.168.0.%d", i)

		// 执行Ping命令检测主机是否可达
		cmd := fmt.Sprintf("ping -c 1 -W 1 %s", ip)
		output, err := execCommand(cmd)
		if err != nil {
			log.Println("Error executing command:", err)
			continue
		}

		// 检查Ping命令输出，判断主机状态
		status := "Offline"
		if strings.Contains(output, "1 packets transmitted, 1 received") {
			status = "Online"
		}

		// 构造主机结构体并添加到切片中
		// 构造主机结构体并添加到切片中
		host := Host{
			ID:       i,
			Name:     "Host " + strconv.Itoa(i),
			IP:       ip,
			Status:   status,
			ScanTime: time.Now().Format("2006-01-02 15:04:05"),
		}
		hosts = append(hosts, host)
	}

	return hosts
}

// 将主机信息写入数据库
func writeHostsToDB(hosts []Host) error {
	// 准备插入语句
	stmt, err := db.Prepare("INSERT INTO hosts (name, ip, status, scan_time) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// 执行插入操作
	for _, host := range hosts {
		_, err := tx.Stmt(stmt).Exec(host.Name, host.IP, host.Status, host.ScanTime)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// 处理主机列表请求
//func hostListHandler(w http.ResponseWriter, r *http.Request) {
//}

// 执行系统命令
func execCommand(cmd string) (string, error) {
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:]
	cmdOut, err := exec.Command(head, parts...).Output()
	if err != nil {
		return "", err
	}
	return string(cmdOut), nil
}

func getIPPrefix(ip string) string {
	parts := strings.Split(ip, ".")
	return parts[0] + "." + parts[1] + "." + parts[2] + "."
}

// 扫描主机状态
func scanHostStatus(ip string) string {
	// 在此处实现扫描主机状态的逻辑
	// 根据实际情况判断主机的在线或离线状态，并返回相应的字符串表示
	return "Online"
}
