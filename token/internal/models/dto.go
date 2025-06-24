package models

type IssueToken struct {
	Namespace string
}

func FromRequest(r IssueTokenRequest) IssueToken {
	return IssueToken{
		Namespace: r.Namespace,
	}
}
