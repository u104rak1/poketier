---
applyTo: '**/*_test.go'
---

# テストの方針

## テストサイズ
Googleが提唱しているテストサイズ分類に基づき、Smallサイズのテストをユニットテスト、Mediumサイズのテストを統合テストとして実装します。
E2Eテストは実装しません。

| サイズ | 実行時間 | ネットワーク | ファイルシステム | データベース | 外部システム |
|--------|----------|--------------|----------------|--------------|--------------|
| **Small** | < 1分 | ❌ 禁止 | ❌ 禁止 | ❌ 禁止 | ❌ 禁止 |
| **Medium** | < 5分 | ✅ localhost のみ | ✅ 許可 | ✅ 許可 | ❌ 禁止 |
| **Large** | < 15分 | ✅ 許可 | ✅ 許可 | ✅ 許可 | ✅ 許可 |

## ユニットテスト

### 基本方針

- ホワイトボックステストを基本とし、コードの内部構造を理解した上でテストを設計します。
- テーブルドリブン形式で書き、様々なケースをテストします。ただしテストケースが一つしかない場合や分岐コードが多くなる場合は無理にテーブルドリブンテストにする必要はありません。ただしサブテストにしてテストケースは日本語で明確に説明します。
- テストケースには日本語で明確な説明を付けて、仕様確認を容易にします。また正常系・異常系どちらかのプレフィックスを付けて、テストの意図を明確にします。
- テストは並行実行（Parallel）を基本とし、テスト間の干渉を避けるようにします。
- Arrange-Act-Assert（3A）方式でテストを構造化し、各ステップの冒頭に `Arrange` `Act` `Assert` のコメントを付けます。コードのどこで準備が始まる。実行される。検証されるを分かりやすくします。 `Arrange` のみ不要なら省略可です。
- testifyなどのアサーションライブラリを使用して、期待値とのマッチングを明確にします。
- 境界値テストを実施して、エッジケースやエラーケースを確認します。型はGoの言語使用を信頼するので、異なる型に関するテストは省略しても構いません。
- カバレッジは80%以上を目安にしますが、再現が困難な分岐やMockを多用するような実施の意味が薄いテストについては省略しても構いません。
- テストでは定数や変数をハードコードし、プロダクトコードの変更に敏感に反応できるようにします。
- テストファイルは対応するプロダクトコードと同じディレクトリに配置し、ファイル名は`*_test.go`とします。package名は{pkg}_testとします。
- 環境変数が必要な場合はt.Setenvを使用して設定します。後処理も必ず行います。
- リポジトリのテストでは `github.com/DATA-DOG/go-sqlmock` を使用して、SQLのMockを行います。これにより、実際のデータベースに依存せずにテストを実行できます。

例
```go
func TestExample(t *testing.T) {
    t.Parallel() // テストを並行実行

    t.Setenv("ENV_VAR", "value") // 必要なら環境変数を設定
    defer t.Unsetenv("ENV_VAR")

    tests := []struct {
        caseName string
        input    string
        want     string
        err      error
    }{
        {
          caseName: "正常系: 小文字が大文字に変換される事",
          input:    "input",
          want:     "INPUT",
          err:      nil,
        },
        {
          caseName: "異常系: 空文字が渡された場合",
          input:    "",
          want:     "",
          err:      errors.New("error: empty input"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.caseName, func(t *testing.T) {
            t.Parallel() // サブテストも並行実行

            // Arrange <= 準備のコードがどこにあるかを明確にする為のコメント
            // ここに必要な準備コードを記述（この例では特に無し）

            // Act <= 実行のコードがどこにあるかを明確にする為のコメント
            got, err := ConvertToUpper(tt.input)

            // Assert <= 検証のコードがどこにあるかを明確にする為のコメント
            if tt.err != nil {
                assert.Error(t, err, "expected error but got none")
                return
            }
            assert.NoError(t, err, "unexpected error occurred")
            assert.Equal(t, tt.want, got, "got does not match expected value")
        })
    }
}
```
