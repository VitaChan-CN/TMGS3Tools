# TMGS3Tools

《ときめきメモリアル Girls Side 3rd Story》相关工具  

目前仅支持DFI(idx img)文件解包和打包，支持OFS3解包与打包。  

具体使用方法参见**Example**  

## Usage
DFI即idx文件，需要与之对应的img文件  
开启日志`-log`和`-ofs3.log`会影响性能，默认**关闭**所有日志  
无论解包还是打包，均需要**对应的原文件**  
```shell
  -append
        [打包]追加写入模式
  -dfi.ofs3
        [DFI解包]递归解包所有OFS3格式文件
  -gz
        解包时是否自动解压gz文件(解压后为.dgz文件，导入需要手动压缩并去掉后缀)
  -i string
        [打包]输入文件夹路径
  -idx string
        [DFI必要]cdimg.idx文件名
  -img string
        [DFI必要]cdimg.idx文件名
  -log
        显示日志
  -o string
        [解包]输出文件夹路径, [打包]输出文件名
  -o2 string
        [打包]输出idx文件名，若为空则为-o后增加.idx
  -ofs3 string
        [OFS3必要] OFS3文件名
  -ofs3.log
        显示OFS3日志
  -patch int
        [打包]对已存在的-o文件的指定位置进行修改而不是创建新的，仅append模式有效。输入原img大小

```

## Example
```shell
# 打包 追加模式，打补丁模式。输出文件名为data/01/a.out.idx和data/01/a.out.img
# 此模式的 -o 可以是使用append后的img文件
# 1072887808 为原cdimg0.img大小
TMGS3Tools -idx=data/_cdimg.idx \
           -img=data/_cdimg0.img \
           -i=data/output \
           -o=data/cdimg0.img \
           -o2=data/cdimg.idx \
           -patch=1072887808 \
           -append 

# OFS3单独解包，输出到data/ofs3/output，开启日志
TMGS3Tools -ofs3=data/ofs3/005.bin \
           -o=data/ofs3/output \
           -ofs3.log
           
# OFS3单独打包，输出到data/ofs3/005.out.bin，开启日志
TMGS3Tools -ofs3=data/ofs3/005.bin \
           -i=data/ofs3/output \
           -o=data/ofs3/005.out.bin # -ofs3.log

# 解包，输出到data/01/output文件夹中
TMGS3Tools -idx=data/01/a.idx \
           -img=data/01/a.img \
           -o=data/01/output 

# 解包（含OFS3文件），输出到data/01/output文件夹中
TMGS3Tools -idx=data/01/a.idx \
           -img=data/01/a.img \
           -o=data/01/output \
           -dfi.ofs3 # -ofs3.log

# 打包 普通模式，输出文件名为data/01/a.out.idx和data/01/a.out.img
TMGS3Tools -idx=data/01/a.idx \
           -img=data/01/a.img \
           -i=data/01/output \
           -o=data/01/a.out  

# 打包 追加模式，输出文件名为data/01/a.out.idx和data/01/a.out.img
TMGS3Tools -idx=data/01/a.idx \
           -img=data/01/a.img \
           -i=data/01/output \
           -o=data/01/a.out \
           -append 

```
## 未来

## 更新日志

### 2022-4-8
- 支持patch模式

### 2022-4-7 2
- 取消支持自动压缩gz数据

### 2022-4-7 1
- 支持OFS3格式解包打包时自动解压、压缩gz数据

### 2022-4-6
- 支持OFS3格式文件递归打包

### 2022-4-5 3
- 支持OFS3格式文件递归解包

### 2022-4-5 2
- 任意DFI(idx img)文件打包
- 支持追加模式打包

### 2022-4-5 1
- 任意DFI(idx img)文件提取

### 2022-3-26
- DFI(idx img)文件提取，仅GS3 cdimg0文件

## 参考

- OFS3-TOOL [https://github.com/Liquid-S/OFS3-TOOL](https://github.com/Liquid-S/OFS3-TOOL)  