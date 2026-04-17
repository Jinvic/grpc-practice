set windows-shell := ["pwsh.exe", "-c"]
set shell := ["/bin/bash", "-uc"]

default_port := "8081"
default_book := '{"title":"The Great Gatsby","author":"F. Scott Fitzgerald","price":10.99,"isbn":"978-0-7432-1967-1","publisher":"Scribner","published_at":"2026-04-17T00:00:00Z"}'
default_book_update := '{"id":1,"title":"The Great Gatsby","author":"F. Scott Fitzgerald","price":10.99,"isbn":"978-0-7432-1967-1","publisher":"Scribner","published_at":"2026-04-17T00:00:00Z"}'
default_book_update_mask := '{"paths":["title","author","price","isbn","publisher","published_at"]}'

default:
  just --list

describe arg="book.v1.BookService" port=default_port:
  grpcurl -plaintext localhost:{{port}} describe {{arg}}

GetBook id="1" port=default_port:
  grpcurl -plaintext -d '{"id":{{id}}}' localhost:{{port}} book.v1.BookService/GetBook

CreateBook book=default_book port=default_port:
  grpcurl -plaintext -d '{"book":{{book}}}' localhost:{{port}} book.v1.BookService/CreateBook

ListBooks page_number="1" page_size="10" port=default_port:
  grpcurl -plaintext -d '{"pageNumber":{{page_number}}, "pageSize":{{page_size}}}' localhost:{{port}} book.v1.BookService/ListBooks

UpdateBook book=default_book_update mask=default_book_update_mask port=default_port:
  grpcurl -plaintext -d '{"book":{{book}}, "updateMask":{{mask}}}' localhost:{{port}} book.v1.BookService/UpdateBook

DeleteBook id="1" port=default_port:
  grpcurl -plaintext -d '{"id":{{id}}}' localhost:{{port}} book.v1.BookService/DeleteBook
