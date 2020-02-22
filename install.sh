#!/bin/bash

if [ ! -f /usr/bin/vim ];then
  apt install -y neovim 2>&1 1>/dev/null;
fi

if [ ! -d $GOPATH ];then
  echo 'export GOPATH="$HOME/go"' >> ~/.bashrc
  echo 'export PATH="$PATH:$GOPATH/bin"' >> ~/.bashrc
  source ~/.bashrc
  apt install -y golang
  go get -v "github.com/Qingluan/VimExtend"
fi

if [ -d ~/.vim/plugged ];then
  mkdir -p ~/.vim/plugged/VimExtend
else
  curl -fLo ~/.vim/autoload/plug.vim --create-dirs 'https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim'
  mkdir -p ~/.vim/plugged/VimExtend
fi

cp -a plugin ~/.vim/plugged/VimExtend/

cat << EOF >> ~/.vimrc
" Specify a directory for plugins
" - For Neovim: stdpath('data') . '/plugged'
" - Avoid using standard Vim directory names like 'plugin'
call plug#begin('~/.vim/plugged')

" Make sure you use single quotes

" Shorthand notation; fetches https://github.com/junegunn/vim-easy-align
Plug 'junegunn/vim-easy-align'

" Any valid git URL is allowed
Plug 'https://github.com/junegunn/vim-github-dashboard.git'

" Multiple Plug commands can be written in a single line using | separators
Plug 'SirVer/ultisnips' | Plug 'honza/vim-snippets'

" On-demand loading
Plug 'scrooloose/nerdtree', { 'on':  'NERDTreeToggle' }
Plug 'tpope/vim-fireplace', { 'for': 'clojure' }

" Using a non-master branch
Plug 'rdnetto/YCM-Generator', { 'branch': 'stable' }

" Using a tagged release; wildcard allowed (requires git 1.9.2 or above)
Plug 'fatih/vim-go', { 'tag': '*' }

" Plugin options
Plug 'nsf/gocode', { 'tag': 'v.20150303', 'rtp': 'vim' }

" Plugin outside ~/.vim/plugged with post-update hook
Plug 'junegunn/fzf', { 'dir': '~/.fzf', 'do': './install --all' }

" Unmanaged plugin (manually installed and updated)
Plug '~/my-prototype-plugin'

" Initialize plugin system
call plug#end()
EOF

mkdir -p ~/.config/nvim
cp ~/.vimrc ~/.config/nvim/init.vim

vim +PlugInstall +qall