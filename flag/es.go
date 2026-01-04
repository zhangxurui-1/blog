package flag

import (
	"bufio"
	"fmt"
	"os"
	"server/model/elasticsearch"
	"server/service"
)

// Elasticsearch 创建 ES 索引
func Elasticsearch() error {
	esService := service.ServiceGroupApp.EsService
	// 检查索引是否已存在
	indexExists, err := esService.IndexExists("flag")
	if err != nil {
		return err
	}
	// 如果索引存在，则打印提示信息并询问是否重建索引
	if indexExists {
		// 打印提示信息
		fmt.Println("The index already exists. Do you want to delete the data and recreate the index? (y/n)")

		// 读取用户输入
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "y":
			// 删除索引
			fmt.Println("Proceeding to delete the data and recreate the index...")
			if err := esService.IndexDelete(elasticsearch.ArticleIndex()); err != nil {
				return err
			}
		case "n":
			fmt.Println("Exiting the program...")
			os.Exit(0)
		default:
			// 如果用户输入无效，提示重新输入
			fmt.Println("Invalid input. Please enter 'y' to delete and recreate the index, or 'n' to exit.")
			// 递归调用，重新输入
			return Elasticsearch()
		}
	}

	// 创建索引
	return esService.IndexCreate(elasticsearch.ArticleIndex(), elasticsearch.ArticleMapping())
}
