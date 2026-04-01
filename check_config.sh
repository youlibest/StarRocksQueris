#!/bin/bash

# 配置文件检查脚本

echo "=== StarRocksQueris 配置文件检查 ==="
echo

# 获取可执行文件路径
EXEC_PATH=$(pwd)
EXEC_NAME="./StarRocksQueris"
CONFIG_FILE=".StarRocksQueris.yaml"

echo "[1] 检查当前目录: $EXEC_PATH"
echo "[1] 预期配置文件: $EXEC_PATH/$CONFIG_FILE"
echo

# 检查文件是否存在
if [ -f "$CONFIG_FILE" ]; then
    echo "[✓] 配置文件存在"
    echo "[i] 文件大小: $(ls -lh $CONFIG_FILE | awk '{print $5}')"
    echo "[i] 文件权限: $(ls -l $CONFIG_FILE | awk '{print $1}')"
    echo
    
    echo "[2] 文件内容前10行:"
    head -10 "$CONFIG_FILE"
    echo
    
    echo "[3] 检查隐藏字符 (显示换行符和制表符):"
    cat -A "$CONFIG_FILE" | head -10
    echo
    
    echo "[4] 文件编码检查:"
    file "$CONFIG_FILE"
    echo
    
    echo "[5] 检查是否有 BOM 头:"
    if head -c 3 "$CONFIG_FILE" | od -An -tx1 | grep -q "efbbbf"; then
        echo "[✗] 警告: 文件包含 UTF-8 BOM 头，需要移除"
        echo "[i] 修复方法: sed -i '1s/^\xEF\xBB\xBF//' $CONFIG_FILE"
    else
        echo "[✓] 文件没有 BOM 头"
    fi
    echo
    
    echo "[6] YAML 语法检查 (需要安装 yq 或 python-yaml):"
    if command -v yq &> /dev/null; then
        yq eval '.' "$CONFIG_FILE" > /dev/null 2>&1 && echo "[✓] YAML 语法正确" || echo "[✗] YAML 语法错误"
    elif python3 -c "import yaml" 2>/dev/null; then
        python3 -c "import yaml; yaml.safe_load(open('$CONFIG_FILE'))" 2>&1 && echo "[✓] YAML 语法正确" || echo "[✗] YAML 语法错误"
    else
        echo "[!] 未安装 yq 或 python-yaml，跳过语法检查"
    fi
else
    echo "[✗] 配置文件不存在: $CONFIG_FILE"
    echo "[i] 请确保配置文件在当前目录下，且文件名为: $CONFIG_FILE"
    echo
    echo "当前目录文件列表:"
    ls -la
fi

echo
echo "=== 检查完成 ==="
