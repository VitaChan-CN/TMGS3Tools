# TMGS3Tools

《ときめきメモリアル Girls Side 3rd Story》相关工具  

目前仅支持DFI(idx img)文件解包和打包，支持OFS3解包，暂不支持OFS3打包。  
OFS3工具可以使用 [https://github.com/Liquid-S/OFS3-TOOL](https://github.com/Liquid-S/OFS3-TOOL)  

具体使用方法参见**Example**  

## Usage
```shell
  -append
        [打包]追加写入模式，待测试
  -i string
        [打包.必要]输入文件夹路径；[OFS3解包]输入OFS3文件名
  -idx string
        [必要]cdimg.idx文件名
  -img string
        [必要]cdimg.idx文件名
  -log
        显示日志 (default true)
  -o string
        [解包]输出文件夹路径；[打包]输出文件名
  -ofs3
        [解包]递归解包所有OFS3格式文件，无法打包
  -ofs3.log
        显示OFS3日志

```

## Example
```shell
# OFS3单独解包
TMGS3Tools -i=data/ofs3/005.bin \
           -o=data/ofs3/output \
           -ofs3 -ofs3log

# 解包，输出到data/01/output文件夹中
TMGS3Tools -idx=data/01/a.idx \
           -img=data/01/a.img \
           -o=data/01/output 

# 解包（含OFS3文件），输出到data/01/output文件夹中
TMGS3Tools -idx=data/01/a.idx \
           -img=data/01/a.img \
           -o=data/01/output \
           -ofs3 # -ofs3log

# 打包 普通模式，输出文件名为data/01/a.out.idx和data/01/a.out.img
TMGS3Tools -idx=data/01/a.idx \
           -img=data/01/a.img \
           -i=data/01/output \
           -o=data/01/a.out  

# 打包 追加模式【未测试】，输出文件名为data/01/a.out.idx和data/01/a.out.img
TMGS3Tools -idx=data/01/a.idx \
           -img=data/01/a.img \
           -i=data/01/output \
           -o=data/01/a.out \
           -append 

```
## 未来
- 支持OFS3打包

## 更新日志

### 2022-4-5 3
- 支持OFS3格式文件递归解包

### 2022-4-5 2
- 任意DFI(idx img)文件打包
- 支持追加模式打包

### 2022-4-5 1
- 任意DFI(idx img)文件提取

### 2022-3-26
- DFI(idx img)文件提取，仅GS3 cdimg0文件
