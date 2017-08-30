package db

import (
	"sts/domain/essence"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const playerCollection = "players"

type PlayerRepo struct {
	session *mgo.Session
}

func NewPlayerRepo(db *mgo.Session) essence.PlayerRepository {
	return PlayerRepo{session: db}
}

func (repo PlayerRepo) collection() *mgo.Collection {
	collection := repo.session.DB(dbName).C(playerCollection)

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

func (repo PlayerRepo) Create(player *essence.Player) error {
	if err := repo.collection().Insert(player); err != nil {
		if mgo.IsDup(err) {
			return essence.ErrPlayerAlreadyExist
		}

		return err
	}

	return nil
}

func (repo PlayerRepo) FindByID(id string) (*essence.Player, error) {
	player := &essence.Player{}
	if err := repo.collection().Find(bson.M{"id": id}).One(player); err != nil {
		if err == mgo.ErrNotFound {
			return nil, essence.ErrPlayerNotFound
		}

		return nil, err
	}

	return player, nil
}

func (repo PlayerRepo) Find(idList []string) ([]*essence.Player, error) {
	players := make([]*essence.Player, 0, 0)

	if err := repo.collection().Find(bson.M{"id": bson.M{"$in": idList}}).All(&players); err != nil {
		if err == mgo.ErrNotFound {
			return nil, essence.ErrPlayerNotFound
		}

		return nil, err
	}

	return players, nil
}

func (repo PlayerRepo) UpdateBalance(player *essence.Player, amount int64) error {
	_, err := repo.FindByID(player.ID)
	if err != nil {
		return err
	}

	player.Balance += amount
	if err := repo.collection().Update(bson.M{"id": player.ID}, player); err != nil {
		return err
	}

	return nil
}
