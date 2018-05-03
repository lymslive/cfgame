#! /bin/env bash

# 生成程序名同项目名，即父目录名
cwd=$(pwd)
project=$(basename $cwd)
echo [$project] $cwd 

# 备份上次生成的程序
if [[ -x $project ]]; then
	echo "bakup executable $project"
	mv $project $project.bak
fi

# 在当前目录生成程序
echo "go build new executable $project"
go build
