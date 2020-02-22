nnoremap <c-k> :PipeTo VimExtend -r
nnoremap <c-s> :% !VimExtend -q
let g:Cmds=["VimExtend -r", "VimExtend", "sort", "text"]
function! PipeToNewBufGo(cmd)
    set splitright
    set cursorline
    highlight Cursorline cterm=underline  ctermbg=green ctermfg=white
    if a:cmd == "sort"
    execute "% !".a:cmd
    elseif a:cmd == "text"
    execute "% !w3m -dump -T text/html"
    "   execute "PipeTo ". 'w3m -dump -T text/html'
    else
        let l:cmdStr='cat << EOF | '.a:cmd."\n". join(getline(1, '$'), "\n")
        " echo l:cmdStr
        let output=system(l:cmdStr)
        vnew
        call setline(1, split(output, "\n"))
    endif
    execute("nnoremap <Space> :call LineToggle() <CR>")
    execute("nnoremap  q :q!  <CR>")
endfunction

function! FuncList(ArgLead, cmdline, cursorpos ) abort
    return join(g:Cmds, "\n")
endfunction


function! LineToggle()
    let line=getline('.')
    echo line
endfunction

command! -nargs=1 -complete=custom,FuncList PipeTo call PipeToNewBufGo(<q-args>)