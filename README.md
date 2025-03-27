# notion2neo4j

## Neo4j数据库模型定义

* 节点

    - **根节点（KnowledgeBase）**：将所有 tag 连接在一起
    - **标签节点（Tag）**：包含标签的名称
    - **文章节点（Article）**：包含文章的标题、内容、创建时间、更新时间、作者、标签、引用的文章

* 关系

    - **包含关系（CONTAINS）**：将标签节点连接到根节点
    - **标签关系（TAGGED_WITH）**：将文章节点连接到标签节点
    - **引用关系（REFERENCES）**：将文章节点连接到引用的文章节点


```js
(:KnowledgeBase) -[:CONTAINS]-> (:Tag)
(:Article) -[:TAGGED_WITH]-> (:Tag)
(:Article) -[:REFERENCES]-> (:Article)
```

