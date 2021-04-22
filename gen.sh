#!/bin/sh

deps() {
go list -f "{{range .Deps}}{{.}}
{{end}}" "$@" \
  | sort -u \
  | grep ^cmd/ \
  | grep -v ^cmd/vendor/
}

gofiles() {
go list -f '{{ $dir := .Dir}}{{range .GoFiles}}{{$dir}}/{{.}}
{{end}}{{range .CgoFiles}}{{$dir}}/{{.}}
{{end}}' "$@"
}

for v in 3 4 5 6 7 8 9 10 11 12 13 14 15 16; do
  git -C ~/wd/go checkout -q go1.$v

  echo "== all cmd (go1.$v) =="
  gofiles $(deps cmd) | ./look | sort | uniq -c
  echo
  echo "== compile+link (go1.$v) =="
  gofiles $(deps cmd/{compile,link}) | ./look | sort | uniq -c
  echo
done
