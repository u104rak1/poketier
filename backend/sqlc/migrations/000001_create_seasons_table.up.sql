-- シーズン集約テーブル
CREATE TABLE seasons (
    season_id UUID PRIMARY KEY,
    name VARCHAR(5) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    is_active BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- アクティブシーズンは1つのみ（ユニーク制約）
CREATE UNIQUE INDEX idx_seasons_active ON seasons (is_active) WHERE is_active = true;

-- updated_atの自動更新用トリガー
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_seasons_updated_at
    BEFORE UPDATE ON seasons
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
