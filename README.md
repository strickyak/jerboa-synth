# jerboa-synth
Jerboa: simulate basic synthesizer modules in GoLang

```
$ go run main.go -db=0 -r=48000 -v=0 | paplay --rate=48000 --channels=1 --format=s16le --raw /dev/stdin
```
