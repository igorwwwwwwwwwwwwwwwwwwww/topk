# topk

A small utility for ad-hoc top-k analysis built in go.

## Install

```
$ go install github.com/igorwwwwwwwwwwwwwwwwwwww/topk@latest
```

## Usage

You can pass a filename:

```
$ topk types.txt
Water        112  ∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
Normal        98  ∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
Grass         70  ∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
Bug           69  ∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
Psychic       57  ∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
Fire          52  ∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
Rock          44  ∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
Electric      44  ∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
Ghost         32  ∎∎∎∎∎∎∎∎∎∎∎∎∎∎
Dragon        32  ∎∎∎∎∎∎∎∎∎∎∎∎∎∎

$ $topk -other types.txt
Water        112  ∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
Normal        98  ∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
Grass         70  ∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
Bug           69  ∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
Psychic       57  ∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
Fire          52  ∎∎∎∎∎∎∎∎∎∎∎∎∎∎
Rock          44  ∎∎∎∎∎∎∎∎∎∎∎∎
Electric      44  ∎∎∎∎∎∎∎∎∎∎∎∎
Ghost         32  ∎∎∎∎∎∎∎∎
Dragon        32  ∎∎∎∎∎∎∎∎
OTHER        191  ∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
```

Alternatively, it's also possible to pipe from stdin:

```
$ cat types.txt | topk
```

## Why

Because my fingers hurt from typing `sort | uniq -c | sort -rn | head`.

## Similar to

- [logtop](https://github.com/igorwwwwwwwwwwwwwwwwwwww/logtop)
- [log2hist](https://github.com/igorwwwwwwwwwwwwwwwwwwww/log2hist)

## License

MIT.
