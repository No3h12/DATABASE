package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// اتصال بقاعدة البيانات الرئيسية (الماستر)
	masterDB, err := sql.Open("mysql", "root:root@tcp(192.168.1.249:3306)/MOST")
	if err != nil {
		log.Fatal("❌ Connection to master database failed:", err)
	}
	defer masterDB.Close()

	// اتصال بقاعدة البيانات المحلية (النسخة الخاصة بك)
	localDB, err := sql.Open("mysql", "root:rootroot@tcp(localhost:3306)/MOST")
	if err != nil {
		log.Fatal("❌ Connection to local database failed:", err)
	}
	defer localDB.Close()

	reader := bufio.NewReader(os.Stdin)

	// الحلقة الرئيسية لاختيار الجدول
mainLoop:
	for {
		// عرض الجداول المتوفرة (نستخدم الماستر كمرجع)
		tables, err := getTables(masterDB)
		if err != nil {
			log.Fatal("❌ Failed to fetch tables:", err)
		}

		if len(tables) == 0 {
			fmt.Println("❌ No tables found.")
			return
		}

		fmt.Println("\n📦 Available Tables:")
		for i, table := range tables {
			fmt.Printf("%d - %s\n", i+1, table)
		}

		// اختيار الجدول
		var tableName string
		for {
			fmt.Print("\n🟢 Enter the number of the table to work with: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			index, err := strconv.Atoi(input)
			if err != nil || index < 1 || index > len(tables) {
				fmt.Println("⚠️ Invalid choice, please try again.")
				continue
			}
			tableName = tables[index-1]
			break
		}

		fmt.Printf("\n✅ You selected table: %s\n", tableName)

		// عرض الأعمدة (نستخدم الماستر كمرجع)
		columns, err := getColumns(masterDB, tableName)
		if err != nil {
			log.Fatal("❌ Failed to get columns:", err)
		}

		// قائمة الخيارات
		for {
			fmt.Println("\n📌 Choice Process:")
			fmt.Println("1 - INSERT")
			fmt.Println("2 - UPDATE")
			fmt.Println("3 - DELETE")
			fmt.Println("4 - SELECT ALL")
			fmt.Println("5 - SEARCH")
			fmt.Println("6 - BACK")
			fmt.Println("7 - Exit")
			fmt.Print("\n  Enter The Number Of The process: ")

			choiceStr, _ := reader.ReadString('\n')
			choiceStr = strings.TrimSpace(choiceStr)

			switch choiceStr {
			case "1":
				// INSERT
				values := []interface{}{}
				placeholders := []string{}
				colNames := []string{}

				for _, col := range columns {
					if col == "id" {
						continue // نتجنب id إذا كان auto_increment
					}
					fmt.Printf("Enter %s: ", col)
					val, _ := reader.ReadString('\n')
					values = append(values, strings.TrimSpace(val))
					placeholders = append(placeholders, "?")
					colNames = append(colNames, col)
				}

				query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(colNames, ","), strings.Join(placeholders, ","))
				
				// تنفيذ على الماستر
				_, err := masterDB.Exec(query, values...)
				if err != nil {
					fmt.Println("\n❌ Error during INSERT on master:", err)
				} else {
					fmt.Println("\n✅ Insert completed on master!")
				}

				// تنفيذ على النسخة المحلية
				_, err = localDB.Exec(query, values...)
				if err != nil {
					fmt.Println("\n❌ Error during INSERT on local:", err)
				} else {
					fmt.Println("\n✅ Insert completed on local!")
				}

			case "2":
				// UPDATE
				fmt.Print("\nEnter ID to update: ")
				idStr, _ := reader.ReadString('\n')
				idStr = strings.TrimSpace(idStr)

				updates := []string{}
				values := []interface{}{}

				for _, col := range columns {
					if col == "id" {
						continue
					}
					fmt.Printf("New value for %s (leave empty to skip): ", col)
					val, _ := reader.ReadString('\n')
					val = strings.TrimSpace(val)
					if val != "" {
						updates = append(updates, fmt.Sprintf("%s = ?", col))
						values = append(values, val)
					}
				}

				if len(updates) == 0 {
					fmt.Println("\n⚠️ No fields to update.")
					continue
				}

				values = append(values, idStr)
				query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", tableName, strings.Join(updates, ", "))
				
				// تنفيذ على الماستر
				_, err := masterDB.Exec(query, values...)
				if err != nil {
					fmt.Println("\n❌ Error during UPDATE on master:", err)
				} else {
					fmt.Println("\n✅ Update completed on master!")
				}

				// تنفيذ على النسخة المحلية
				_, err = localDB.Exec(query, values...)
				if err != nil {
					fmt.Println("\n❌ Error during UPDATE on local:", err)
				} else {
					fmt.Println("\n✅ Update completed on local!")
				}

			case "3":
				// DELETE
				fmt.Print("\nEnter ID to delete: ")
				idStr, _ := reader.ReadString('\n')
				idStr = strings.TrimSpace(idStr)

				query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)
				
				// تنفيذ على الماستر
				_, err := masterDB.Exec(query, idStr)
				if err != nil {
					fmt.Println("\n❌ Error during DELETE on master:", err)
				} else {
					fmt.Println("\n✅ Record deleted from master!")
				}

				// تنفيذ على النسخة المحلية
				_, err = localDB.Exec(query, idStr)
				if err != nil {
					fmt.Println("\n❌ Error during DELETE on local:", err)
				} else {
					fmt.Println("\n✅ Record deleted from local!")
				}

			case "4":
				// SELECT ALL - عرض البيانات من كلا القاعدتين
				query := fmt.Sprintf("SELECT * FROM %s", tableName)
				
				// عرض بيانات الماستر
				fmt.Println("\n📋 Records from MASTER database:")
				rowsMaster, err := masterDB.Query(query)
				if err != nil {
					fmt.Println("\n❌ Error during SELECT from master:", err)
				} else {
					printResults(rowsMaster)
					rowsMaster.Close()
				}

				// عرض بيانات النسخة المحلية
				fmt.Println("\n📋 Records from LOCAL database:")
				rowsLocal, err := localDB.Query(query)
				if err != nil {
					fmt.Println("\n❌ Error during SELECT from local:", err)
				} else {
					printResults(rowsLocal)
					rowsLocal.Close()
				}

			case "5":
				// SEARCH
				fmt.Print("\nEnter keyword to search: ")
				keyword, _ := reader.ReadString('\n')
				keyword = strings.TrimSpace(keyword)

				searchCol := columns[1] // نفترض البحث في العمود الثاني
				query := fmt.Sprintf("SELECT * FROM %s WHERE %s LIKE ?", tableName, searchCol)
				
				// بحث في الماستر
				fmt.Println("\n🔍 Search Results from MASTER database:")
				rowsMaster, err := masterDB.Query(query, "%"+keyword+"%")
				if err != nil {
					fmt.Println("\n❌ Error during SEARCH in master:", err)
				} else {
					printResults(rowsMaster)
					rowsMaster.Close()
				}

				// بحث في النسخة المحلية
				fmt.Println("\n🔍 Search Results from LOCAL database:")
				rowsLocal, err := localDB.Query(query, "%"+keyword+"%")
				if err != nil {
					fmt.Println("\n❌ Error during SEARCH in local:", err)
				} else {
					printResults(rowsLocal)
					rowsLocal.Close()
				}

			case "6":
				fmt.Println("\n🔙 Going back to table selection...")
				continue mainLoop

			case "7":
				fmt.Println("\n👋 Exit program.")
				return

			default:
				fmt.Println("\n⚠️ Invalid choice.")
			}
		}
	}
}

// دالة مساعدة لعرض النتائج
func printResults(rows *sql.Rows) {
	colNames, _ := rows.Columns()
	values := make([]sql.RawBytes, len(colNames))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			fmt.Println("❌ Error:", err)
			continue
		}
		for i, val := range values {
			fmt.Printf("%s: %s  ", colNames[i], string(val))
		}
		fmt.Println()
	}
}

func getTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	var tableName string
	for rows.Next() {
		err := rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}
	return tables, nil
}

func getColumns(db *sql.DB, tableName string) ([]string, error) {
	rows, err := db.Query(fmt.Sprintf("SHOW COLUMNS FROM %s", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var field, colType, null, key, defVal, extra sql.NullString
		err := rows.Scan(&field, &colType, &null, &key, &defVal, &extra)
		if err != nil {
			return nil, err
		}
		columns = append(columns, field.String)
	}
	return columns, nil
}