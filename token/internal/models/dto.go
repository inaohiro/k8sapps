package models

type IssueToken struct {
	Namespace string
}

// todo: validator で namespace は空文字列でないことを確認する
func FromRequest(r IssueTokenRequest) IssueToken {
	return IssueToken{
		Namespace: r.Namespace,
	}
}
