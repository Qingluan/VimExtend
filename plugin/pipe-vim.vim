nnoremap <c-k> :PipeTo

function! PipeToNewBufGo(cmd)
    set splitright
    let output=system('cat << EOF | '.a:cmd."\n". join(getline(1, '$'), "\n") )
    vnew
    call setline(1, split(output, "\n"))
endfunction

function! FuncList(ArgLead, cmdline, cursorpos ) abort
    return join(["VimExtend", "sort"], "\n")
endfunction

command! -nargs=1 -complete=custom,FuncList PipeTo call PipeToNewBufGo(<q-args>)