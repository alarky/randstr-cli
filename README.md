# randstr

A CLI tool for generating cryptographically secure random strings.

暗号論的に安全なランダム文字列を生成する CLI ツール。

## Installation / インストール

Download a binary from [Releases](https://github.com/alarky/randstr-cli/releases).

[Releases](https://github.com/alarky/randstr-cli/releases) からバイナリをダウンロードしてください。

Or build from source / またはソースからビルド:

```bash
go install github.com/alarky/randstr-cli@latest
```

## Usage / 使い方

```bash
randstr [flags]
```

### Flags / フラグ

| Flag | Default | Description |
|------|---------|-------------|
| `-c` | `12` | Number of strings to generate / 生成する文字列の個数 |
| `-n` | `16` | Length of each string / 各文字列の長さ |
| `-a` | `false` | Include lowercase letters (a-z) / 小文字を含める |
| `-A` | `false` | Include uppercase letters (A-Z) / 大文字を含める |
| `-0` | `false` | Include digits (0-9) / 数字を含める |
| `-s` | `false` | Include all symbols / 全記号を含める |
| `-symbols` | `""` | Include specific symbols / 特定の記号を含める |
| `-symbol-max` | `20` | Max percentage of symbols in output / 記号の最大割合(%) |

When no character type flags are specified, a-z + A-Z + 0-9 are used by default.

文字種フラグを指定しない場合、デフォルトで a-z + A-Z + 0-9 が使用されます。

### Examples / 使用例

```bash
# Default: 12 strings of 16 alphanumeric chars
# デフォルト: 英数字16文字 × 12個
randstr

# 5 strings of 32 chars
# 32文字 × 5個
randstr -c 5 -n 32

# Lowercase only
# 小文字のみ
randstr -a

# Uppercase + digits
# 大文字 + 数字
randstr -A -0

# With all symbols (max 20% symbols)
# 全記号付き (記号は最大20%)
randstr -s

# With specific symbols, symbol ratio up to 50%
# 特定の記号を指定、記号割合を最大50%に
randstr -symbols '!@#$' -symbol-max 50
```

## Development / 開発

```bash
make build        # Build binary / バイナリをビルド
make test         # Run tests / テストを実行
make cross-build  # Cross-build for all platforms / 全プラットフォーム向けクロスビルド
make release      # Create GitHub release / GitHub リリースを作成 (requires git tag)
make clean        # Remove build artifacts / ビルド成果物を削除
```

## License

MIT
