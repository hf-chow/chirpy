package database

type Chirp struct {
	ID          int    `json:"id"`
	Body        string `json:"body"`
    AuthorID    int    `json:"author_id"`
}

func (db *DB) CreateChirp(body string, author_id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:         id,
		Body:       body,
        AuthorID:   author_id,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(id int) (error) {
    dbStructure, err := db.loadDB()
    if err != nil {
        return err
    }

    delete(dbStructure.Chirps, id)
    err = db.writeDB(dbStructure)

    if err != nil {
        return err
    }
    return nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, ErrNotExist
	}

	return chirp, nil
}

func (db *DB) GetChirpAuthorID(id int) (int, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return -1, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return -1, ErrNotExist
	}

    author_id := chirp.AuthorID

	return author_id, nil
}
