# サイバーセキュリティ演習

このアプリには以下のユーザーが登録されています。

| ユーザー名 | パスワード |
| ---------- | ---------- |
| user1      | password1  |
| user2      | password2  |

## 環境構築と起動

1. DockerCompose でビルドする

   ```bash
   docker-compose up -d
   ```

1. DockerCompose でコンテナを起動する

   ```bash
   docker-compose up -d
   ```

   <http://localhost:3000/login> では、Web アプリケーションに直接アクセスできます。
   <http://localhost:8080/login> は、WAF(Web Application Firewall)によって保護されています。

<http://localhost:3000/login> にアクセスし、ログインフォームに以下の値を入力してログインを試みて、成功することを確認しましょう。また、パスワードを適当なものに変更してみて、ログインが失敗することを確認しましょう。

- Username: `user1`
- Password: `password1`

## SQL インジェクション攻撃

この Web アプリケーションには、SQL インジェクションの脆弱性があります。app/main.go の 42 行目で、ユーザー入力をそのままクエリに組み込んでいます。

```go
42 | db.Raw("SELECT * FROM users WHERE username = '" + username + "' AND password = '" + password + "'").Scan(&user)
```

### 攻撃が成功する例

1. <http://localhost:3000/login> にアクセスし、ログインフォームに以下の値を入力する

   - Username: `user1`
   - Password: `' OR '1'='1`

2. 「ようこそ、user1 さん！」と表示される

攻撃が成功する理由は以下の通りです：

1. パスワードに入力された `' OR '1'='1` が SQL クエリに組み込まれると、以下のようなクエリが生成されます：

   SELECT \* FROM users WHERE username = 'user1' AND password = '' OR '1'='1'

2. このクエリは以下のように解釈されます：

   - `password = ''` は偽となりますが
   - `'1'='1'` は常に真となります
   - `OR` 演算子により、条件全体が真となります

3. その結果、WHERE 句の条件が常に真となり、username が 'user1' のユーザーレコードが返されてしまいます。

このように、入力値を直接 SQL クエリに組み込むことで、攻撃者は認証をバイパスすることができてしまいます。

### 攻撃が失敗する例

1. WAF(Web Application Firewall) によって保護された <http://localhost:8080/login> にアクセスし、ログインフォームに以下の値を入力する

- Username: `user1`
- Password: `' OR '1'='1`

1. 「403 Forbidden」が表示される

## ModSecurity のログ確認

SQL インジェクションの攻撃を検知した場合、`modsecurity/log/audit.log`にログが記録されます。

ログファイルには以下のように SQL Injection の攻撃が検知されたことが記録されています。(一部抜粋)

```json
{
  "message": "SQL Injection Attack Detected via libinjection",
  "details": {
    "match": "detected SQLi using libinjection.",
    "reference": "v839,13",
    "ruleId": "942100",
    "file": "/etc/modsecurity.d/owasp-crs/rules/REQUEST-942-APPLICATION-ATTACK-SQLI.conf",
    "lineNumber": "46",
    "data": "Matched Data: s&sos found within ARGS:password: ' OR '1' = '1",
    "severity": "2",
    "ver": "OWASP_CRS/4.7.0",
    "rev": "",
    "tags": [
      "application-multi",
      "language-multi",
      "platform-multi",
      "attack-sqli",
      "paranoia-level/1",
      "OWASP_CRS",
      "capec/1000/152/248/66",
      "PCI/6.5.2"
    ],
    "maturity": "0",
    "accuracy": "0"
  }
}
```
