# PokeTier ユビキタス言語辞典

## 概要
PokeTierプロジェクトで使用するドメイン固有の用語と概念を定義します。
コード内の変数名、関数名、API名は、この辞典の用語に準拠して命名します。

---

## ドメインエンティティ

### コアビジネス概念

#### PokePoke(ポケポケ)
**定義**: ポケモンカードのスマホアプリの略称 
**英語**: `pokepoke`
**日本語**: ポケポケ

---

#### Season（シーズン）
**定義**: ポケポケのランクマッチの環境期間を表す時間的な区切り  
**期間**: 通常1ヶ月。新弾リリースなどで切り替わる  
**英語**: `season`  
**日本語**: シーズン、環境期間  
**DB名**: `seasons`
**属性**:
- `season_id`: UUID - シーズン一意識別子
- `name`: string - シーズン名（例: A2b, A3）
- `is_active`: boolean - アクティブフラグ
- `start_date`: timestamp - 開始日時
- `end_date`: timestamp - 終了日時（null = 進行中）
**関連概念**:
- `SeasonMigration` - シーズン移行処理
- `ActiveSeason` - 現在アクティブなシーズン

---

#### Card（カード）
**定義**: ポケポケの個別カード、またはその画像の事  
**粒度**: 1枚のカード。同名カードでも拡張・レアリティが異なれば別エンティティ  
**英語**: `card`
**日本語**: カード  
**DB名**: `cards`
**属性**:
- `card_id`: UUID - カード一意識別子
- `expansion_id`: UUID - 拡張パックの一位識別子
- `name`: string - カード名（例: "ピカチュウex"）
- `image_url`: string - カード画像URL
- `category`: enum - カテゴリー（pokemon/trainer）
- `type`: enum - タイプ（grass/fire/water/lightning/psychic/fighting/darkness/metal/dragon/colorless/item/pokemon_tool/fossil/supporter）
- `rarity`: enum - レアリティ enum (dia1/dia2/dia3/dia4/promo)

**関連概念**:
- `CardExpansion` - 拡張パック

---

#### Expansion（拡張パック）
**定義**: ポケポケのカードパック・拡張セット  
**粒度**: 1つの拡張パック。例: 「黒炎の支配者」「レイジングサーフ」等  
**英語**: `expansion`  
**日本語**: 拡張パック  
**DB名**: `expansions`
**属性**:
- `expansion_id`: UUID - 拡張パック一意識別子
- `name`: string - 拡張パック名（例: "黒炎の支配者"）
- `code`: string - 略称コード（例: "sv4a"）
- `release_date`: date - リリース日
- `is_active`: boolean - 現在パック販売中フラグ
- `series`: string - シリーズ名（例: "スカーレット&バイオレット"）

**関連概念**:
- `ExpansionSeries` - 拡張シリーズ
- `CardCollection` - 収録カード一覧

---

#### Deck（デッキ）
**定義**: ティアリストで配置される個別の構築要素  
**構成**: 1つ以上のcardの組み合わせ（例: リザードンexのカード + ニンフィアexのカード = リザニンフ）  
**英語**: `deck`  
**日本語**: デッキ  
**DB名**: `decks`
**属性**:
- `deck_id`: UUID - デッキ一意識別子
- `season_id`: UUID - 所属シーズン
- `primary_card_id`: UUID - 1番目のカードのID
- `secondary_card_id`: UUID - 2番目のカードのID
- `tertiary_card_id`: UUID - 3番目のカードのID
- `nickname`: string - 表示名（例: "リザニンフ"）
- `card_names`: text - 使用したカード名のカンマ区切りの組み合わせ（検索用）
- `image_url`: string - 合成画像URL
**関連概念**:
- `DeckImage` - 合成画像

---

#### TierList（ティアリスト）
**定義**: ユーザーが作成するデッキ強度ランキング  
**形式**: デッキをSS/S/A/B/C/D/Eのティアに配置  
**英語**: `tier_list`  
**日本語**: ティアリスト  
**DB名**: `tier_lists`
**属性**:
- `tier_list_id`: UUID - ティアリスト一意識別子
- `title`: string - タイトル
- `description`: text - 説明
- `season_id`: UUID - 対象シーズン
- `author_name`: string - 作成者名
- `author_id`: UUID - 作成者ID（任意）
- `view_count`: integer - 閲覧数
**関連概念**:
- `TierListCreation` - 作成プロセス
- `TierListSharing` - 共有機能

---

#### TierPlacement（ティア配置）
**定義**: ティアリスト内でのデッキの配置情報  
**関係**: TierList と Deck の多対多関係  
**英語**: `tier_placement`  
**日本語**: ティア配置  
**DB名**: `tier_placements`
**属性**:
- `tier_placement_id`: UUID - 配置一意識別子
- `tier_list_id`: UUID - 所属ティアリスト
- `deck_id`: UUID - 配置されるデッキ
- `tier_rank`: integer - ティアランク
- `position`: integer - ティア内での順序
**関連概念**:
- `TierRank` - ティアランク列挙型
- `TierPosition` - ティア内配置順序

---

### 分析・統計概念

#### ConsensusTierList（集計ティアリスト）
**定義**: 全ユーザーの投稿から統計的に算出された総合ティアリスト  
**算出**: 重み付きスコア計算による自動生成  
**英語**: `consensus_tier_list`  
**日本語**: 集計ティアリスト 
**DB名**: 算出結果（永続化しない）

**関連概念**:
- `TierScore` - ティアスコア算出
- `WeightedCalculation` - 重み付き計算
- `ConsensusGeneration` - 集計処理

---

#### TierStatistics（ティア統計）
**定義**: デッキのティア配置に関する統計データ  
**期間**: シーズン単位で集計  
**英語**: `tier_statistics`
**日本語**: ティア統計  
**DB名**: `tier_statistics`（デッキ詳細システム）

**属性**:
- `deck_id`: UUID - 対象デッキ
- `season_id`: string - 対象シーズン
- `tier_rank`: integer - ランクの平均値。近似値のランクに配置される。
- `calculated_at`: timestamp - 計算日時

---

### 拡張機能概念

#### DeckDetail（デッキ詳細）
**定義**: デッキの詳細情報と分析結果  
**内容**: カードリスト、統計情報、メタ分析  
**英語**: `deck_detail`  
**日本語**: デッキ詳細

**関連概念**:
- `DeckCardList` - カードリスト詳細
- `DeckAnalysis` - デッキ分析
- `MatchupData` - 対戦相性データ

---

#### User（ユーザー）
**定義**: システム利用者アカウント  
**種類**: guest, user, moderator, admin  
**英語**: `user`  
**日本語**: ユーザー  

**関連概念**:
- `UserRole` - ユーザー権限
- `UserProfile` - プロフィール
- `UserCommunity` - コミュニティ参加

---

## 値オブジェクト・列挙型

### TierRank（ティアランク）
**定義**: ティアリストでの強度ランク  
**値**: `SS`, `S`, `A`, `B`, `C`, `D`, `E`  
**順序**: SS（最強）> S > A > B > C > D > E（最弱）  
**英語**: `tier_rank`  
**日本語**: ティアランク

**数値変換**:
```go
const (
    TierSS = 7
    TierS  = 6
    TierA  = 5
    TierB  = 4
    TierC  = 3
    TierD  = 2
    TierE  = 1
)
```

---

### CardCategory（カードカテゴリ）
**定義**: ポケポケカードの基本分類  
**値**: `pokemon`, `trainer`  
**英語**: `card_category`  
**日本語**: カードカテゴリ

---

### CardType（カードタイプ）
**定義**: ポケモンカードの詳細タイプ分類  
**英語**: `card_type`  
**日本語**: カードタイプ

**値**:
- **ポケモン**: `grass`, `fire`, `water`, `lightning`, `psychic`, `fighting`, `darkness`, `metal`, `dragon`, `colorless`
- **トレーナー**: `item`, `pokemon_tool`, `fossil`, `supporter`

---

### CardRarity（カードレアリティ）
**定義**: ポケポケカードのレアリティ分類  
**英語**: `card_rarity`  
**日本語**: カードレアリティ

**値**:
- `dia1`: ダイヤ1
- `dia2`: ダイヤ2  
- `dia3`: ダイヤ3
- `dia4`: ダイヤ4
- `promo`: プロモ

### ExpansionSeries（拡張シリーズ）
**定義**: 拡張パックのシリーズ分類  
**英語**: `expansion_series`  
**日本語**: 拡張シリーズ

**値**:
- `scarlet_violet`: スカーレット&バイオレット
- `sword_shield`: ソード&シールド
- `sun_moon`: サン&ムーン
- `xy`: XY
- `black_white`: ブラック&ホワイト

---

## プロセス・アクション

### TierListCreation（ティアリスト作成）
**定義**: ユーザーがティアリストを作成する一連のプロセス  
**フロー**: デッキ選択 → ティア配置 → 保存・公開  
**英語**: `tier_list_creation`  
**日本語**: ティアリスト作成

**関連アクション**:
- `PlaceDeck` - デッキ配置
- `MoveBetweenTiers` - ティア間移動
- `PublishTierList` - ティアリスト公開

---

### ImageComposition（画像合成）
**定義**: 複数カードから1つのデッキ画像を生成  
**処理**: カード画像 → レイアウト → 合成 → キャッシュ  
**英語**: `image_composition`  
**日本語**: 画像合成

**関連概念**:
- `CardImageLayout` - カード配置レイアウト
- `CompositeImage` - 合成画像
- `ImageCache` - 画像キャッシュ

---

### StatisticsCalculation（統計計算）
**定義**: ティア配置データから統計情報を計算  
**種類**: 集計ティアリスト、使用率、トレンド分析  
**英語**: `statistics_calculation`  
**日本語**: 統計計算

**関連プロセス**:
- `ConsensusGeneration` - 集計ティアリスト生成
- `TrendAnalysis` - トレンド分析
- `UsageTracking` - 使用率追跡

---

### SeasonMigration（シーズン移行）
**定義**: 環境変化時のデータ移行処理  
**処理**: 前シーズンデータ → フィルタ → 減衰 → 新シーズン  
**英語**: `season_migration`  
**日本語**: シーズン移行

**関連概念**:
- `MigrationRules` - 移行ルール
- `DataDecay` - データ減衰
- `CarryOver` - 持ち越し処理

---

## 集約境界

### TierListAggregate（ティアリスト集約）
**ルート**: TierList  
**含有**: TierPlacement[]  
**不変条件**:
- 同一デッキは1ティアリスト内で1箇所のみ配置
- ティアランクは定義された値のみ
- 配置位置は非負整数

---

### DeckAggregate（デッキ集約）
**ルート**: Deck  
**含有**: Card[], DeckAnalysis（デッキ詳細システム）  
**不変条件**:
- 画像URLは有効なリンク
- シーズンは存在するもののみ
- 最低1つのカードを含む

### ExpansionAggregate（拡張パック集約）
**ルート**: Expansion  
**含有**: Card[]  
**不変条件**:
- リリース日は過去または現在
- コードは一意
- 同一シリーズ内で順序性を持つ

---

### SeasonAggregate（シーズン集約）
**ルート**: Season  
**含有**: Deck[], TierList[]  
**不変条件**:
- アクティブシーズンは同時に1つのみ
- 開始日 < 終了日
- 過去シーズンは変更不可

---

## API命名規則

### REST エンドポイント
```
GET    /api/v1/seasons                   - ListSeasons
GET    /api/v1/seasons/active            - GetActiveSeason  
GET    /api/v1/seasons/{season_id}       - GetSeason
GET    /api/v1/expansions                - ListExpansions
GET    /api/v1/cards                     - ListCards
POST   /api/v1/decks                     - CreateDeck
GET    /api/v1/tier-lists                - ListTierLists
GET    /api/v1/tier-lists/{tier_list_id} - GetTierList
POST   /api/v1/tier-lists                - CreateTierList
PATCH  /api/v1/tier-lists/{tier_list_id} - UpdateTierList
GET    /api/v1/consensus/{season_id}     - GetConsensusTierList
```

### 内部処理
```
CalculateConsensus()     - 集計ティアリスト計算
GenerateCompositeImage() - 合成画像生成
MigrateToNewSeason()     - シーズン移行
UpdateTierStatistics()   - ティア統計更新
```

---

**文書作成日**: 2025年8月8日  
**更新頻度**: プロジェクト進行に合わせて随時更新  
**適用範囲**: 全PokeTierプロジェクト（0, 1, 2）  
**バージョン**: 1.2
