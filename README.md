## simples

* sync Mac Notes to web page


![img.png](img.png)
 
#### 简单干净

* https://firsh.me


#### quick start

* install python3 and deps

* install golang 1.25.0 or later

* web page server

```bash
go run main.go
```

* Gen Markdown file from Mac Notes

> modify `HOME_BASE` and os.system(...) in notes_export.py

```bash
python watch_notes_service/notes_export.py
```

* open `http://localhost:8080`


#### 感谢（复制了这个大佬的代码）

* https://github.com/keithvassallomt/taskbridge