let SessionLoad = 1
let s:so_save = &g:so | let s:siso_save = &g:siso | setg so=0 siso=0 | setl so=-1 siso=-1
let v:this_session=expand("<sfile>:p")
silent only
silent tabonly
cd /Dev/go-revise
if expand('%') == '' && !&modified && line('$') <= 1 && getline(1) == ''
  let s:wipebuf = bufnr('%')
endif
let s:shortmess_save = &shortmess
if &shortmess =~ 'A'
  set shortmess=aoOA
else
  set shortmess=aoO
endif
badd +15 internal/domain/reviseitem/repository_sqlite.go
badd +5 /Dev/go-revise/sqlc.yaml
badd +3 /Dev/go-revise/internal/db/migrations/000001_init.up.sql
badd +12 /Dev/go-revise/internal/domain/user/repository.go
badd +4 /Dev/go-revise/internal/db/queries/user.sql
badd +1 /Dev/go-revise/internal/db/sqlc/db.go
badd +1 /Dev/go-revise/internal/db/sqlc/models.go
badd +1 /Dev/go-revise/internal/db/sqlc/user.sql.go
badd +12 /Dev/go-revise/internal/domain/user/repository_sqlite.go
badd +43 /Dev/go-revise/internal/domain/user/user.go
argglobal
%argdel
edit /Dev/go-revise/internal/domain/user/repository_sqlite.go
let s:save_splitbelow = &splitbelow
let s:save_splitright = &splitright
set splitbelow splitright
wincmd _ | wincmd |
vsplit
1wincmd h
wincmd w
wincmd _ | wincmd |
split
1wincmd k
wincmd w
let &splitbelow = s:save_splitbelow
let &splitright = s:save_splitright
wincmd t
let s:save_winminheight = &winminheight
let s:save_winminwidth = &winminwidth
set winminheight=0
set winheight=1
set winminwidth=0
set winwidth=1
exe 'vert 1resize ' . ((&columns * 30 + 158) / 316)
exe '2resize ' . ((&lines * 36 + 38) / 76)
exe 'vert 2resize ' . ((&columns * 285 + 158) / 316)
exe '3resize ' . ((&lines * 36 + 38) / 76)
exe 'vert 3resize ' . ((&columns * 285 + 158) / 316)
argglobal
enew
file NvimTree_1
balt /Dev/go-revise/internal/db/sqlc/db.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal nofen
wincmd w
argglobal
balt /Dev/go-revise/internal/domain/user/user.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let &fdl = &fdl
let s:l = 12 - ((5 * winheight(0) + 18) / 36)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 12
normal! 0
wincmd w
argglobal
enew | setl bt=help
help session-file@en
balt /Dev/go-revise/internal/domain/user/repository_sqlite.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal nofen
silent! normal! zE
let &fdl = &fdl
let s:l = 800 - ((6 * winheight(0) + 18) / 36)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 800
normal! 059|
wincmd w
3wincmd w
exe 'vert 1resize ' . ((&columns * 30 + 158) / 316)
exe '2resize ' . ((&lines * 36 + 38) / 76)
exe 'vert 2resize ' . ((&columns * 285 + 158) / 316)
exe '3resize ' . ((&lines * 36 + 38) / 76)
exe 'vert 3resize ' . ((&columns * 285 + 158) / 316)
tabnext 1
if exists('s:wipebuf') && len(win_findbuf(s:wipebuf)) == 0 && getbufvar(s:wipebuf, '&buftype') isnot# 'terminal'
  silent exe 'bwipe ' . s:wipebuf
endif
unlet! s:wipebuf
set winheight=1 winwidth=20
let &shortmess = s:shortmess_save
let &winminheight = s:save_winminheight
let &winminwidth = s:save_winminwidth
let s:sx = expand("<sfile>:p:r")."x.vim"
if filereadable(s:sx)
  exe "source " . fnameescape(s:sx)
endif
let &g:so = s:so_save | let &g:siso = s:siso_save
set hlsearch
nohlsearch
doautoall SessionLoadPost
unlet SessionLoad
" vim: set ft=vim :
