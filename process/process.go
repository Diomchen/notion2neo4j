package process

import (
	"context"
	"encoding/csv"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// 读取 csv 文件
func ProcessCSV(filePath string) ([]DataEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var entries []DataEntry
	for index, record := range records {
		if len(record) != 3 || index == 0 {
			continue
		}

		var id = index

		tags := strings.Split(strings.TrimSpace(record[1]), ",")

		timeStr := strings.TrimSpace(record[2])
		t, err := time.Parse("2006年1月2日 15:04", timeStr)
		if err != nil {
			t = time.Now()
		}

		entries = append(entries, DataEntry{
			ID:    id,
			Title: strings.TrimSpace(record[0]),
			Tags:  tags,
			Time:  t,
		})

	}

	return entries, nil
}

// 解析文章
func ProcessArticle(dataFolder string, entries []DataEntry) (map[int]Article, error) {
	articles := make(map[int]Article)
	titleToID := make(map[string]int)

	for _, entry := range entries {
		titleToID[entry.Title] = entry.ID
	}

	files, err := ioutil.ReadDir(dataFolder)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		fileName := file.Name()
		var fileTitle string
		// 找到最后一个空格
		lastSpaceIndex := strings.LastIndex(fileName, " ")
		if lastSpaceIndex == -1 {
			fileTitle = fileName
		} else {
			fileTitle = strings.TrimSpace(fileName[:lastSpaceIndex])
		}

		id, ok := titleToID[fileTitle]
		if !ok {
			continue
		}

		content, err := os.ReadFile(filepath.Join(dataFolder, fileName))
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			if entry.ID == id {
				articles[id] = Article{
					ID:      id,
					Title:   fileTitle,
					Content: string(content),
					Time:    entry.Time,
				}
				break
			}
		}
	}
	return articles, nil
}

func CreateNeo4jGraph(entries []DataEntry, articles map[int]Article) error {

	ctx := context.Background()
	defer NEO4J_DIVER.Close(ctx)

	session := NEO4J_DIVER.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	// 添加根节点
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		res, err := tx.Run(
			ctx,
			"CREATE (kb:KnowledgeBase {name: 'Knowledge Root'}) RETURN kb",
			map[string]interface{}{},
		)
		if err != nil {
			return nil, err
		}
		return res.Consume(ctx)
	})
	if err != nil {
		return err
	}

	// 添加 tag 节点
	allTags := make(map[string]bool)
	for _, entry := range entries {
		for _, tag := range entry.Tags {
			allTags[tag] = true
		}
	}

	for tag := range allTags {
		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
			result, err := tx.Run(
				ctx,
				"MATCH (kb:KnowledgeBase {name: 'Knowledge Root'}) "+
					"MERGE (t:Tag {name: $tagName}) "+
					"MERGE (kb)-[:CONTAINS]->(t) "+
					"RETURN t",
				map[string]interface{}{"tagName": tag},
			)
			if err != nil {
				return nil, err
			}
			return result.Consume(ctx)
		})
		if err != nil {
			return err
		}
	}

	// 添加数据节点与 tag 之间的关系
	for _, entry := range entries {
		article, ok := articles[entry.ID]
		if !ok {
			continue
		}

		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
			res, err := tx.Run(
				ctx,
				"CREATE (a:Article {id: $id, title: $title, content: $content, createdAt: $createdAt}) "+"RETURN a",
				map[string]interface{}{
					"id":        article.ID,
					"title":     article.Title,
					"content":   article.Content,
					"createdAt": article.Time.Format(time.RFC3339),
				},
			)
			if err != nil {
				return nil, err
			}
			return res.Consume(ctx)
		})
		if err != nil {
			return err
		}

		for _, tag := range entry.Tags {
			_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
				res, err := tx.Run(
					ctx,
					"MATCH (a:Article {id: $articleID}) "+
						"MATCH (t:Tag {name: $tagName}) "+
						"MERGE (a)-[:TAGGED_WITH]->(t) "+
						"RETURN a, t",
					map[string]interface{}{
						"articleID": article.ID,
						"tagName":   tag,
					},
				)
				if err != nil {
					return nil, err
				}

				return res.Consume(ctx)
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
