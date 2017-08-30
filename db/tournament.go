package db

import (
	"sts/domain/essence"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const TournamentCollection = "tournaments"

type TournamentRepo struct {
	session *mgo.Session
}

func NewTournamentRepo(db *mgo.Session) essence.TournamentRepository {
	return TournamentRepo{session: db}
}

func (repo TournamentRepo) collection() *mgo.Collection {
	collection := repo.session.DB(dbName).C(TournamentCollection)

	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := collection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	return collection
}

func (repo TournamentRepo) Create(tournament *essence.Tournament) error {
	if err := repo.collection().Insert(tournament); err != nil {
		if mgo.IsDup(err) {
			return essence.ErrTournamentAlreadyExist
		}

		return err
	}

	return nil
}

func (repo TournamentRepo) FindByID(id uint) (*essence.Tournament, error) {
	tournament := &essence.Tournament{}
	if err := repo.collection().Find(bson.M{"id": id}).One(tournament); err != nil {
		if err == mgo.ErrNotFound {
			return nil, essence.ErrTournamentNotFound
		}

		return nil, err
	}

	return tournament, nil
}

func (repo TournamentRepo) Update(tournament *essence.Tournament) error {
	_, err := repo.FindByID(tournament.ID)
	if err != nil {
		return err
	}

	if err := repo.collection().Update(bson.M{"id": tournament.ID}, tournament); err != nil {
		return err
	}

	return nil
}
