-- データベース初期化ファイル

-- アイテムテーブルの作成
CREATE TABLE IF NOT EXISTS items (
    id BIGINT NOT NULL AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    brand VARCHAR(100) NOT NULL,
    purchase_price INT NOT NULL,
    purchase_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- サンプルデータの挿入
INSERT INTO items (name, category, brand, purchase_price, purchase_date) VALUES 
('ロレックス デイトナ', '時計', 'ROLEX', 1500000, '2023-01-15'),
('エルメス バーキン', 'バッグ', 'HERMÈS', 2000000, '2023-02-20'),
('ティファニー ネックレス', 'ジュエリー', 'TIFFANY & CO.', 300000, '2023-03-10'),
('ルブタン パンプス', '靴', 'Christian Louboutin', 150000, '2023-04-05'),
('アップルウォッチ', 'その他', 'Apple', 50000, '2023-05-01'); 