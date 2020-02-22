nnmap <c-k> :'<,'>call  RunInGo()

function! RunInGo() range
    echo system('echo '.shellescape(join(getline(a:firstline, a:lastline), "\n")).'| VimExtend')
endfunction