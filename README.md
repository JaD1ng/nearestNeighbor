# 最近邻搜索算法

## 项目概述
本项目实现了多种最近邻搜索算法，并通过比较它们在处理相同数据集时的性能来评估它们的效率。这些算法包括KD树、R树、局部敏感哈希（LSH）、并行搜索以及暴力搜索方法。

## 结构
- `algorithms/`: 包含所有最近邻搜索算法的实现。
    - `kdtree`: KD树算法实现。
    - `rtree`: R树算法实现。
    - `lsh`: 局部敏感哈希算法实现。
    - `parallel`: 并行搜索算法实现。
    - `violent`: 暴力搜索算法实现。
- `common/`: 包含项目中使用的公共配置和工具。
- `config/`: 包含配置文件。
- `tests/`: 包含测试数据和测试脚本。
- `tool/`: 包含数据加载和处理的工具函数。

## 待完成
- rtree使用github.com/dhconnelly/rtreego库，性能较差；自己实现的rtree结果不正确，需要进一步调试。
- 尝试希尔伯特曲线
- ...