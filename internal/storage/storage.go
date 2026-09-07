package storage

type Page struct {
	ID     int
	URL    string
	UserID int
}

type User struct {
	ID       int
	Username string
}

//func (p Page) Hash() (string, error) {
//	const op = "storage.Page.Hash"
//
//	hash := sha1.New()
//
//	if _, err := hash.Write([]byte(p.URL)); err != nil {
//		return "", fmt.Errorf("%s: %w", op, err)
//	}
//
//	if _, err := hash.Write([]byte(p.Username)); err != nil {
//		return "", fmt.Errorf("%s: %w", op, err)
//	}
//
//	return hex.EncodeToString(hash.Sum(nil)), nil
//}
