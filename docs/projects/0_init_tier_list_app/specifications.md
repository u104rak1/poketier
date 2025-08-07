# PokeTier - システム仕様書 (Phase 0: 初期ティアリスト作成システム)

## 1. システム概要

### 1.1 システム構成図
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   User          │    │  Cloudflare     │    │   Next.js       │
│   (Browser)     │◄──►│   Pages/CDN     │◄──►│   Frontend      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                       │
                                                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Cloudflare    │    │   Cloudflare    │    │   Cloudflare    │
│   R2 Storage    │◄──►│   Workers/Go    │◄──►│   D1 Database   │
│   (Images)      │    │   (API Server)  │    │   (SQLite)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### 1.2 技術スタック詳細（Phase 0）
- **フロントエンド**: Next.js 15+ (SSG), TypeScript 5.6+, Tailwind CSS
- **バックエンド**: Go 1.23+
- **データベース**: Cloudflare D1 (SQLite互換) (仮)
- **ストレージ**: Cloudflare R2 (画像保存) (仮)
- **インフラ**: Cloudflare Pages/Workers (無料枠) (仮)
- **CI/CD**: GitHub Actions

## 2. データベース仕様（ユビキタス言語準拠）

### 2.1 テーブル設計

#### 2.1.1 seasons（シーズン集約）
| カラム名 | データ型 | 制約 | 説明 |
|---------|----------|------|------|
| season_id | TEXT | PRIMARY KEY | UUID文字列形式 |
| name | TEXT | NOT NULL | シーズン名 |
| start_date | INTEGER | NOT NULL | UNIX timestamp |
| end_date | INTEGER | | UNIX timestamp (null = 進行中) |
| is_active | INTEGER | DEFAULT 0, CHECK (0,1) | boolean |
| created_at | INTEGER | DEFAULT (strftime('%s', 'now')) | 作成日時 |
| updated_at | INTEGER | DEFAULT (strftime('%s', 'now')) | 更新日時 |

**インデックス**:
- `idx_seasons_active`: is_active (UNIQUE WHERE is_active = 1)
- `idx_seasons_dates`: start_date, end_date

#### 2.1.2 expansions（拡張パック集約）
| カラム名 | データ型 | 制約 | 説明 |
|---------|----------|------|------|
| expansion_id | TEXT | PRIMARY KEY | UUID文字列形式 |
| name | TEXT | NOT NULL | 拡張パック名 |
| code | TEXT | NOT NULL UNIQUE | 略称コード |
| release_date | INTEGER | NOT NULL | UNIX timestamp |
| is_active | INTEGER | DEFAULT 1, CHECK (0,1) | 販売中フラグ |
| series | TEXT | NOT NULL | シリーズ名 |
| created_at | INTEGER | DEFAULT (strftime('%s', 'now')) | 作成日時 |
| updated_at | INTEGER | DEFAULT (strftime('%s', 'now')) | 更新日時 |

**インデックス**:
- `idx_expansions_series`: series
- `idx_expansions_active`: is_active (WHERE is_active = 1)

#### 2.1.3 cards（カード集約）
| カラム名 | データ型 | 制約 | 説明 |
|---------|----------|------|------|
| card_id | TEXT | PRIMARY KEY | UUID文字列形式 |
| expansion_id | TEXT | NOT NULL, FK | 拡張パックID |
| name | TEXT | NOT NULL | カード名 |
| image_url | TEXT | NOT NULL | カード画像URL |
| category | TEXT | NOT NULL, CHECK | 'pokemon', 'trainer' |
| type | TEXT | NOT NULL, CHECK | カードタイプ |
| rarity | TEXT | NOT NULL, CHECK | 'dia1', 'dia2', 'dia3', 'dia4', 'promo' |
| created_at | INTEGER | DEFAULT (strftime('%s', 'now')) | 作成日時 |
| updated_at | INTEGER | DEFAULT (strftime('%s', 'now')) | 更新日時 |

**インデックス**:
- `idx_cards_expansion`: expansion_id
- `idx_cards_name`: name
- `idx_cards_type`: type
- `idx_cards_category`: category

#### 2.1.4 decks（デッキ集約）
| カラム名 | データ型 | 制約 | 説明 |
|---------|----------|------|------|
| deck_id | TEXT | PRIMARY KEY | UUID文字列形式 |
| season_id | TEXT | NOT NULL, FK | シーズンID |
| primary_card_id | TEXT | NOT NULL, FK | 1番目のカードID |
| secondary_card_id | TEXT | FK | 2番目のカードID |
| tertiary_card_id | TEXT | FK | 3番目のカードID |
| nickname | TEXT | NOT NULL | デッキニックネーム |
| card_names | TEXT | NOT NULL | 検索用カンマ区切り |
| image_url | TEXT | | 合成画像URL |
| created_at | INTEGER | DEFAULT (strftime('%s', 'now')) | 作成日時 |
| updated_at | INTEGER | DEFAULT (strftime('%s', 'now')) | 更新日時 |

**インデックス**:
- `idx_decks_season`: season_id
- `idx_decks_primary_card`: primary_card_id
- `idx_decks_nickname`: nickname

#### 2.1.5 tier_lists（ティアリスト集約）
| カラム名 | データ型 | 制約 | 説明 |
|---------|----------|------|------|
| tier_list_id | TEXT | PRIMARY KEY | UUID文字列形式 |
| title | TEXT | NOT NULL | タイトル |
| description | TEXT | | 説明 |
| season_id | TEXT | NOT NULL, FK | シーズンID |
| author_name | TEXT | DEFAULT '匿名ユーザー' | 作成者名 |
| view_count | INTEGER | DEFAULT 0 | 閲覧数 |
| created_at | INTEGER | DEFAULT (strftime('%s', 'now')) | 作成日時 |
| updated_at | INTEGER | DEFAULT (strftime('%s', 'now')) | 更新日時 |

**インデックス**:
- `idx_tier_lists_season`: season_id
- `idx_tier_lists_created`: created_at DESC

#### 2.1.6 tier_placements（ティア配置）
| カラム名 | データ型 | 制約 | 説明 |
|---------|----------|------|------|
| tier_placement_id | TEXT | PRIMARY KEY | UUID文字列形式 |
| tier_list_id | TEXT | NOT NULL, FK | ティアリストID |
| deck_id | TEXT | NOT NULL, FK | デッキID |
| tier_rank | INTEGER | NOT NULL, CHECK (1-7) | 1=E, 7=SS |
| position | INTEGER | NOT NULL DEFAULT 0 | ティア内順序 |
| created_at | INTEGER | DEFAULT (strftime('%s', 'now')) | 作成日時 |

**制約**:
- `UNIQUE(tier_list_id, deck_id)`: 同一ティアリスト内での重複配置禁止

**インデックス**:
- `idx_tier_placements_tier_list`: tier_list_id
- `idx_tier_placements_deck`: deck_id
- `idx_tier_placements_rank`: tier_rank

#### 2.1.7 tier_statistics（ティア統計）
| カラム名 | データ型 | 制約 | 説明 |
|---------|----------|------|------|
| deck_id | TEXT | NOT NULL, FK | デッキID |
| season_id | TEXT | NOT NULL, FK | シーズンID |
| tier_rank | REAL | NOT NULL | 平均ティアランク |
| placement_count | INTEGER | NOT NULL DEFAULT 0 | 配置回数 |
| calculated_at | INTEGER | DEFAULT (strftime('%s', 'now')) | 計算日時 |

**制約**:
- `PRIMARY KEY (deck_id, season_id)`: 複合主キー

**インデックス**:
- `idx_tier_statistics_season`: season_id
- `idx_tier_statistics_rank`: tier_rank DESC

## 3. API仕様（ユビキタス言語準拠）

### 3.1 REST エンドポイント設計

#### 3.1.1 シーズン管理API（SeasonAggregate）
```
GET    /api/v1/seasons                   - ListSeasons
GET    /api/v1/seasons/active            - GetActiveSeason  
GET    /api/v1/seasons/{season_id}       - GetSeason
```

#### 3.1.2 拡張パック・カード管理API
```
GET    /api/v1/expansions                - ListExpansions
GET    /api/v1/cards                     - ListCards
```

#### 3.1.3 デッキ管理API（DeckAggregate）
```
POST   /api/v1/decks                     - CreateDeck
GET    /api/v1/decks                     - ListDecks (季節別)
```

#### 3.1.4 ティアリストAPI（TierListAggregate）
```
GET    /api/v1/tier-lists                - ListTierLists
GET    /api/v1/tier-lists/{tier_list_id} - GetTierList
POST   /api/v1/tier-lists                - CreateTierList
PATCH  /api/v1/tier-lists/{tier_list_id} - UpdateTierList
```

#### 3.1.5 統計・集計API（StatisticsCalculation）
```
GET    /api/v1/consensus/{season_id}     - GetConsensusTierList
GET    /api/v1/statistics/tier/{deck_id} - GetTierStatistics
```

### 3.2 主要なリクエスト/レスポンス仕様

#### 3.2.1 POST /api/v1/decks - CreateDeck
**リクエスト**:
```json
{
  "season_id": "550e8400-e29b-41d4-a716-446655440000",
  "primary_card_id": "550e8400-e29b-41d4-a716-446655440001",
  "secondary_card_id": "550e8400-e29b-41d4-a716-446655440002",
  "tertiary_card_id": null,
  "nickname": "リザニンフ"
}
```

**レスポンス**:
```json
{
  "deck_id": "550e8400-e29b-41d4-a716-446655440003",
  "season_id": "550e8400-e29b-41d4-a716-446655440000",
  "primary_card_id": "550e8400-e29b-41d4-a716-446655440001",
  "secondary_card_id": "550e8400-e29b-41d4-a716-446655440002",
  "tertiary_card_id": null,
  "nickname": "リザニンフ",
  "card_names": "リザードンex,ニンフィアex",
  "image_url": "https://r2.example.com/decks/550e8400-e29b-41d4-a716-446655440003.png",
  "created_at": 1691513600
}
```

#### 3.2.2 POST /api/v1/tier-lists - CreateTierList
**リクエスト**:
```json
{
  "title": "8月環境ティアリスト",
  "description": "新弾環境での評価",
  "season_id": "550e8400-e29b-41d4-a716-446655440000",
  "author_name": "配信者A",
  "placements": [
    {
      "deck_id": "550e8400-e29b-41d4-a716-446655440003",
      "tier_rank": 7,
      "position": 0
    }
  ]
}
```

**レスポンス**:
```json
{
  "tier_list_id": "550e8400-e29b-41d4-a716-446655440004",
  "title": "8月環境ティアリスト",
  "description": "新弾環境での評価",
  "season_id": "550e8400-e29b-41d4-a716-446655440000",
  "author_name": "配信者A",
  "view_count": 0,
  "created_at": 1691513600,
  "placements": [
    {
      "tier_placement_id": "550e8400-e29b-41d4-a716-446655440005",
      "deck_id": "550e8400-e29b-41d4-a716-446655440003",
      "tier_rank": 7,
      "position": 0,
      "deck": {
        "deck_id": "550e8400-e29b-41d4-a716-446655440003",
        "nickname": "リザニンフ",
        "image_url": "https://r2.example.com/decks/550e8400-e29b-41d4-a716-446655440003.png"
      }
    }
  ]
}
```

#### 3.2.3 GET /api/v1/consensus/{season_id} - GetConsensusTierList
**レスポンス**:
```json
{
  "season_id": "550e8400-e29b-41d4-a716-446655440000",
  "generated_at": 1691513600,
  "total_tier_lists": 25,
  "tiers": {
    "SS": [
      {
        "deck_id": "550e8400-e29b-41d4-a716-446655440003",
        "nickname": "リザニンフ",
        "image_url": "https://r2.example.com/decks/550e8400-e29b-41d4-a716-446655440003.png",
        "average_tier_rank": 6.8,
        "placement_count": 20
      }
    ],
    "S": [],
    "A": [],
    "B": [],
    "C": [],
    "D": [],
    "E": []
  }
}
```

### 3.3 エラーハンドリング仕様

#### 3.3.1 HTTPステータスコード
- `200 OK`: 正常処理
- `201 Created`: リソース作成成功
- `400 Bad Request`: リクエスト形式エラー
- `404 Not Found`: リソース未発見
- `500 Internal Server Error`: サーバーエラー

#### 3.3.2 エラーレスポンス形式
```json
{
  "error": {
    "code": "INVALID_TIER_RANK",
    "message": "tier_rank must be between 1 and 7",
    "details": {
      "field": "tier_rank",
      "value": 8,
      "allowed_values": [1, 2, 3, 4, 5, 6, 7]
    }
  }
}
```

## 4. 画面仕様（Phase 0）

### 4.1 画面一覧

#### 4.1.1 トップページ（/）
**目的**: サービス概要と最新ティアリストの表示
**表示内容**:
- ヒーローセクション（PokeTier説明）
- アクティブシーズン情報（GetActiveSeason）
- アクティブシーズンの集計ティアリスト表示
- ティアリスト作成画面へのボタン
- アクティブシーズンのティアリスト一覧 view_countの人気順（ListTierLists）

#### 4.1.2 ティアリスト作成画面（/create）
**目的**: TierListCreation プロセスの実行
**表示内容**:
- タイトル・説明入力フォーム
- デッキ検索・選択機能（ListDecks）
- ドラッグ&ドロップによるティア配置エリア（SS/S/A/B/C/D/E）
- リアルタイムプレビュー
**その他**:
- デッキがない場合はその場でモーダル表示で作れる。

#### 4.1.3 ティアリスト詳細画面（/tier-lists/:id）
**目的**: TierList の詳細表示
**表示内容**:
- ティアリスト情報（タイトル、作成者、作成日時）
- ティア配置の表示（TierPlacement）
- 閲覧数統計
- シンプルな共有機能

### 4.2 レスポンシブ対応（Phase 0）
- **Mobile**: 320px - 767px（基本対応）
- **Desktop**: 768px+（メイン対応）

## 5. 画像処理仕様（ImageComposition）

### 5.1 デッキ画像合成（Phase 0）
**目的**: 複数カードから1つのデッキ画像を生成
**処理フロー**:
1. カード画像取得（1-3枚）
2. レイアウト決定（横並び配置）スラッシュ区切りでカッコよくする
3. 画像合成処理
4. Cloudflare R2 へアップロード
5. デッキエンティティに URL 保存
6. 4:3くらいの横長
7. 下部にはデッキのニックネームを表示

**画像仕様**:
- **サイズ**: 800x600px（固定）
- **フォーマット**: webp
- **レイアウト**: 横並び、等間隔配置
- **品質**: Web表示最適化

### 5.2 画像キャッシュ戦略
- **TTL**: 24時間（Cloudflare R2）
- **命名規則**: `decks/{deck_id}.png`

## 6. デプロイメント仕様（Phase 0）

### 6.1 Cloudflare 構成 (仮)
**Cloudflare Pages（フロントエンド）**:
- Next.js SSG ビルド
- 自動デプロイ（GitHub連携）

**Cloudflare Workers（バックエンド）**:
- Go コンパイル→WASM
- D1 データベース接続
- R2 ストレージ接続

**Cloudflare D1（データベース）**:
- SQLite 互換
- 無料枠: 100MB

**Cloudflare R2（ストレージ）**:
- 画像ファイル保存
- 無料枠: 10GB

### 6.3 環境設定
**開発環境**:
- Local D1 database
- Local R2 storage (minio)
- wrangler dev

**本番環境**:
- Cloudflare D1
- Cloudflare R2
- カスタムドメイン（予定）

## 7. セキュリティ仕様（Phase 0）

### 7.1 基本セキュリティ
- **HTTPS強制**: Cloudflare SSL
- **CORS設定**: フロントエンドドメインのみ許可
- **入力検証**: 全APIでバリデーション実装
- **SQLインジェクション対策**: パラメーター化クエリ使用

### 7.2 レート制限（Phase 0）
- **API制限**: 100リクエスト/分/IP
- **画像生成制限**: 10リクエスト/分/IP
- **Cloudflare Workers制限**: CPU時間10ms

## 8. 監視・ログ仕様（Phase 0）

### 8.1 基本監視
- **Cloudflare Analytics**: トラフィック監視
- **Workers Analytics**: API使用状況
- **D1 Analytics**: データベースパフォーマンス

### 8.2 ログ設定
- **アプリケーションログ**: console.log（Workers）
- **アクセスログ**: Cloudflare標準ログ
- **エラーログ**: sentry.io（将来検討）
