nnoremap <c-k> :PipeTo
nnoremap <c-s> :% !VimExtend -r false  -q

let g:Cmds=["VimExtend -r", "VimExtend", "sort", "text", "StartServer"]
let g:if_start_proxy_server=0
let g:if_set_proxy_listen=0
" let g:last_cmd=""
" let g:cur_cmd=""
" augroup testgroup
"     autocmd  BufEnter * :echom "Baz"
" augroup END
set splitbelow
function! PipeToNewBufGo(cmd)
    if a:cmd == "StartServer"
        silent execute("ListProxy")
        " let g:last_cmd=cur_cmd
        " let g:cur_cmd="ListProxy"
        return
    endif
    setlocal splitright
    setlocal cursorline
    highlight Cursorline cterm=underline  ctermbg=green ctermfg=white
    if a:cmd == "sort"
        execute "% !".a:cmd
    elseif a:cmd == "text"
        execute "% !w3m -dump -T text/html"
        return
    "   execute "PipeTo ". 'w3m -dump -T text/html'
    else
        let l:cmdStr='cat << EOF | '.a:cmd."\n". join(getline(1, '$'), "\n")
        " echo l:cmdStr
        let output=system(l:cmdStr)
        vnew pipe
        nnoremap <buffer>  q :q!  <CR>
        call setline(1, split(output, "\n"))
    endif
    nnoremap <buffer><Space> :call LineToggle() <CR>
    nnoremap <buffer>  q :q!  <CR>
endfunction

function! FuncList(ArgLead, cmdline, cursorpos ) abort
    return join(g:Cmds, "\n")
endfunction


function! LineToggle()
    let line=getline('.')
    if line =~ '^[\d'
        let l:data=split(line, " ")
        let l:output = system('VimExtend -g '.data[1])
        vnew
        nnoremap <buffer>  q :q!  <CR>
        nnoremap <buffer><Space> :call LineToggle() <CR>
        call setline(1, split(output, "\n"))
    elseif getline(1) =~ '^GET'
        let l:cmdStr='cat << EOF | VimExtend'."\n". join(getline(1, '$'), "\n")
        let output=system(l:cmdStr)
        new pipe
        nnoremap <buffer>  q :q!  <CR>
        call setline(1, split(output, "\n"))
    elseif getline(1) =~ '^POST'
        let l:cmdStr='cat << EOF | VimExtend'."\n". join(getline(1, '$'), "\n")
        let output=system(l:cmdStr)
        new pipe
        nnoremap <buffer>  q :q!  <CR>
        call setline(1, split(output, "\n"))
    elseif getline(1) =~ '^PUT'
        let l:cmdStr='cat << EOF | VimExtend'."\n". join(getline(1, '$'), "\n")
        let output=system(l:cmdStr)
        new pipe
        nnoremap <buffer>  q :q!  <CR>
        call setline(1, split(output, "\n"))
    elseif getline(1) =~ '^HEAD'
        let l:cmdStr='cat << EOF | VimExtend'."\n". join(getline(1, '$'), "\n")
        let output=system(l:cmdStr)
        new pipe
        nnoremap <buffer>  q :q!  <CR>
        call setline(1, split(output, "\n"))
    elseif getline(1) =~ '^DELETE'
        let l:cmdStr='cat << EOF | VimExtend'."\n". join(getline(1, '$'), "\n")
        let output=system(l:cmdStr)
        new pipe
        nnoremap <buffer>  q :q!  <CR>
        call setline(1, split(output, "\n"))
    
    endif
endfunction

if v:version >= 800
    let g:JobBack = {}
    fun! JobBack.on_stdout(job_id, data, event)
        if self.name == "ProxyServer"
            execute("ListProxy")
        endif
        if a:data != [""]
            " call setline(2,  shellescape(join(a:data, '\n')) )
            " echom shellescape(join(a:data, "\n"))
            let l:buf_no=len(getline("display-status","$"))+1
            for l in a:data
                if l == ""
                    continue
                endif
                call setbufline("display-status",l:buf_no,  l )
                let l:buf_no += 1
            endfor
            
        endif
    endfunction

    let JobBack.on_stderr = JobBack.on_stdout
    function JobBack.on_exit(job_id, _data, event)
      let msg = printf('job %d "%s" finished', a:job_id, self.name)
      call setbufline("display-status",1, printf('[%s] %s!', a:event, msg))
      bdelete! display-status
    endfunction

    fun! JobBack.new(name, cmd)
        if  bufwinnr("display-status") == -1
            new display-status
            resize 2
            "  if g:if_set_proxy_listen == 0 
            "      autocmd TextChanged <buffer>  :ListProxy
            "     let g:if_set_proxy_listen=1
            " endif
        endif
        let obj = extend(copy(g:JobBack), {'name': a:name})
        let obj.cmd = ["/bin/bash", "-c", a:cmd]
        let obj.id = jobstart(obj.cmd, obj)
        $
        return obj
    endfunction
    fun! g:StartJob(cmd)
        let g:run_job=g:JobBack.new("ProxyServer", a:cmd)
    endfunction
    command! -nargs=1 -complete=custom,FuncList Work call StartJob(<q-args>) 
    command! -nargs=0 -bang WorkStop call jobstop(g:run_job.id)
endif
fun! s:listProxy()
    
    if g:if_start_proxy_server == 0
        setlocal splitright
        setlocal cursorline
        nnoremap <buffer><Space> :call LineToggle() <CR>
        nnoremap <buffer>q :ListProxy <CR>
        nnoremap b<buffer>  :!VimExtend -l
        execute 'Work VimExtend -S'
        let g:if_start_proxy_server = 1
        
    endif
   
    let output=system('VimExtend -ls')
    silent call setbufline(1, 1, split(output,"\n" ))
endfunction


command! -nargs=1 -complete=custom,FuncList PipeTo call PipeToNewBufGo(<q-args>)
command! -nargs=0 ListProxy call s:listProxy()
