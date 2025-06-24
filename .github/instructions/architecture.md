以下のアーキテクチャガイドを遵守した実装としてください

# Web API アーキテクチャガイド

## 1. アーキテクチャの概要

本プロジェクトでは、Golang を使用したアプリケーションで、以下の レイヤードアーキテクチャを採用しています

### 1.1 レイヤー構造

```
[Presentation Layer]
  └─ Controller

[Application Layer]
  ├─ Service
  └─ DTO

[Infrastructure Layer]
  └─ Model
```

## 2. パッケージ構造

基本パッケージ構造は以下のとおりです

```
core/                # 機能間で共通利用したい関数など
server/
  ├─ public.go
  └─ internal/
     ├─ controller/
     ├─ model/
     ├─ dto/
     └─ service/
```

## 3. コンポーネントの役割と責務

### 3.1 API との連携

- 外部 API との通信は専用の Client クラスを作成
- Client 層は Infrastructure Layer に配置
- レスポンスは DTO 形式で返却

## 4. レイヤー間の依存関係ルール

### 4.1 基本ルール

- 上位レイヤーから下位レイヤーへの依存のみ許可
- 同一レイヤー内での依存は許可
- 下位レイヤーから上位レイヤーへの依存は禁止

### 4.2 具体的な依存ルール

1. Controller
   - Service への依存可
2. Service
   - Controller への依存禁止

## 5. 命名規則

### 5.1 クラス名

- Controller: `〇〇Controller`
- Service: `〇〇Service`

### 5.2 メソッド名

- Controller
  - 一覧表示: `index`
  - 詳細表示: `detail`
  - 作成: `create`
  - 発行: `issue`
  - 更新: `update`
  - 削除: `delete`
- Service
  - 取得: `get〇〇`
  - 作成: `create〇〇`
  - 発行: `issue〇〇`
  - 更新: `update〇〇`
  - 削除: `delete〇〇`
