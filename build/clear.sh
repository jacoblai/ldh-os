# 清理内核编译
cd ../kernel
make clean    # 清理大部分编译文件
make mrproper # 完全清理，包括.config文件

# 清理构建输出目录
rm -rf ./output/*