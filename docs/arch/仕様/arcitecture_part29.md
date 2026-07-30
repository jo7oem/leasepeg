---
Title: 29. Conformance Checklist
Description: 本章は「Conformance Checklist」を対象とし、適用範囲と満たすべき制約を規定する。
Summary: |
  「29. Conformance Checklist」の要件、判断基準、関連する前提を定義する。
  本章で規定する制約に違反する実装または運用は許容しない。
---

# 29. Conformance Checklist

## 29.1 概要

本章では本仕様への適合性確認項目を定義する。

実装は以下のチェックリストを用いて仕様準拠性を評価できる。

各項目は以下のいずれかで評価する。

- PASS
- FAIL
- NOT APPLICABLE

---

# 29.2 Authentication / Authorization

## UserSession

- [ ] UserSession の最大有効期間は 30 分以下である。
- [ ] User 権限を持つ長寿命 Credential を発行しない。
- [ ] UserSession は管理操作専用である。
- [ ] UserSession は Cookie Security 要件を満たす。

---

## DeviceCredential

- [ ] DeviceCredential は単一 Device にのみ紐づく。
- [ ] DeviceCredential は Write Only である。
- [ ] DeviceCredential で Address Read ができない。
- [ ] DeviceCredential で User Management ができない。
- [ ] DeviceCredential は即時失効できる。
- [ ] DeviceCredential は expires_at を持つ。
- [ ] DeviceCredential は last_used_at を保持する。

---

## GroupReadToken

- [ ] GroupReadToken は単一 ManagementGroup に紐づく。
- [ ] GroupReadToken は Read Only である。
- [ ] GroupReadToken で Update ができない。
- [ ] GroupReadToken は即時失効できる。
- [ ] GroupReadToken は expires_at を持つ。
- [ ] GroupReadToken は last_used_at を保持する。

---

## Authorization

- [ ] Authentication と Authorization が分離されている。
- [ ] UUIDv7 に依存して Authorization を省略していない。
- [ ] User は自分の Resource のみ操作できる。
- [ ] Device は自 Device のみ操作できる。
- [ ] GroupReadToken は自 Group のみ参照できる。

---

# 29.3 Credential Security

## Secret Storage

- [ ] Bearer Token を平文保存していない。
- [ ] Password を平文保存していない。
- [ ] Session Secret を平文保存していない。
- [ ] Credential Secret は Hash 保存される。

---

## Token Policy

- [ ] Token 長は 32 文字以上である。
- [ ] Token Prefix は 16 Random Bytes 以上である。
- [ ] Token Prefix は認証に使用されない。
- [ ] Token Prefix の有効性を Error Response に含めない。

---

## Lifetime Policy

- [ ] DeviceCredential に expires_at がある。
- [ ] GroupReadToken に expires_at がある。
- [ ] Default Cap は 365 日以下で設定されている。
- [ ] 最大 Cap は 99 年以下である。

---

# 29.4 Address Model

## Address Family

- [ ] IPv4 をサポートする。
- [ ] IPv6 をサポートする。
- [ ] Dual Stack をサポートする。

---

## Device Address

- [ ] Device ごとに IPv4 は最大 1 件。
- [ ] Device ごとに IPv6 は最大 1 件。
- [ ] UNIQUE(device_id, family) を保証している。

---

## submitted_ip

- [ ] submitted_ip の形式を検証する。
- [ ] submitted_ip 所有性検証を行わない。
- [ ] submitted_ip 到達性検証を行わない。

---

## Lifecycle

- [ ] active を実装している。
- [ ] expired を実装している。
- [ ] deleted を実装している。
- [ ] Read API は active のみ返す。

---

# 29.5 Address Update Semantics

## update

- [ ] update は対象 Family のみ変更する。

---

## replace

- [ ] replace は対象 Family のみ残す。

---

## nop

- [ ] nop は Address Value を変更しない。
- [ ] nop + submitted_ip を拒否する。
- [ ] nop + submitted_ip は 400 を返す。

---

## Idempotency

- [ ] update は Retry Safe である。
- [ ] replace は Retry Safe である。
- [ ] nop は Retry Safe である。

---

# 29.6 Metadata

## Metadata Constraints

- [ ] Metadata Key は UTF-8 である。
- [ ] Metadata Key は 256 Bytes 以下である。
- [ ] Metadata Value は UTF-8 である。
- [ ] Metadata Value は 128 KiB 以下である。
- [ ] Device あたり Metadata は最大 3 件である。

---

## Metadata Operations

- [ ] Empty String を有効値として扱う。
- [ ] Null を Delete として扱う。
- [ ] UNIQUE(device_id, key) を保証している。

---

## Metadata History

- [ ] Metadata Value を UpdateHistory に保存しない。
- [ ] Metadata Create を記録する。
- [ ] Metadata Update を記録する。
- [ ] Metadata Delete を記録する。

---

# 29.7 Deletion Semantics

## Device Delete

- [ ] DeviceCredential を失効する。
- [ ] Device Address を削除する。
- [ ] Device Metadata を削除する。
- [ ] Membership を削除する。
- [ ] AuditLog を削除しない。

---

## ManagementGroup Delete

- [ ] GroupReadToken を失効する。
- [ ] Membership を削除する。
- [ ] Device を削除しない。
- [ ] AuditLog を削除しない。

---

## User Delete

- [ ] UserSession を失効する。
- [ ] Device Delete Semantics を適用する。
- [ ] Group Delete Semantics を適用する。
- [ ] AuditLog を保持する。

---

# 29.8 AuditLog

## Recording

- [ ] Authentication Success を記録する。
- [ ] Authentication Failure を記録する。
- [ ] Authorization Failure を記録する。
- [ ] Address Update を記録する。
- [ ] Metadata Update を記録する。
- [ ] Credential Operation を記録する。
- [ ] Delete Operation を記録する。

---

## Retention

- [ ] 保持期間は既定 13 ヶ月である。
- [ ] CleanupJob による削除が可能である。
- [ ] User Delete で削除されない。
- [ ] Device Delete で削除されない。
- [ ] Group Delete で削除されない。

---

## Security

- [ ] Secret を保存していない。
- [ ] Metadata Value を保存していない。

---

# 29.9 UpdateHistory

- [ ] Address Update History を保持する。
- [ ] Address 値を記録する。
- [ ] Metadata Create を記録する。
- [ ] Metadata Update を記録する。
- [ ] Metadata Delete を記録する。
- [ ] Metadata Value を保存していない。

---

# 29.10 Rate Limit

## Authentication Phase

- [ ] Source IP の Rate Limit がある。
- [ ] Token Prefix の Rate Limit がある。

---

## Authorization Phase

- [ ] Credential ID の Rate Limit がある。
- [ ] Device ID の Rate Limit がある。
- [ ] User ID の Rate Limit がある。

---

## Security

- [ ] Rate Limit で Resource 情報を漏洩しない。
- [ ] Authentication Failure を AuditLog に記録する。

---

# 29.11 CleanupJob

- [ ] Expired Address を削除できる。
- [ ] 古い UpdateHistory を削除できる。
- [ ] 古い AuditLog を削除できる。
- [ ] CleanupJob は AuditLog を記録する。
- [ ] CleanupJob は Retry Safe である。
- [ ] Cleanup API は既定で非公開である。

---

# 29.12 External Integration

- [ ] 非同期実行である。
- [ ] Retry Policy が存在する。
- [ ] Dead Letter Policy が存在する。
- [ ] Failure Count を保持する。
- [ ] Suspension State を保持する。
- [ ] Resume Operation が存在する。
- [ ] Failure を AuditLog に記録する。

---

# 29.13 Monitoring

- [ ] Monitoring API が存在する。
- [ ] Monitoring API は別 Port を利用する。
- [ ] Monitoring API は既定で非公開である。
- [ ] Liveness Endpoint が存在する。
- [ ] Readiness Endpoint が存在する。
- [ ] Metrics Endpoint が存在する。

---

# 29.14 SQL Constraints

- [ ] UNIQUE(device_id, family) が存在する。
- [ ] UNIQUE(device_id, key) が存在する。
- [ ] UNIQUE(management_group_id, device_id) が存在する。
- [ ] UNIQUE(public_id) が存在する。
- [ ] UNIQUE(type, token_prefix) が存在する。
- [ ] Foreign Key 制約が存在する。
- [ ] Ownership Constraint が保証される。

---

# 29.15 Final Conformance

実装は以下を満たす場合、本仕様に適合するとみなせる。

- [ ] Authentication Requirements
- [ ] Authorization Requirements
- [ ] Address Requirements
- [ ] Metadata Requirements
- [ ] Audit Requirements
- [ ] UpdateHistory Requirements
- [ ] Cleanup Requirements
- [ ] External Integration Requirements
- [ ] Monitoring Requirements
- [ ] SQL Constraint Requirements
- [ ] Security Requirements

すべて PASS の場合、本仕様への適合と判定する。
