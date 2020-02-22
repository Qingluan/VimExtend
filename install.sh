#!/bin/bash

if [ -d ~/.vim/plugged ];then
  mkdir -p ~/.vim/plugged/VimExtend
else
  curl -fLo ~/.vim/autoload/plug.vim --create-dirs 'https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim'
fi

cp -a plugin ~/.vim/plugged/VimExtend/
