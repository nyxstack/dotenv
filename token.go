package dotenv

// TokenType represents the type of token in the .env file
type TokenType int

const (
	TokenComment TokenType = iota
	TokenKey
	TokenAssign
	TokenValue
	TokenNewline
	TokenEOF
)

// Token represents a lexical token
type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

// String returns a string representation of the token type
func (t TokenType) String() string {
	switch t {
	case TokenComment:
		return "COMMENT"
	case TokenKey:
		return "KEY"
	case TokenAssign:
		return "ASSIGN"
	case TokenValue:
		return "VALUE"
	case TokenNewline:
		return "NEWLINE"
	case TokenEOF:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}
